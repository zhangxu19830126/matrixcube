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

package client

import (
	"context"
	"sync"

	"github.com/fagongzi/util/hack"
	"github.com/fagongzi/util/protoc"
	"github.com/matrixorigin/matrixcube/components/log"
	"github.com/matrixorigin/matrixcube/pb/metapb"
	"github.com/matrixorigin/matrixcube/pb/rpcpb"
	"github.com/matrixorigin/matrixcube/pb/txnpb"
	"github.com/matrixorigin/matrixcube/raftstore"
	"github.com/matrixorigin/matrixcube/util/uuid"
	"go.uber.org/zap"
)

// Option client option
type Option func(*rpcpb.Request)

// WithShardGroup set shard group to execute the request
func WithShardGroup(group uint64) Option {
	return func(req *rpcpb.Request) {
		req.Group = group
	}
}

// WithRouteKey use the specified key to route request
func WithRouteKey(key []byte) Option {
	return func(req *rpcpb.Request) {
		req.Key = key
	}
}

// WithKeysRange If the current request operates on multiple Keys, set the range [from, to) of Keys
// operated by the current request. The client needs to split the request again if it wants
// to re-route according to KeysRange after the data management scope of the Shard has
// changed, or if it returns the specified error.
func WithKeysRange(from, to []byte) Option {
	return func(req *rpcpb.Request) {
		req.KeysRange = &rpcpb.Range{From: from, To: to}
	}
}

// WithShard use the specified shard to route request
func WithShard(shard uint64) Option {
	return func(req *rpcpb.Request) {
		req.ToShard = shard
	}
}

// WithReplicaSelectPolicy set the ReplicaSelectPolicy for request, default is SelectLeader
func WithReplicaSelectPolicy(policy rpcpb.ReplicaSelectPolicy) Option {
	return func(req *rpcpb.Request) {
		req.ReplicaSelectPolicy = policy
	}
}

// Client is a cube client, providing read and write access to the external.
type Client interface {
	// Start start the cube client
	Start() error
	// Stop stop the cube client
	Stop() error

	// Router returns a Router with real-time updated routing table information
	// inside for custom message routing
	Router() raftstore.Router

	// Admin exec the admin request, and use the `Future` to get the response.
	Admin(ctx context.Context, requestType uint64, payload []byte, opts ...Option) *Future
	// Write exec the write request, and use the `Future` to get the response.
	Write(ctx context.Context, requestType uint64, payload []byte, opts ...Option) *Future
	// Read exec the read request, and use the `Future` to get the response
	Read(ctx context.Context, requestType uint64, payload []byte, opts ...Option) *Future
	// Txn exec the transaction request, and use the `Future` to get the response
	Txn(ctx context.Context, request txnpb.TxnBatchRequest, opts ...Option) *Future

	// AddLabelToShard add lable to shard, and use the `Future` to get the response
	AddLabelToShard(ctx context.Context, name, value string, shard uint64) *Future
}

var _ Client = (*client)(nil)

// client a tcp application server
type client struct {
	logger      *zap.Logger
	shardsProxy raftstore.ShardsProxy
	inflights   sync.Map // request id -> *Future

}

// NewClient creates and return a cube client
func NewClient(cfg Cfg) Client {
	return NewClientWithOptions(CreateWithLogger(cfg.Store.GetConfig().Logger.Named("cube-client")),
		CreateWithShardsProxy(cfg.Store.GetShardsProxy()))
}

// NewClientWithOptions create client with options
func NewClientWithOptions(options ...CreateOption) Client {
	c := &client{}
	for _, opt := range options {
		opt(c)
	}
	c.adjust()
	return c
}

func (s *client) adjust() {
	if s.logger == nil {
		s.logger = log.Adjust(nil).Named("cube-client")
	}

	if s.shardsProxy == nil {
		s.logger.Fatal("ShardsProxy not set")
	}
}

func (s *client) Start() error {
	s.logger.Info("begin to start cube client")
	s.shardsProxy.SetCallback(s.done, s.doneError)
	s.shardsProxy.SetRetryController(s)
	s.logger.Info("cube client started")
	return nil
}

