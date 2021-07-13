package grpc4go

import (
	"google.golang.org/grpc"
)

type ServerOption struct {
	grpc.EmptyServerOption
	apply func(option *serverOption)
}

type serverOption struct {
	registerTTL int64
	addr        string
}

func WithRegisterTTL(ttl int64) ServerOption {
	return ServerOption{apply: func(opt *serverOption) {
		opt.registerTTL = ttl
	}}
}

func WithAddress(addr string) ServerOption {
	return ServerOption{apply: func(opt *serverOption) {
		opt.addr = addr
	}}
}

func filterServerOptions(inOpts []grpc.ServerOption) (grpcOptions []grpc.ServerOption, serverOptions []ServerOption) {
	for _, inOpt := range inOpts {
		if opt, ok := inOpt.(ServerOption); ok {
			serverOptions = append(serverOptions, opt)
		} else {
			grpcOptions = append(grpcOptions, inOpt)
		}
	}
	return grpcOptions, serverOptions
}

func mergeServerOptions(serverOpt *serverOption, serverOptions []ServerOption) *serverOption {
	if len(serverOptions) == 0 {
		return serverOpt
	}
	var nOpt = &serverOption{}
	*nOpt = *serverOpt
	for _, f := range serverOptions {
		f.apply(nOpt)
	}
	return nOpt
}
