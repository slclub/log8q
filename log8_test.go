package log8q

import (
	"context"
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
	})
	log.Print("stdin print")
	log.Debug("stdin debug", "b", "c", "d")
	log.Warn("stdin warn Get File line")
	time.Sleep(time.Millisecond * 10)
}