func (s *client) Stop() error {
	s.logger.Info("cube client stopped")
	return nil
}

func (s *client) Router() raftstore.Router {
	return s.shardsProxy.Router()
}

func (s *client) Write(ctx context.Context, requestType uint64, payload []byte, opts ...Option) *Future {
	return s.exec(ctx, requestType, payload, rpcpb.Write, nil, opts...)
}

func (s *client) Read(ctx context.Context, requestType uint64, payload []byte, opts ...Option) *Future {
	return s.exec(ctx, requestType, payload, rpcpb.Read, nil, opts...)
}

func (s *client) Admin(ctx context.Context, requestType uint64, payload []byte, opts ...Option) *Future {
	return s.exec(ctx, requestType, payload, rpcpb.Admin, nil, opts...)
}

func (s *client) Txn(ctx context.Context, request txnpb.TxnBatchRequest, opts ...Option) *Future {
	return s.exec(ctx, 0, nil, rpcpb.Txn, &request, opts...)
}

func (s *client) AddLabelToShard(ctx context.Context, name, value string, shard uint64) *Future {
	payload := protoc.MustMarshal(&rpcpb.UpdateLabelsRequest{
		Labels: []metapb.Label{{Key: name, Value: value}},
		Policy: rpcpb.Add,
	})
	return s.exec(ctx, uint64(rpcpb.CmdUpdateLabels), payload, rpcpb.Admin, nil, WithShard(shard))
}

func (s *client) exec(ctx context.Context, requestType uint64, payload []byte, cmdType rpcpb.CmdType, txnRequest *txnpb.TxnBatchRequest, opts ...Option) *Future {
	f := newFuture(ctx)

	req := rpcpb.Request{}
	req.ID = uuid.NewV4().Bytes()
	req.Type = cmdType
	req.CustomType = requestType
	req.Cmd = payload
	req.TxnBatchRequest = txnRequest
	for _, opt := range opts {
		opt(&req)
	}
	s.inflights.Store(hack.SliceToString(req.ID), asyncCtx{f: f, req: req})

	if len(req.Key) > 0 && req.ToShard > 0 {
		s.logger.Fatal("route with key and route with shard cannot be set at the same time")
	}
	if _, ok := ctx.Deadline(); !ok {
		s.logger.Fatal("cube client must use timeout context")
	}

	if ce := s.logger.Check(zap.DebugLevel, "begin to send request"); ce != nil {
		ce.Write(log.RequestIDField(req.ID))
	}

	if err := s.shardsProxy.Dispatch(req); err != nil {
		f.done(nil, nil, err)
	}
	return f
}

func (s *client) Retry(requestID []byte) (rpcpb.Request, bool) {
	id := hack.SliceToString(requestID)
	if v, ok := s.inflights.Load(id); ok {
		c := v.(asyncCtx)
		if c.f.canRetry() {
			return c.req, true
		}
	}

	return rpcpb.Request{}, false
}

func (s *client) done(resp rpcpb.Response) {
	if ce := s.logger.Check(zap.DebugLevel, "response received"); ce != nil {
		ce.Write(log.RequestIDField(resp.ID))
	}

	id := hack.SliceToString(resp.ID)
	if c, ok := s.inflights.Load(hack.SliceToString(resp.ID)); ok {
		s.inflights.Delete(id)
		c.(asyncCtx).f.done(resp.Value, resp.TxnBatchResponse, nil)
	} else {
		if ce := s.logger.Check(zap.DebugLevel, "response skipped"); ce != nil {
			ce.Write(log.RequestIDField(resp.ID), log.ReasonField("missing ctx"))
		}
	}
}

func (s *client) doneError(requestID []byte, err error) {
	if ce := s.logger.Check(zap.DebugLevel, "error response received"); ce != nil {
		ce.Write(log.RequestIDField(requestID), zap.Error(err))
	}

	id := hack.SliceToString(requestID)
	if c, ok := s.inflights.Load(id); ok {
		s.inflights.Delete(id)
		c.(asyncCtx).f.done(nil, nil, err)
	}
}

type asyncCtx struct {
	f   *Future
	req rpcpb.Request
}
