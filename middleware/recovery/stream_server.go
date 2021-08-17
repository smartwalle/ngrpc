package recovery

import "google.golang.org/grpc"

// WithStreamServer 服务端捕获流调用异常
func WithStreamServer(opts ...Option) grpc.ServerOption {
	var defaultOption = &option{}
	defaultOption = mergeOptions(defaultOption, opts)
	return grpc.ChainStreamInterceptor(streamServerRecovery(defaultOption))
}

func streamServerRecovery(defaultOption *option) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = errorFrom(ss.Context(), defaultOption, r)
			}
		}()
		err = handler(srv, ss)
		return err
	}
}
