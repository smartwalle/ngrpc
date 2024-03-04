package ngrpc

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
	retry int
}

func (cc *ClientConn) Prepare() {
	cc.pool.Prepare()
}

func (cc *ClientConn) Invoke(ctx context.Context, method string, args, reply interface{}, opts ...grpc.CallOption) error {
	for i := 0; i <= cc.retry; i++ {
		var conn, err = cc.pool.Get()
		if err != nil {
			return err
		}
		return conn.Invoke(ctx, method, args, reply, opts...)
	}
	return ErrServerNotFound
}

func (cc *ClientConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	for i := 0; i <= cc.retry; i++ {
		var conn, err = cc.pool.Get()
		if err != nil {
			return nil, err
		}
		return conn.NewStream(ctx, desc, method, opts...)
	}
	return nil, ErrServerNotFound
}

func (cc *ClientConn) Close() {
	cc.pool.Close()
}

type DialFun func() (*grpc.ClientConn, error)

type ClientPool struct {
	dial     DialFun
	connList []*grpc.ClientConn
	mu       sync.Mutex
	size     int
	next     uint32
}

func NewClientPool(size int, fn DialFun) *ClientPool {
	var p = &ClientPool{}
	p.dial = fn
	p.size = size
	p.connList = make([]*grpc.ClientConn, p.size)
	return p
}

func (p *ClientPool) Prepare() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for idx := range p.connList {
		var conn = p.connList[idx]
		if conn == nil {
			nConn, _ := p.dial()
			if nConn != nil {
				p.connList[idx] = nConn
			}
		}
	}
}

func (p *ClientPool) Get() (*grpc.ClientConn, error) {
	var index = int(atomic.AddUint32(&p.next, 1)-1) % p.size

	p.mu.Lock()
	defer p.mu.Unlock()

	var conn = p.connList[index]
	if conn != nil && p.checkState(conn) {
		return conn, nil
	}

	if conn != nil {
		conn.Close()
	}

	conn = p.connList[index]
	if conn != nil && p.checkState(conn) {
		return conn, nil
	}

	conn, err := p.dial()
	if err != nil {
		return nil, err
	}
	p.connList[index] = conn
	return conn, nil
}

func (p *ClientPool) checkState(conn *grpc.ClientConn) bool {
	var state = conn.GetState()
	switch state {
	case connectivity.TransientFailure, connectivity.Shutdown:
		return false
	}
	return true
}

func (p *ClientPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, conn := range p.connList {
		if conn == nil {
			continue
		}
		conn.Close()
	}
}
