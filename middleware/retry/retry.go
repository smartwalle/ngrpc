// https://github.com/grpc-ecosystem/go-grpc-middleware/blob/master/retry/retry.go

package retry

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

// WithUnaryCallRetry 普通方法调用重试处理
func WithUnaryCallRetry(opts ...CallOption) grpc.DialOption {
	var defaultOption = &option{
		max:     1,
		timeout: 5 * time.Second,
		codes:   []codes.Code{codes.ResourceExhausted, codes.Unavailable},
		backoff: func(ctx context.Context, i int) time.Duration {
			return time.Second * 1
		},
	}
	defaultOption = mergeOptions(defaultOption, opts)
	return grpc.WithUnaryInterceptor(unaryClientRetry(defaultOption))
}

func unaryClientRetry(defaultOption *option) grpc.UnaryClientInterceptor {
	return func(pCtx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var grpcOpts, retryOpts = filterOptions(opts)
		var callOpt = mergeOptions(defaultOption, retryOpts)

		var err error

		for i := 0; i <= callOpt.max; i++ {
			if i > 0 {
				if err = retryBackoff(i, pCtx, callOpt); err != nil {
					return err
				}
			}

			var nCtx = callContext(pCtx, defaultOption)
			err = invoker(nCtx, method, req, reply, cc, grpcOpts...)

			if err == nil {
				return nil
			}

			if isContextError(err) {
				if pCtx.Err() != nil {
					return err
				} else if callOpt.timeout != 0 {
					continue
				}
			}

			if isRetriable(err, callOpt) == false {
				return err
			}
		}
		return err
	}
}

func retryBackoff(i int, ctx context.Context, callOpt *option) error {
	var waitTime time.Duration = 0
	if i > 0 {
		waitTime = callOpt.backoff(ctx, i)
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

func callContext(pCtx context.Context, callOpt *option) context.Context {
	var nCtx = pCtx
	if callOpt.timeout != 0 {
		nCtx, _ = context.WithTimeout(nCtx, callOpt.timeout)
	}
	return nCtx
}

func isContextError(err error) bool {
	code := status.Code(err)
	return code == codes.DeadlineExceeded || code == codes.Canceled
}

func isRetriable(err error, callOpts *option) bool {
	var errCode = status.Code(err)
	if isContextError(err) {
		return false
	}
	for _, code := range callOpts.codes {
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
