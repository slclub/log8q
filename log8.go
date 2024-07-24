package log8q

import (
	"context"
	"fmt"
	"github.com/slclub/log8q/filer"
	"io"
	"sync"
	"time"
)

const (
	TIME_FORMAT = "2006-01-02 15:04:05"
)

type Config struct {
	Filename          string
	Format            string // 时间格式
	Permmison         Level
	Depth             int // 设置深度
	Writer            io.Writer
	CacheBucketLength int
	RotateTime        int64 // 日志保留秒数 单位s
}
type Log8 struct {
	w       io.Writer
	cache   *Cache
	logPool sync.Pool
	Option  Config
	ctx     context.Context
	depth   int
}

func New(ctx context.Context, config *Config) *Log8 {
	if config == nil {
		config = &Config{}
	}
	if config.Writer == nil {
		config.Writer = filer.New(ctx, &filer.Config{FileName: config.Filename, RotateTime: config.RotateTime}, nil)
	}
	config.Init()
	lg := &Log8{
		w:      config.Writer,
		depth:  config.Depth,
		Option: *config,
		cache:  NewCache(10, config.CacheBucketLength),
		ctx:    ctx,
	}
	lg.Init()
	return lg
}

// -------------------------------------------------------
// class log8
// -------------------------------------------------------

func (self *Log8) Init() {
	self.logPool.New = func() any {
		return newLine(self, 0)
	}
	go self.deamon()
}
func (self *Log8) Output(o io.Writer) {
	self.w = o
}

func (self *Log8) deamon() {
	ticker := time.NewTicker(time.Millisecond * 2)
	defer ticker.Stop()

	if self.w == nil || self.cache == nil {
		return
	}
	for {
		select {
		case <-self.ctx.Done():
			return
		case <-ticker.C:
			self.readTo()
		}
	}
}

func (self *Log8) readTo() {
	data := make([]byte, 8192)
	if self.cache.ReadSize() == 0 {
		return
	}
	i := 1
	for i > 0 {
		ii, err := self.cache.Read(data)
		i = ii
		if err != nil {
			continue
		}
		ii, err = self.w.Write(data[:i])
	}
}

func (self *Log8) depthPrint(depth int, level Level, args ...any) {
	if len(args) == 0 {
		return
	}

	if !level.Check(self.Option.Permmison) {
		return
	}
	line := self.logPool.Get().(*logLine)
	line.SetLevel(level)
	line.Handle(depth)
	fmt.Fprint(line, args...)
	// read to
	line.ReadTo(self.cache)
	self.afterWrite()
	line.Reset()
	self.logPool.Put(line)
}

func (self *Log8) depthPrintf(depth int, level Level, format string, args ...any) {
	if len(args) == 0 {
		return
	}

	if !level.Check(self.Option.Permmison) {
		return
	}
	line := self.logPool.Get().(*logLine)
	line.SetLevel(level)
	line.Handle(depth)
	fmt.Fprintf(line, format, args...)
	// read to
	line.ReadTo(self.cache)
	self.afterWrite()
}

func (self *Log8) afterWrite() {
	// 0.8 cap  内容时自动写入文件
	if 10*self.cache.ReadSize() <= 8*self.cache.Cap() {
		return
	}
	self.readTo()
}

func (self *Log8) Depth(depth int) {
	self.depth = depth
}

// 直接打印类
func (self *Log8) Print(args ...any) {
	self.depthPrint(self.depth, LEVEL_INFO, args...)
}

func (self *Log8) Info(args ...any) {
	self.depthPrint(self.depth, LEVEL_INFO, args...)
}

func (self *Log8) Debug(args ...any) {
	self.depthPrint(self.depth, LEVEL_DEBUG, args...)
}

func (self *Log8) Warn(args ...any) {
	self.depthPrint(self.depth, LEVEL_WARNNING, args...)
}

func (self *Log8) Error(args ...any) {
	self.depthPrint(self.depth, LEVEL_ERROR, args...)
}

func (self *Log8) Fatal(args ...any) {
	self.depthPrint(self.depth, LEVEL_FATAL, args...)
}

// format 参数类
func (self *Log8) Printf(format string, args ...any) {
	self.depthPrintf(self.depth, LEVEL_INFO, format, args...)
}

func (self *Log8) Infof(format string, args ...any) {
	self.depthPrintf(self.depth, LEVEL_INFO, format, args...)
}

func (self *Log8) Debugf(format string, args ...any) {
	self.depthPrintf(self.depth, LEVEL_DEBUG, format, args...)
}

func (self *Log8) Warnf(format string, args ...any) {
	self.depthPrintf(self.depth, LEVEL_WARNNING, format, args...)
}

func (self *Log8) Errorf(format string, args ...any) {
	self.depthPrintf(self.depth, LEVEL_ERROR, format, args...)
}

func (self *Log8) Fatalf(format string, args ...any) {
	self.depthPrintf(self.depth, LEVEL_FATAL, format, args...)
}

// ---------------------------------------------------------------
// Config
// ---------------------------------------------------------------

func (self *Config) Init() {
	if self.Format == "" {
		self.Format = TIME_FORMAT
	}
	if self.Filename == "" {
		self.Filename = "./log8q.log"
	}
	if self.Permmison == 0 {
		self.Permmison = ALL_LEVEL
	}
	if self.CacheBucketLength == 0 {
		self.CacheBucketLength = LEN_BUCKET
	}
}
