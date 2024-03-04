package logging

import "context"

type Logger interface {
	Println(ctx context.Context, args ...interface{})
	Printf(ctx context.Context, format string, args ...interface{})
}

type nilLogger struct {
}

func (logger *nilLogger) Println(ctx context.Context, args ...interface{}) {
}

func (logger *nilLogger) Printf(ctx context.Context, format string, args ...interface{}) {
}
