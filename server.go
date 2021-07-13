package grpc4go

import (
	"context"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	opt      *serverOption
	domain   string
	service  string
	node     string
	addr     net.Addr
	listener net.Listener
	registry Registry
	server   *grpc.Server
}

func NewServer(domain, service, node string, registry Registry, opts ...grpc.ServerOption) (*Server, error) {
	var defaultOption = &serverOption{
		registerTTL: 15,
	}

	var grpcOpts, serverOpts = filterServerOptions(opts)
	var serverOpt = mergeServerOptions(defaultOption, serverOpts)

	var listener, err = listen(serverOpt.addr)
	if err != nil {
		return nil, err
	}

	var s = &Server{}
	s.opt = serverOpt
	s.domain = domain
	s.service = service
	s.node = node
	s.addr = listener.Addr()
	s.listener = listener
	s.registry = registry
	s.server = grpc.NewServer(grpcOpts...)
	return s, nil
}

func getInternalIP() net.IP {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP
			}
		}
	}
	return nil
}

func listen(addr string) (net.Listener, error) {
	nAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	if len(nAddr.IP) == 0 {
		nAddr.IP = getInternalIP()
	}
	listener, err := net.ListenTCP("tcp", nAddr)
	if err != nil {
		return nil, err
	}
	return listener, nil
}

func (this *Server) Name() string {
	return this.registry.BuildPath(this.domain, this.service, this.node)
}

func (this *Server) Domain() string {
	return this.domain
}

func (this *Server) Service() string {
	return this.service
}

func (this *Server) Node() string {
	return this.node
}

func (this *Server) Addr() string {
	return this.addr.String()
}

func (this *Server) Registry() Registry {
	return this.registry
}

func (this *Server) Server() *grpc.Server {
	return this.server
}

func (this *Server) Run() error {
	if this.registry != nil {
		this.registry.Register(context.Background(), this.domain, this.service, this.node, this.Addr(), this.opt.registerTTL)
	}
	if err := this.server.Serve(this.listener); err != nil {
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
