package ngrpc

import (
	"context"
	"google.golang.org/grpc"
	"net"
	"path/filepath"
)

type Server struct {
	listener net.Listener
	registry Registry
	options  *serverOptions
	server   *grpc.Server
	domain   string
	service  string
	node     string
}

func NewServer(domain, service, node string, registry Registry, opts ...grpc.ServerOption) (*Server, error) {
	var defaultOptions = &serverOptions{
		registerTTL: 15,
	}

	var grpcOpts, nOpts = filterServerOptions(opts)
	var opt = mergeServerOptions(defaultOptions, nOpts)

	var listener, err = listen(opt.addr)
	if err != nil {
		return nil, err
	}

	var s = &Server{}
	s.options = opt
	s.domain = domain
	s.service = service
	s.node = node
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

func (s *Server) Name() string {
	if s.registry != nil {
		return s.registry.BuildPath(s.domain, s.service, s.node)
	}
	return filepath.Join(s.domain, s.service, s.node)
}

func (s *Server) Domain() string {
	return s.domain
}

func (s *Server) Service() string {
	return s.service
}

func (s *Server) Node() string {
	return s.node
}

func (s *Server) Addr() string {
	return s.listener.Addr().String()
}

func (s *Server) Registry() Registry {
	return s.registry
}

func (s *Server) Server() *grpc.Server {
	return s.server
}

func (s *Server) Start(ctx context.Context) (err error) {
	if s.registry != nil {
		if _, err = s.registry.Register(ctx, s.domain, s.service, s.node, s.Addr(), s.options.registerTTL); err != nil {
			return err
		}
	}
	if err = s.server.Serve(s.listener); err != nil {
		s.Stop(ctx)
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) (err error) {
	if s.registry != nil {
		err = s.registry.Unregister(ctx, s.domain, s.service, s.node)
	}
	s.server.GracefulStop()
	return err
}

func (s *Server) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
	s.server.RegisterService(desc, impl)
}
