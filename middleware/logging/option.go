package logging

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	kLogUUID = "log-uuid"
)

type Option struct {
	grpc.EmptyCallOption
	apply func(*option)
}

type option struct {
	logger  Logger
	payload bool
}

func Disable() Option {
	return WithLogger(&nilLogger{})
}

func WithLogger(logger Logger) Option {
	if logger == nil {
		logger = &nilLogger{}
	}
	return Option{
		apply: func(opt *option) {
			opt.logger = logger
		},
	}
}

func WithPayload(payload bool) Option {
	return Option{
		apply: func(opt *option) {
			opt.payload = payload
		},
	}
}

func mergeOptions(opt *option, callOptions []Option) *option {
	if len(callOptions) == 0 {
		return opt
	}
	var nOpt = &option{}
	*nOpt = *opt
	for _, f := range callOptions {
		f.apply(nOpt)
	}
	return nOpt
}

func filterOptions(inOpts []grpc.CallOption) (grpcOptions []grpc.CallOption, retryOptions []Option) {
	for _, inOpt := range inOpts {
		if opt, ok := inOpt.(Option); ok {
			retryOptions = append(retryOptions, opt)
		} else {
			grpcOptions = append(grpcOptions, inOpt)
		}
	}
	return grpcOptions, retryOptions
}

func getUUID(ctx context.Context) (string, context.Context) {
	var in, _ = metadata.FromIncomingContext(ctx)
	if in == nil {
		in = metadata.MD{}
	}
	var ids = in.Get(kLogUUID)

	var id string
	if len(ids) > 0 && ids[0] != "" {
		id = ids[0]
	} else {
		id = uuid.NewString()
	}

	var out, _ = metadata.FromOutgoingContext(ctx)
	if out == nil {
		out = metadata.MD{}
	}
	out.Set(kLogUUID, id)

	var nCtx = metadata.NewOutgoingContext(ctx, out)
	return id, nCtx
}
