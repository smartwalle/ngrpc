package logging

import "context"

type Logger interface {
	Println(ctx context.Context, args ...interface{})
	Printf(ctx context.Context, format string, args ...interface{})
}

type nilLogger struct {
}

func (this *nilLogger) Println(ctx context.Context, args ...interface{}) {
}

func (this *nilLogger) Printf(ctx context.Context, format string, args ...interface{}) {
}
