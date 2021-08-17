package ratelimit

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Limiter interface {
	Allow() bool
}

func errorFrom(opt *option, method string) error {
	if opt.handler != nil {
		return opt.handler(method)
	}
	return status.Errorf(codes.ResourceExhausted, "%s is rejected by rate limit middleware", method)
}
