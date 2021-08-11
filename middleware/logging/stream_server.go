package logging

import (
	"google.golang.org/grpc"
	"io"
	"time"
)

// WithStreamServer 服务端流调用日志处理
func WithStreamServer(opts ...Option) grpc.ServerOption {
	var defaultOption = &option{
		logger:  &nilLogger{},
		payload: true,
	}
	defaultOption = mergeOptions(defaultOption, opts)
	return grpc.ChainStreamInterceptor(streamServerLog(defaultOption))
}

func streamServerLog(defaultOption *option) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		var id, _ = getLogUUID(ss.Context())

		defaultOption.logger.Printf("[%s] GRPC 流 [%s] 建立成功 \n", id, info.FullMethod)

		var nStream = &serverStream{
			ServerStream: ss,
			logId:        id,
			opt:          defaultOption,
		}

		var start = time.Now()
		var err = handler(srv, nStream)
		var end = time.Now()
		if err != nil && err != io.EOF {
			defaultOption.logger.Printf("[%s] GRPC 流 [%s] 异常关闭, 流持续时间 [%v], 错误信息 [%v] \n", id, info.FullMethod, end.Sub(start), err)
		} else {
			defaultOption.logger.Printf("[%s] GRPC 流 [%s] 正常关闭, 流持续时间 [%v] \n", id, info.FullMethod, end.Sub(start))
		}
		return err
	}
}

type serverStream struct {
	grpc.ServerStream
	logId string
	opt   *option
}

func (this *serverStream) SendMsg(m interface{}) error {
	var err = this.ServerStream.SendMsg(m)
	if this.opt.payload {
		if err != nil {
			this.opt.logger.Printf("[%s] GRPC 流发送消息 [%v] 发生错误 [%v] \n", this.logId, m, err)
		} else {
			this.opt.logger.Printf("[%s] GRPC 流发送消息 [%v] 成功 \n", this.logId, m)
		}
	}
	return err
}

func (this *serverStream) RecvMsg(m interface{}) error {
	var err = this.ServerStream.RecvMsg(m)
	if this.opt.payload {
		if err != nil {
			this.opt.logger.Printf("[%s] GRPC 流接收消息发生错误 [%v] \n", this.logId, err)
		} else {
			this.opt.logger.Printf("[%s] GRPC 流接收消息 [%v] 成功 \n", this.logId, m)
		}
	}
	return err
}
