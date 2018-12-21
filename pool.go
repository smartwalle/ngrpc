package grpc4go

import (
	"github.com/smartwalle/pool4go"
	"google.golang.org/grpc"
	"sync"
)

// --------------------------------------------------------------------------------
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

// --------------------------------------------------------------------------------
type PoolHub struct {
	mu sync.Mutex
	ps map[string]*Pool
}

func NewPoolHub() *PoolHub {
	var h = &PoolHub{}
	h.ps = make(map[string]*Pool)
	return h
}

func (this *PoolHub) GetPool(target string) *Pool {
	this.mu.Lock()
	defer this.mu.Unlock()
	var p = this.ps[target]
	return p
}

func (this *PoolHub) AddPool(target string, p *Pool) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.ps[target] = p
}

func (this *PoolHub) RemovePool(target string) (p *Pool) {
	this.mu.Lock()
	defer this.mu.Unlock()
	if p = this.ps[target]; p != nil {
		delete(this.ps, target)
	}
	return p
}
