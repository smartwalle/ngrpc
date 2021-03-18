package grpc4go

import (
	"context"
	"github.com/smartwalle/pool4go"
	"google.golang.org/grpc"
)

type ClientConn struct {
	target string
	opts   []grpc.DialOption
	pool   pool4go.Pool
}

func NewClientConn(target string, maxIdle, maxOpen int, opts ...grpc.DialOption) *ClientConn {
	if maxIdle <= 0 {
		maxIdle = 1
	}
	if maxOpen <= 0 {
		maxOpen = 1
	}

	var c = &ClientConn{}
	c.target = target
	c.opts = opts
	c.pool = pool4go.New(func() (pool4go.Conn, error) {
		return grpc.Dial(target, opts...)
	}, pool4go.WithMaxIdle(maxIdle), pool4go.WithMaxOpen(maxOpen))
	return c
}

func (this *ClientConn) Invoke(ctx context.Context, method string, args, reply interface{}, opts ...grpc.CallOption) error {
	var conn, err = this.pool.Get()
	if err != nil {
		return err
	}

	var nConn = conn.(*grpc.ClientConn)
	err = nConn.Invoke(ctx, method, args, reply, opts...)
	this.pool.Put(conn)
	return err
}

func (this *ClientConn) Close() {
	this.pool.Close()
}
