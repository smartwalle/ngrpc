package logging

import (
	"context"
	"google.golang.org/grpc"
)

// WithStreamClient 客户端流操作日志处理
func WithStreamClient(opts ...Option) grpc.DialOption {
	var defaultOption = &option{
		logger:  &nilLogger{},
		payload: true,
	}
	defaultOption = mergeOptions(defaultOption, opts)
	return grpc.WithChainStreamInterceptor(streamClientLog(defaultOption))
}

func streamClientLog(defaultOption *option) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		var id, nCtx = getLogId(ctx)

		var grpcOpts, logOpts = filterOptions(opts)
		var opt = mergeOptions(defaultOption, logOpts)

		opt.logger.Printf("[%s] GRPC 请求建立流: [%s - %s] \n", id, cc.Target(), method)

		var stream, err = streamer(nCtx, desc, cc, method, grpcOpts...)
		if err != nil {
			opt.logger.Printf("[%s] GRPC 建立流失败: [%s - %s], 错误信息: [%v] \n", id, cc.Target(), method, err)
		} else {
			opt.logger.Printf("[%s] GRPC 建立流成功: [%s - %s] \n", id, cc.Target(), method)
		}

		if stream == nil {
			return nil, err
		}

		var nStream = &clientStream{
			ClientStream: stream,
			logId:        id,
			opt:          opt,
		}
		return nStream, err
	}
}

type clientStream struct {
	grpc.ClientStream
	logId string
	opt   *option
}

func (this *clientStream) SendMsg(m interface{}) error {
	var err = this.ClientStream.SendMsg(m)
	if this.opt.payload {
		if err != nil {
			this.opt.logger.Printf("[%s] GRPC 流发送消息失败: [%v], 错误信息: [%v] \n", this.logId, m, err)
		} else {
			this.opt.logger.Printf("[%s] GRPC 流发送消息成功: [%v] \n", this.logId, m)
		}
	}
	return err
}

func (this *clientStream) RecvMsg(m interface{}) error {
	var err = this.ClientStream.RecvMsg(m)
	if this.opt.payload {
		if err != nil {
			this.opt.logger.Printf("[%s] GRPC 流接收消息失败: [%v] \n", this.logId, err)
		} else {
			this.opt.logger.Printf("[%s] GRPC 流接收消息成功: [%v] \n", this.logId, m)
		}
	}
	return err
}
