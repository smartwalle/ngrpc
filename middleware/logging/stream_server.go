package logging

import (
	"context"
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
		var id, nCtx = getUUID(ss.Context())

		opt.logger.Printf("[%s] GRPC 流建立成功: [%s] \n", id, info.FullMethod)

		var nStream = &serverStream{
			ServerStream: ss,
			ctx:          nCtx,
			logId:        id,
			opt:          opt,
		}

		var start = time.Now()
		var err = handler(srv, nStream)
		var end = time.Now()
		if err != nil && err != io.EOF {
			opt.logger.Printf("[%s] GRPC 流异常关闭: [%s], 持续时间: [%v], 错误信息: [%v] \n", id, info.FullMethod, end.Sub(start), err)
		} else {
			opt.logger.Printf("[%s] GRPC 流正常关闭: [%s], 持续时间: [%v] \n", id, info.FullMethod, end.Sub(start))
		}
		return err
	}
}

type serverStream struct {
	grpc.ServerStream
	ctx   context.Context
	logId string
	opt   *option
}

func (this *serverStream) Context() context.Context {
	return this.ctx
}

func (this *serverStream) SendMsg(m interface{}) error {
	var err = this.ServerStream.SendMsg(m)
	if this.opt.payload {
		if err != nil {
			this.opt.logger.Printf("[%s] GRPC 流发送消息失败: [%v], 错误信息: [%v] \n", this.logId, m, err)
		} else {
			this.opt.logger.Printf("[%s] GRPC 流发送消息成功: [%v] \n", this.logId, m)
		}
	}
	return err
}

func (this *serverStream) RecvMsg(m interface{}) error {
	var err = this.ServerStream.RecvMsg(m)
	if this.opt.payload {
		if err != nil {
			this.opt.logger.Printf("[%s] GRPC 流接收消息失败: [%v] \n", this.logId, err)
		} else {
			this.opt.logger.Printf("[%s] GRPC 流接收消息成功: [%v] \n", this.logId, m)
		}
	}
	return err
}
