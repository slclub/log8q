package log8q

import "context"

var default_log *Log8

func init() {
	default_log = New(context.Background(), &Config{
		Filename: "./log8q.log",
		Depth:    1,
	})
}

func Info(args ...any) {
	default_log.Info(args...)
}

func Debug(args ...any) {
	default_log.Debug(args...)
}

func Warn(args ...any) {
	default_log.Warn(args...)
}

func Error(args ...any) {
	default_log.Error(args...)
}

func Fatal(args ...any) {
	default_log.Fatal(args...)
}

func Infof(format string, args ...any) {
	default_log.Infof(format, args...)
}

func Debugf(format string, args ...any) {
	default_log.Debugf(format, args...)
}

func Warnf(format string, args ...any) {
	default_log.Warnf(format, args...)
}

func Errorf(format string, args ...any) {
	default_log.Errorf(format, args...)
}

func Fatalf(format string, args ...any) {
	default_log.Fatalf(format, args...)
}
