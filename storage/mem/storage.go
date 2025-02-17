// Copyright 2020 MatrixOrigin.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package mem

import (
	"encoding/binary"
	"fmt"
	"io"
	"sync/atomic"
	"time"

	"github.com/fagongzi/goetty/buf"
	"github.com/matrixorigin/matrixcube/pb/bhmetapb"
	"github.com/matrixorigin/matrixcube/storage/stats"
	"github.com/matrixorigin/matrixcube/util"
	"github.com/matrixorigin/matrixcube/vfs"
)

// Storage memory storage
type Storage struct {
	fs    vfs.FS
	kv    *util.KVTree
	stats stats.Stats

	// SyncCount number of `Sync` method called
	SyncCount uint64
}

// NewStorage returns a mem data storage with snapshot data backed by the
// specified vfs
func NewStorage(fs vfs.FS) *Storage {
	return &Storage{
		kv: util.NewKVTree(),
		fs: fs,
	}
}

func (s *Storage) Stats() stats.Stats {
	return s.stats
}

// Set put the key, value pair to the storage
func (s *Storage) Set(key []byte, value []byte) error {
	atomic.AddUint64(&s.stats.WrittenKeys, 1)
	atomic.AddUint64(&s.stats.WrittenBytes, uint64(len(value)+len(key)))
	return s.SetWithTTL(key, value, 0)
}

// SetWithTTL put the key, value pair to the storage with a ttl in seconds
func (s *Storage) SetWithTTL(key []byte, value []byte, ttl int32) error {
	atomic.AddUint64(&s.stats.WrittenKeys, 1)
	atomic.AddUint64(&s.stats.WrittenBytes, uint64(len(value)+len(key)))

	if ttl > 0 {
		s.kv.Put(key, encodeValue(value, time.Now().Add(time.Second*time.Duration(ttl)).Unix()))
	} else {
		s.kv.Put(key, encodeValue(value, 0))
	}
	return nil
}

// BatchSet batch set
func (s *Storage) BatchSet(pairs ...[]byte) error {
	if len(pairs)%2 != 0 {
		return fmt.Errorf("invalid args len: %d", len(pairs))
	}

	atomic.AddUint64(&s.stats.WrittenKeys, uint64(len(pairs)/2))
	for i := 0; i < len(pairs)/2; i++ {
		s.Set(pairs[2*i], pairs[2*i+1])
		atomic.AddUint64(&s.stats.WrittenBytes, uint64(len(pairs[2*i])+len(pairs[2*i+1])))
	}

	return nil
}

// Get returns the value of the key
func (s *Storage) Get(key []byte) ([]byte, error) {
	v := s.kv.Get(key)
	return decodeValue(v), nil
}

// MGet returns multi values
func (s *Storage) MGet(keys ...[]byte) ([][]byte, error) {
	var values [][]byte
	for _, key := range keys {
		values = append(values, decodeValue(s.kv.Get(key)))
	}

	return values, nil
}

// Delete remove the key from the storage
func (s *Storage) Delete(key []byte) error {
	atomic.AddUint64(&s.stats.WrittenKeys, 1)
	atomic.AddUint64(&s.stats.WrittenBytes, uint64(len(key)))
	s.kv.Delete(key)
	return nil
}

// BatchDelete batch delete
func (s *Storage) BatchDelete(keys ...[]byte) error {
	n := 0
	for _, key := range keys {
		s.kv.Delete(key)
		n += len(key)
	}

	atomic.AddUint64(&s.stats.WrittenKeys, uint64(len(keys)))
	atomic.AddUint64(&s.stats.WrittenBytes, uint64(n))
	return nil
}

// RangeDelete remove data in [start,end)
func (s *Storage) RangeDelete(start, end []byte) error {
	atomic.AddUint64(&s.stats.WrittenKeys, 2)
	atomic.AddUint64(&s.stats.WrittenBytes, uint64(len(start)+len(end)))
	s.kv.RangeDelete(start, end)
	return nil
}

// Scan scans the key-value paire in [start, end), and perform with a handler function, if the function
// returns false, the scan will be terminated, if the `pooledKey` is true, raftstore will call `Free` when
// scan completed.
func (s *Storage) Scan(start, end []byte, handler func(key, value []byte) (bool, error), pooledKey bool) error {
	return s.kv.Scan(start, end, func(key, value []byte) (bool, error) {
		value = decodeValue(value)
		if len(value) == 0 {
			return true, nil
		}

		return handler(key, value)
	})
}

// PrefixScan scans the key-value pairs starts from prefix but only keys for the same prefix,
// while perform with a handler function, if the function returns false, the scan will be terminated.
// if the `pooledKey` is true, raftstore will call `Free` when
// scan completed.
func (s *Storage) PrefixScan(prefix []byte, handler func(key, value []byte) (bool, error), pooledKey bool) error {
	return s.kv.PrefixScan(prefix, func(key, value []byte) (bool, error) {
		value = decodeValue(value)
		if len(value) == 0 {
			return true, nil
		}

		return handler(key, value)
	})
}

// Free free the pooled bytes
func (s *Storage) Free(pooled []byte) {

}

