package grpc4go

import (
	"github.com/smartwalle/pool4go"
	"google.golang.org/grpc"
)

func NewPool(target string, maxActive, maxIdle int, opts ...grpc.DialOption) (p *Pool) {
	var dialFunc = func() (pool4go.Conn, error) {
		c, err := grpc.Dial(target, opts...)
		if err != nil {
			return nil, err
		}
		return c, nil
	}

	p = &Pool{}
	p.Pool = pool4go.NewPool(dialFunc)
	p.Pool.SetMaxIdleConns(maxIdle)
	p.Pool.SetMaxOpenConns(maxActive)
	return p
}

type Pool struct {
	*pool4go.Pool
}

func (this *Pool) GetConn() *grpc.ClientConn {
	var c, err = this.Pool.Get()
	if err != nil {
		return nil
	}
	if c == nil {
		return nil
	}
	var cc = c.(*grpc.ClientConn)
	return cc
}

func (this *Pool) Release(c *grpc.ClientConn) error {
	this.Pool.Release(c, false)
	return nil
}
