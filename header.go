package grpc4go

import (
	"context"
	"google.golang.org/grpc/metadata"
)

type Header struct {
	md metadata.MD
}

func NewHeader() *Header {
	var h = &Header{}
	h.md = metadata.MD{}
	return h
}

func HeaderFrom(ctx context.Context) *Header {
	var h = &Header{}
	h.md, _ = metadata.FromIncomingContext(ctx)
	return h
}

func (this *Header) Add(key, value string) {
	this.md.Set(key, value)
}

func (this *Header) Get(key string) string {
	var vs = this.md.Get(key)
	if len(vs) > 0 {
		return vs[0]
	}
	return ""
}

func (this *Header) Context(ctx context.Context) context.Context {
	return metadata.NewOutgoingContext(ctx, this.md)
}
