package grpc4go

import (
	"context"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	domain   string
	service  string
	node     string
	addr     string
	registry Registry
	server   *grpc.Server
}

func NewServer(domain, service, node, addr string, registry Registry, opts ...grpc.ServerOption) *Server {
	var s = &Server{}
	s.domain = domain
	s.service = service
	s.node = node
	s.addr = addr
	s.registry = registry
	s.server = grpc.NewServer(opts...)
	return s
}

func (this *Server) Service() string {
	return this.service
}

func (this *Server) Node() string {
	return this.node
}

func (this *Server) Addr() string {
	return this.addr
}

func (this *Server) Server() *grpc.Server {
	return this.server
}

func (this *Server) Run() error {
	listen, err := net.Listen("tcp", this.addr)
	if err != nil {
		return err
	}
	this.addr = listen.Addr().String()

	if this.registry != nil {
		this.registry.Register(context.Background(), this.domain, this.service, this.node, this.addr, 15)
	}
	if err = this.server.Serve(listen); err != nil {
		this.Stop()
		return err
	}
	return nil
}

func (this *Server) Stop() {
	if this.registry != nil {
		this.registry.Deregister(context.Background(), this.domain, this.service, this.node)
	}
	this.server.Stop()
}

func (this *Server) GracefulStop() {
	if this.registry != nil {
		this.registry.Deregister(context.Background(), this.domain, this.service, this.node)
	}
	this.server.GracefulStop()
}
