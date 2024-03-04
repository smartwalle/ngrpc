package ngrpc

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

func HeaderFromMetadata(md metadata.MD) *Header {
	if md == nil {
		md = metadata.MD{}
	}
	var h = &Header{}
	h.md = md
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

func (h *Header) Raw() metadata.MD {
	return h.md
}

func (h *Header) Set(key, value string) {
	h.md.Set(key, value)
}

func (h *Header) Get(key string) string {
	var vs = h.md.Get(key)
	if len(vs) > 0 {
		return vs[0]
	}
	return ""
}

func (h *Header) Len() int {
	return h.md.Len()
}

func (h *Header) Del(key string) {
	delete(h.md, strings.ToLower(key))
}

func (h *Header) ForeachKey(handler func(key, val string) error) error {
	for key, values := range h.md {
		for _, value := range values {
			if err := handler(key, value); err != nil {
				return err
			}
		}
	}
	return nil
}

func (h *Header) Context(ctx context.Context) context.Context {
	return metadata.NewOutgoingContext(ctx, h.md)
}
