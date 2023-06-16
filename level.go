package log8q

import "strings"

type Level int

const (
	// LEVEL
	LEVEL_INFO     Level = 1
	LEVEL_DEBUG    Level = 2
	LEVEL_WARNNING Level = 4
	LEVEL_ERROR    Level = 8
	LEVEL_FATAL    Level = 16

	TRACE_INFO     Level = 32
	TRACE_DEBUG    Level = 64
	TRACE_WARNNING Level = 128
	TRACE_ERROR    Level = 256
	TRACE_FATAL    Level = 512
)
const (
	CONST_QUEUE_CAPACITY = 2 * 1024 * 1024
)

var level_name = map[Level]string{
	LEVEL_INFO:     "INFO ",
	LEVEL_DEBUG:    "DEBUG",
	LEVEL_WARNNING: "WARN ",
	LEVEL_ERROR:    "ERROR",
	LEVEL_FATAL:    "FATAL",
}

var (
	// default level. also is most commonly used.
	ALL_LEVEL = LEVEL_INFO + LEVEL_DEBUG + LEVEL_WARNNING + LEVEL_ERROR + LEVEL_FATAL + TRACE_ERROR + TRACE_FATAL
	// debuging. = 1024 -1
	ALL_TRACE = ALL_LEVEL + TRACE_INFO + TRACE_DEBUG + TRACE_WARNNING
)

func (level Level) String() string {
	in := level.Int()
	if in < 32 {
		return level_name[level]
	}
	n := in / 32
	n = 1<<n - 1
	return "TRACE " + Level(n).String()
}

func (level Level) Int() int {
	return int(level)
}

func (level Level) Trim() string {
	return strings.Trim(level.String(), " ")
}

func (level Level) Len() int {
	return len(level.String())
}

func (level Level) Trace() bool {
	if level < LEVEL_ERROR {
		return false
	}
	return true
}

func (level Level) Check(l Level) bool {
	if level&l > 0 {
		return true
	}
	return false
}

// functions
func ConvLevel(v any) Level {
	switch val := v.(type) {
	case int:
		return Level(val)
	case uint:
		return Level(val)
	case int64:
		return Level(val)
	case uint64:
		return Level(val)
	case int32:
		return Level(val)
	case uint32:
		return Level(val)
	case int8:
		return Level(val)
	case uint8:
		return Level(val)
	case int16:
		return Level(val)
	case uint16:
		return Level(val)
	}
	return 0
}
