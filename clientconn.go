package grpc4go

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"sync"
	"sync/atomic"
)

var (
	ErrServerNotFound = errors.New("server not found")
)

type ClientConn struct {
	pool  *ClientPool
	retry int32
}

func (this *ClientConn) Invoke(ctx context.Context, method string, args, reply interface{}, opts ...grpc.CallOption) error {
	for i := int32(0); i <= this.retry; i++ {
		var conn, err = this.pool.Get()
		if err != nil {
			return err
		}
		return conn.Invoke(ctx, method, args, reply, opts...)
	}
	return ErrServerNotFound
}

func (this *ClientConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	for i := int32(0); i <= this.retry; i++ {
		var conn, err = this.pool.Get()
		if err != nil {
			return nil, err
		}
		return conn.NewStream(ctx, desc, method, opts...)
	}
	return nil, ErrServerNotFound
}

func (this *ClientConn) Close() {
	this.pool.Close()
}

type DialFun func() (*grpc.ClientConn, error)

type ClientPool struct {
	dial  DialFun
	size  int32
	next  int32
	mu    sync.Mutex
	conns []*grpc.ClientConn
}

func NewClientPool(size int32, fn DialFun) *ClientPool {
	var p = &ClientPool{}
	p.dial = fn
	p.size = size
	p.conns = make([]*grpc.ClientConn, p.size)
	return p
}

func (this *ClientPool) Get() (*grpc.ClientConn, error) {
	var next = atomic.AddInt32(&this.next, 1)
	var index = next % this.size

	var conn = this.conns[index]
	if conn != nil && this.checkState(conn) {
		return conn, nil
	}

	if conn != nil {
		conn.Close()
	}

	this.mu.Lock()
	defer this.mu.Unlock()

	conn = this.conns[index]
	if conn != nil && this.checkState(conn) {
		return conn, nil
	}

	conn, err := this.dial()
	if err != nil {
		return nil, err
	}
	this.conns[index] = conn
	return conn, nil
}

func (this *ClientPool) checkState(conn *grpc.ClientConn) bool {
	var state = conn.GetState()
	switch state {
	case connectivity.TransientFailure, connectivity.Shutdown:
		return false
	}
	return true
}

func (this *ClientPool) Close() {
	this.mu.Lock()
	defer this.mu.Unlock()
	for _, conn := range this.conns {
		if conn == nil {
			continue
		}
		conn.Close()
	}
}
