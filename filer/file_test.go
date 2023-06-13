package filer

import (
	"context"
	"testing"
)

func TestFileLog(t *testing.T) {
	filename := "/tmp/log8q/log8q.log"
	ctx, cancel := context.WithCancel(context.Background())
	w := New(ctx, &Config{
		FileName: filename,
	}, nil)
	_, err := w.Write([]byte("你好啊"))

	if err != nil {
		t.Error("Log.File.Write ", err)
	}

	if ok, _ := isFileExist(filename); !ok {
		t.Error("Log.File is not exist")
	}
	cancel()
}

func TestConfig(t *testing.T) {
	filename := "/tmp/log8q/log8q.log"
	config := &Config{
		FileName: filename,
	}
	config.Init()
	if config.BaseName() != "log8q" {
		t.Error("Log.Config.BaseName", config.BaseName())
	}
}
