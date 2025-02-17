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

package prophet

import (
	"time"

	"github.com/fagongzi/goetty"
	"github.com/fagongzi/goetty/buf"
	"github.com/matrixorigin/matrixcube/components/prophet/codec"
	"github.com/matrixorigin/matrixcube/components/prophet/pb/metapb"
	"github.com/matrixorigin/matrixcube/components/prophet/util"
)

// Option client option
type Option func(*options)

type options struct {
	leaderGetter func() *metapb.Member
	rpcTimeout   time.Duration
}

func (opts *options) adjust() {
	if opts.rpcTimeout == 0 {
		opts.rpcTimeout = time.Second * 10
	}
}

// WithLeaderGetter set a func to get a leader
func WithLeaderGetter(value func() *metapb.Member) Option {
	return func(opts *options) {
		opts.leaderGetter = value
	}
}

// WithRPCTimeout set rpc timeout
func WithRPCTimeout(value time.Duration) Option {
	return func(opts *options) {
		opts.rpcTimeout = value
	}
}

func createConn() goetty.IOSession {
	encoder, decoder := codec.NewClientCodec(10 * buf.MB)
	return goetty.NewIOSession(goetty.WithCodec(encoder, decoder),
		goetty.WithLogger(util.GetLogger()),
		goetty.WithEnableAsyncWrite(16))

}
