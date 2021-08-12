package grpc4go

import (
	"context"
	"google.golang.org/grpc/metadata"
	"strings"
)

type Header struct {
	md metadata.MD
}

func NewHeader() *Header {
	var h = &Header{}
	h.md = metadata.MD{}
	return h
}

func HeaderFromIncoming(ctx context.Context) *Header {
	var h = &Header{}
	h.md, _ = metadata.FromIncomingContext(ctx)
	if h.md == nil {
		h.md = metadata.MD{}
	}
	return h
}

func HeaderFromOutgoing(ctx context.Context) *Header {
	var h = &Header{}
	h.md, _ = metadata.FromOutgoingContext(ctx)
	if h.md == nil {
		h.md = metadata.MD{}
	}
	return h
}

func (this *Header) Set(key, value string) {
	this.md.Set(key, value)
}

func (this *Header) Get(key string) string {
	var vs = this.md.Get(key)
	if len(vs) > 0 {
		return vs[0]
	}
	return ""
}

func (this *Header) Len() int {
	return this.md.Len()
}

func (this *Header) Del(key string) {
	delete(this.md, strings.ToLower(key))
}

func (this *Header) Context(ctx context.Context) context.Context {
	return metadata.NewOutgoingContext(ctx, this.md)
}
