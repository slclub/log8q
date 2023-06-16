package log8q

import (
	"testing"
	"time"
)

func TestDefaultLog(t *testing.T) {
	Info("a", "bcd", "efg")
	Debug("a", "bcd", "efg")
	Warn("a", "bcd", "efg")

	time.Sleep(time.Millisecond * 10)
}
