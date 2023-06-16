package log8q

import (
	"errors"
	"fmt"
	"github.com/slclub/go-tips/stringbyte"
	"io"
	"math"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

var (
	SEP = []byte(" ")
)

type Trace struct {
	filename   []byte
	line       int
	stack_info []byte
}

type logLine struct {
	log8   *Log8
	data   []byte
	level  Level
	trace  *Trace
	lenght int
}

func newLine(log8 *Log8, level Level) *logLine {
	line := &logLine{
		log8:  log8,
		level: level,
		data:  make([]byte, 0),
		trace: &Trace{
			filename: make([]byte, 0),
		},
	}
	line.SetLevel(level)
	return line
}

// ------------------------------------------------
// class logLine
func (self *logLine) Reset() {
	self.data = nil
	self.level = 0
	self.lenght = 0
}

func (self *logLine) SetLevel(level Level) {
	if self.level == 0 {
		self.lenght += level.Len()
	}
	self.level = level

	if (level).Trace() {
		self.trace.filename = make([]byte, 0)
	}
}

func (self *logLine) Handle(depth int) {
	if self.level >= LEVEL_WARNNING {
		self.handleHeadLine(depth)
	}
}

func (self *logLine) handleHeadLine(depth int) {
	_, file, line, ok := runtime.Caller(5 + depth)
	if !ok {
		self.trace.filename = []byte("???")
		self.trace.line = -1
		return
	}
	slash := strings.LastIndex(file, "/")
	if slash >= 0 {
		self.trace.filename = stringbyte.StringToBytes(file[slash+1:])
	}

	self.trace.line = line
	self.lenght += len(self.trace.filename)
	self.lenght += 8
}

func (self *logLine) handleStack() {
	self.trace.stack_info = stack(-1)
	self.lenght += len(self.trace.stack_info)
}

func (self *logLine) handleTime() []byte {
	now := time.Now()
	ts := []byte(now.Format(self.log8.Option.Format))
	self.lenght += len(ts)
	return ts
}

func (self *logLine) ReadTo(o io.Writer) (int, error) {
	if o == nil {
		return 0, errors.New("log8q.WriteWith is nil")
	}
	ww := [][]byte{}
	ww = append(ww, self.handleTime(), SEP)
	ww = append(ww, stringbyte.StringToBytes(self.level.String()), SEP)
	ww = append(ww, self.data)
	ww = append(ww, SEP)
	var (
		n   int
		err error
	)
	self.lenght += 3
	defer func() {
		self.writeDataTo(ww, o)
	}()
	if self.trace == nil {
		return n, err
	}
	if self.level >= LEVEL_WARNNING {
		ww = append(ww, self.trace.filename, SEP)
		self.lenght += 1
		lb, err := IntToBytes(self.trace.line)
		if err != nil {
			ww = append(ww, lb, SEP)
		}
	}
	if self.level >= TRACE_INFO {
		self.handleStack()
		ww = append(ww, self.trace.stack_info)
	}

	ww = append(ww, lineBreak())
	self.lenght += 1
	return n, err
}

func (self *logLine) writeDataTo(ww [][]byte, o io.Writer) (int, error) {
	if wm, ok := o.(WriteMany); ok {
		n, err := wm.WriteMany(ww...)
		return n, err
	}
	data := make([]byte, self.lenght)
	i := 0
	for _, v := range ww {
		n := copy(data[i:], v)
		i += n
	}
	return o.Write(data)
}

func (self *logLine) Read(b []byte) (int, error) {
	return 0, nil
}

func (self *logLine) Write(p []byte) (int, error) {
	n := len(p)
	self.data = append(self.data, p...)
	self.lenght += n
	return n, nil
}

// ------------------------------------------------
// functions
func stack(size int) []byte {
	//var trace = make([]byte, 1024, size)
	//runtime.Stack(trace, true)
	trace := debug.Stack()
	s := trace[600:]
	if size == -1 {
		return s
	}
	if size >= len(s) {
		return s
	}

	return s[:size]
}

func lineBreak() []byte {
	if IsUnix() {
		return []byte{'\n'}
	}
	return []byte{'\r', '\n'}
}

// Judging unix system.
// Unix*
func IsUnix() bool {
	sys_type := runtime.GOOS

	if sys_type == "windows" {
		return false
	}
	return true
}

func IntToBytes(a int) ([]byte, error) {
	if a > math.MaxInt32 {
		return nil, errors.New(fmt.Sprintf("a>math.MaxInt32, a is %d\n", a))
	}
	buf := make([]byte, 4)
	for i := 0; i < 4; i++ {
		var b uint8 = uint8(a & 0xff)
		buf[i] = b
		a = a >> 8
	}
	return buf, nil
}
