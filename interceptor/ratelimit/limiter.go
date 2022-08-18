package ratelimit

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Limiter interface {
	Allow() bool
}

func errorFrom(opts *options, method string) error {
	if opts.handler != nil {
		return opts.handler(method)
	}
	return status.Errorf(codes.ResourceExhausted, "%s is rejected by rate limit middleware", method)
}
