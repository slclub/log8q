package log8q

import (
	"context"
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
