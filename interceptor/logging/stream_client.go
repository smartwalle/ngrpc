package logging

import (
	"context"
	"google.golang.org/grpc"
)

// WithStreamClient 客户端流操作日志处理
func WithStreamClient(opts ...Option) grpc.DialOption {
	var defaultOptions = &options{
		logger:  &nilLogger{},
		payload: true,
	}
	defaultOptions = mergeOptions(defaultOptions, opts)
	return grpc.WithChainStreamInterceptor(streamClientLog(defaultOptions))
}

func streamClientLog(defaultOptions *options) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		var grpcOpts, nOpts = filterOptions(opts)
		var opt = mergeOptions(defaultOptions, nOpts)

		opt.logger.Printf(ctx, "GRPC 请求建立流: [%s - %s] \n", cc.Target(), method)

		var stream, err = streamer(ctx, desc, cc, method, grpcOpts...)
		if err != nil {
			opt.logger.Printf(ctx, "GRPC 建立流失败: [%s - %s], 错误信息: [%v] \n", cc.Target(), method, err)
		} else {
			opt.logger.Printf(ctx, "GRPC 建立流成功: [%s - %s] \n", cc.Target(), method)
		}

		if stream == nil {
			return nil, err
		}

		var nStream = &clientStream{
			ClientStream: stream,
			opts:         opt,
		}
		return nStream, err
	}
}

type clientStream struct {
	grpc.ClientStream
	opts *options
}

func (stream *clientStream) SendMsg(m interface{}) error {
	var err = stream.ClientStream.SendMsg(m)
	if stream.opts.payload {
		if err != nil {
			stream.opts.logger.Printf(stream.Context(), "GRPC 流发送消息失败: [%v], 错误信息: [%v] \n", m, err)
		} else {
			stream.opts.logger.Printf(stream.Context(), "GRPC 流发送消息成功: [%v] \n", m)
		}
	}
	return err
}

func (stream *clientStream) RecvMsg(m interface{}) error {
	var err = stream.ClientStream.RecvMsg(m)
	if stream.opts.payload {
		if err != nil {
			stream.opts.logger.Printf(stream.Context(), "GRPC 流接收消息失败: [%v] \n", err)
		} else {
			stream.opts.logger.Printf(stream.Context(), "GRPC 流接收消息成功: [%v] \n", m)
		}
	}
	return err
}
