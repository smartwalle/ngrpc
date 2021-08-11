package logging

type Logger interface {
	SetPrefix(prefix string)
	Prefix() string
	Println(args ...interface{})
	Printf(format string, args ...interface{})
	Output(calldepth int, s string) error
}

type nilLogger struct {
}

func (this *nilLogger) SetPrefix(prefix string) {
}

func (this *nilLogger) Prefix() string {
	return ""
}

func (this *nilLogger) Println(args ...interface{}) {
}

func (this *nilLogger) Printf(format string, args ...interface{}) {
}

func (this *nilLogger) Output(calldepth int, s string) error {
	return nil
}
