package log8q

import (
	"context"
	"github.com/slclub/log8q/filer"
	"os"
	"testing"
	"time"
)

func TestLog8(t *testing.T) {
	log := New(context.Background(), nil)
	log.Print("a", "b", "c")
	log.Debug("a", "b", "c", "d")
	log.Warn("Get File line")
	time.Sleep(time.Millisecond * 10)
}

func TestLog8NewIo(t *testing.T) {
	log := New(context.Background(), &Config{
		Writer: os.Stdout,
		//Permmison: ALL_TRACE, // 设置 trace 才有 TraceWarn 和TraceError 等才能输出相关 内容 模式
	})
	log.Print("stdin print")
	log.Debug("stdin debug", "b", "c", "d")
	log.Warn("stdin warn Get File line")
	time.Sleep(time.Millisecond * 10)
}

func TestLogPath(t *testing.T) {
	conf := &Config{
		Filename: "logs/log8q.log",
	}
	file := filer.New(context.Background(), &filer.Config{FileName: conf.Filename, RotateTime: conf.RotateTime}, nil)
	conf.Writer = file
	log8 := New(context.Background(), conf)
	log8.Infof("TestLogPath ")
}
