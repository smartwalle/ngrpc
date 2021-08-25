package logging

import (
	"google.golang.org/grpc"
	"io"
	"time"
)

// WithStreamServer 服务端流操作日志处理
func WithStreamServer(opts ...Option) grpc.ServerOption {
	var defaultOption = &option{
		logger:  &nilLogger{},
		payload: true,
	}
	defaultOption = mergeOptions(defaultOption, opts)
	return grpc.ChainStreamInterceptor(streamServerLog(defaultOption))
}

func streamServerLog(opt *option) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		opt.logger.Printf(ss.Context(), "GRPC 流建立成功: [%s] \n", info.FullMethod)

		var nStream = &serverStream{
			ServerStream: ss,
			opt:          opt,
		}

		var start = time.Now()
		var err = handler(srv, nStream)
		var end = time.Now()
		if err != nil && err != io.EOF {
			opt.logger.Printf(ss.Context(), "GRPC 流异常关闭: [%s], 持续时间: [%v], 错误信息: [%v] \n", info.FullMethod, end.Sub(start), err)
		} else {
			opt.logger.Printf(ss.Context(), "GRPC 流正常关闭: [%s], 持续时间: [%v] \n", info.FullMethod, end.Sub(start))
		}
		return err
	}
}

type serverStream struct {
	grpc.ServerStream
	opt *option
}

func (this *serverStream) SendMsg(m interface{}) error {
	var err = this.ServerStream.SendMsg(m)
	if this.opt.payload {
		if err != nil {
			this.opt.logger.Printf(this.Context(), "GRPC 流发送消息失败: [%v], 错误信息: [%v] \n", m, err)
		} else {
			this.opt.logger.Printf(this.Context(), "GRPC 流发送消息成功: [%v] \n", m)
		}
	}
	return err
}

func (this *serverStream) RecvMsg(m interface{}) error {
	var err = this.ServerStream.RecvMsg(m)
	if this.opt.payload {
		if err != nil {
			this.opt.logger.Printf(this.Context(), "GRPC 流接收消息失败: [%v] \n", err)
		} else {
			this.opt.logger.Printf(this.Context(), "GRPC 流接收消息成功: [%v] \n", m)
		}
	}
	return err
}
