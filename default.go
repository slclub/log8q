package log8q

import "context"

var default_log *Log8

func defaultLog() *Log8 {
	if default_log != nil {
		return default_log
	}
	default_log = New(context.Background(), &Config{
		Filename: "logs/log8q.log",
		Depth:    1,
	})
	return default_log
}

func Info(args ...any) {
	defaultLog().Info(args...)
}

func Debug(args ...any) {
	defaultLog().Debug(args...)
}

func Warn(args ...any) {
	defaultLog().Warn(args...)
}

func Error(args ...any) {
	defaultLog().Error(args...)
}

func Fatal(args ...any) {
	defaultLog().Fatal(args...)
}

func Infof(format string, args ...any) {
	defaultLog().Infof(format, args...)
}

func Debugf(format string, args ...any) {
	defaultLog().Debugf(format, args...)
}

func Warnf(format string, args ...any) {
	defaultLog().Warnf(format, args...)
}

func Errorf(format string, args ...any) {
	defaultLog().Errorf(format, args...)
}

func Fatalf(format string, args ...any) {
	defaultLog().Fatalf(format, args...)
}
