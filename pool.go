package grpc4go

import (
	"github.com/smartwalle/pool4go"
	"google.golang.org/grpc"
	"strings"
	"sync"
)

type BuildTargetFunc func(target, node string) string

// --------------------------------------------------------------------------------
func NewPool(target string, maxOpen, maxIdle int, opts ...grpc.DialOption) (p *Pool) {
	var dialFunc = func() (pool4go.Conn, error) {
		c, err := grpc.Dial(target, opts...)
		if err != nil {
			return nil, err
		}
		return c, nil
	}

	p = &Pool{}
	p.target = target
	p.opts = opts
	p.Pool = pool4go.NewPool(dialFunc)
	p.Pool.SetMaxIdleConns(maxIdle)
	p.Pool.SetMaxOpenConns(maxOpen)
	p.hub = NewPoolHub()
	return p
}

type Pool struct {
	mu            sync.Mutex
	target        string
	TargetBuilder BuildTargetFunc
	opts          []grpc.DialOption
	hub           *PoolHub
	*pool4go.Pool
}

func (this *Pool) GetConn(node ...string) *grpc.ClientConn {
	if len(node) > 0 {
		var node1 = strings.TrimSpace(node[0])
		if node1 != "" {
			var nTarget string
			if this.TargetBuilder != nil {
				nTarget = this.TargetBuilder(this.target, node1)
			} else {
				nTarget = this.target + "/" + node1
			}

			this.mu.Lock()
			defer this.mu.Unlock()

			var p = this.hub.GetPool(nTarget)
			if p == nil {
				p = NewPool(nTarget, this.MaxOpenConns(), this.MaxIdleConns(), this.opts...)
				this.hub.AddPool(nTarget, p)
			}
			return p.GetConn()
		}
	}

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
	if c.Target() != this.target {
		var p = this.hub.GetPool(c.Target())
		if p != nil {
			p.Release(c)
			return nil
		}

		c.Close()
		return nil
	}
	this.Pool.Release(c, false)
	return nil
}

func (this *Pool) Target() string {
	return this.target
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

	var op = this.ps[target]
	if op != nil {
		op.Close()
	}

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
