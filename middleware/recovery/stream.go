package recovery

import "google.golang.org/grpc"

func WithStream(opts ...Option) grpc.ServerOption {
	var defaultOption = &option{}
	defaultOption = mergeOptions(defaultOption, opts)
	return grpc.ChainStreamInterceptor(streamRecovery(defaultOption))
}

func streamRecovery(defaultOption *option) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = recoverFrom(ss.Context(), defaultOption, r)
			}
		}()
		err = handler(srv, ss)
		return err
	}
}
