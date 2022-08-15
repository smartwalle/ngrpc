package retry

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

// WithUnaryClient 客户端普通方法调用重试处理
func WithUnaryClient(opts ...Option) grpc.DialOption {
	var defaultOptions = &option{
		max:         1,
		callTimeout: 5 * time.Second,
		codes:       []codes.Code{codes.ResourceExhausted, codes.Unavailable},
		backoff: func(ctx context.Context, i int) time.Duration {
			return time.Second * 1
		},
	}
	defaultOptions = mergeOptions(defaultOptions, opts)
	return grpc.WithChainUnaryInterceptor(unaryClientRetry(defaultOptions))
}

func unaryClientRetry(defaultOptions *option) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var grpcOpts, nOpts = filterOptions(opts)
		var opt = mergeOptions(defaultOptions, nOpts)

		var err error

		for i := 0; i <= opt.max; i++ {
			if i > 0 {
				if err = retryBackoff(i, ctx, opt); err != nil {
					return err
				}
			}

			var nCtx = callContext(ctx, opt)
			err = invoker(nCtx, method, req, reply, cc, grpcOpts...)

			if err == nil {
				return nil
			}

			if isContextError(err) {
				if ctx.Err() != nil {
					return err
				} else if opt.callTimeout != 0 {
					continue
				}
			}

			if isRetriable(err, opt) == false {
				return err
			}
		}
		return err
	}
}

func retryBackoff(i int, ctx context.Context, opts *option) error {
	var waitTime time.Duration = 0
	if i > 0 && opts.backoff != nil {
		waitTime = opts.backoff(ctx, i)
	}
	if waitTime > 0 {
		var timer = time.NewTimer(waitTime)
		select {
		case <-ctx.Done():
			timer.Stop()
			return contextErrToGRPCErr(ctx.Err())
		case <-timer.C:
		}
	}
	return nil
}

func callContext(ctx context.Context, opts *option) context.Context {
	var nCtx = ctx
	if opts.callTimeout > 0 {
		nCtx, _ = context.WithTimeout(nCtx, opts.callTimeout)
	}
	return nCtx
}

func isContextError(err error) bool {
	code := status.Code(err)
	return code == codes.DeadlineExceeded || code == codes.Canceled
}

func isRetriable(err error, callOption *option) bool {
	var errCode = status.Code(err)
	if isContextError(err) {
		return false
	}
	for _, code := range callOption.codes {
		if code == errCode {
			return true
		}
	}
	return false
}

func contextErrToGRPCErr(err error) error {
	switch err {
	case context.DeadlineExceeded:
		return status.Error(codes.DeadlineExceeded, err.Error())
	case context.Canceled:
		return status.Error(codes.Canceled, err.Error())
	default:
		return status.Error(codes.Unknown, err.Error())
	}
}