// SplitCheck Find a key from [start, end), so that the sum of bytes of the value of [start, key) <=size,
// returns the current bytes in [start,end), and the founded key
func (s *Storage) SplitCheck(start []byte, end []byte, size uint64) (uint64, uint64, [][]byte, error) {
	total := uint64(0)
	keys := uint64(0)
	sum := uint64(0)
	appendSplitKey := false
	var splitKeys [][]byte
	s.kv.Scan(start, end, func(key, value []byte) (bool, error) {
		value = decodeValue(value)
		if len(value) == 0 {
			return true, nil
		}

		if appendSplitKey {
			splitKeys = append(splitKeys, key)
			appendSplitKey = false
			sum = 0
		}

		n := uint64(len(key) + len(value))
		total += n
		sum += n
		keys++
		if sum >= size {
			appendSplitKey = true
		}
		return true, nil
	})

	return total, keys, splitKeys, nil
}

// Seek returns the first key-value that >= key
func (s *Storage) Seek(key []byte) ([]byte, []byte, error) {
	k, v := s.kv.Seek(key)
	return k, decodeValue(v), nil
}

// Sync sync data
func (s *Storage) Sync() error {
	atomic.AddUint64(&s.SyncCount, 1)
	return nil
}

// Write write the data in batch
func (s *Storage) Write(wb *util.WriteBatch, sync bool) error {
	if len(wb.Ops) == 0 {
		return nil
	}

	for idx, op := range wb.Ops {
		switch op {
		case util.OpDelete:
			s.Delete(wb.Keys[idx])
		case util.OpSet:
			s.SetWithTTL(wb.Keys[idx], wb.Values[idx], wb.TTLs[idx])
		}
	}
	return nil
}

// RemoveShardData remove shard data
func (s *Storage) RemoveShardData(shard bhmetapb.Shard, encodedStartKey, encodedEndKey []byte) error {
	return s.RangeDelete(encodedStartKey, encodedEndKey)
}

// CreateSnapshot create a snapshot file under the giving path
func (s *Storage) CreateSnapshot(path string, start, end []byte) error {
	err := s.fs.MkdirAll(path, 0755)
	if err != nil {
		return err
	}
	f, err := s.fs.Create(s.fs.PathJoin(path, "db.data"))
	if err != nil {
		return err
	}
	defer f.Close()

	err = writeBytes(f, start)
	if err != nil {
		return err
	}

	err = writeBytes(f, end)
	if err != nil {
		return err
	}

	return s.kv.Scan(start, end, func(key, value []byte) (bool, error) {
		err := writeBytes(f, key)
		if err != nil {
			return false, err
		}

		err = writeBytes(f, value)
		if err != nil {
			return false, err
		}

		return true, nil
	})
}

// ApplySnapshot apply a snapshort file from giving path
func (s *Storage) ApplySnapshot(path string) error {
	f, err := s.fs.Open(s.fs.PathJoin(path, "db.data"))
	if err != nil {
		return err
	}
	defer f.Close()

	start, err := readBytes(f)
	if err != nil {
		return err
	}
	if len(start) == 0 {
		return fmt.Errorf("error format, missing start field")
	}

	end, err := readBytes(f)
	if err != nil {
		return err
	}
	if len(end) == 0 {
		return fmt.Errorf("error format, missing end field")
	}

	s.kv.RangeDelete(start, end)
	for {
		key, err := readBytes(f)
		if err != nil {
			return err
		}
		if len(key) == 0 {
			break
		}

		value, err := readBytes(f)
		if err != nil {
			return err
		}
		if len(value) == 0 {
			return fmt.Errorf("error format, missing value field")
		}

		atomic.AddUint64(&s.stats.WrittenKeys, 1)
		atomic.AddUint64(&s.stats.WrittenBytes, uint64(len(key)+len(value)))
		s.kv.Put(key, value)
	}

	return nil
}

// Close close the storage
func (s *Storage) Close() error {
	return nil
}

func writeBytes(f vfs.File, data []byte) error {
	size := make([]byte, 4)
	binary.BigEndian.PutUint32(size, uint32(len(data)))
	_, err := f.Write(size)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func readBytes(f vfs.File) ([]byte, error) {
	size := make([]byte, 4)
	n, err := f.Read(size)
	if n == 0 && err == io.EOF {
		return nil, nil
	}

	total := int(binary.BigEndian.Uint32(size))
	written := 0
	data := make([]byte, total)
	for {
		n, err = f.Read(data[written:])
		if err != nil && err != io.EOF {
			return nil, err
		}
		written += n
		if written == total {
			return data, nil
		}
	}
}

func encodeValue(value []byte, expireAt int64) []byte {
	data := make([]byte, len(value)+8)
	buf.Int64ToBytesTo(expireAt, data)
	copy(data[8:], value)
	return data
}

func decodeValue(value []byte) []byte {
	if len(value) == 0 {
		return nil
	}

	expireAt := buf.Byte2Int64(value)
	now := time.Now().Unix()
	if expireAt != 0 && expireAt < now {
		return nil
	}

	return value[8:]
}
