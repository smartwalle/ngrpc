package ngrpc

import (
	"google.golang.org/grpc"
)

type ServerOption struct {
	grpc.EmptyServerOption
	apply func(opts *serverOptions)
}

type serverOptions struct {
	addr        string
	registerTTL int64
}

func WithRegisterTTL(ttl int64) ServerOption {
	return ServerOption{apply: func(opts *serverOptions) {
		opts.registerTTL = ttl
	}}
}

func WithAddress(addr string) ServerOption {
	return ServerOption{apply: func(opts *serverOptions) {
		opts.addr = addr
	}}
}

func mergeServerOptions(dOpts *serverOptions, opts []ServerOption) *serverOptions {
	if len(opts) == 0 {
		return dOpts
	}
	var nOpts = &serverOptions{}
	*nOpts = *dOpts
	for _, f := range opts {
		f.apply(nOpts)
	}
	return nOpts
}

func filterServerOptions(inOpts []grpc.ServerOption) (grpcOptions []grpc.ServerOption, nOptions []ServerOption) {
	for _, inOpt := range inOpts {
		if opt, ok := inOpt.(ServerOption); ok {
			nOptions = append(nOptions, opt)
		} else {
			grpcOptions = append(grpcOptions, inOpt)
		}
	}
	return grpcOptions, nOptions
}
