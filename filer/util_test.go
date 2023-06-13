package filer

import (
	"testing"
	"time"
)

func TestUtil(t *testing.T) {
	if ok, err := isFileExist("/tmp/nothing/nothing"); ok {
		t.Error("Log.Util", err)
	}

	now := time.Now()
	if !isSameDate(now, now) {
		t.Error("Log.Util.isSameDate ")
	}
	if isSameDate(now, now.Add(86400*time.Second)) {
		t.Error("Log.Util.isSameDate")
	}

	d := getFileModifyTime("/tmp/nothing/nothing")
	d1 := d.Add(-1 * time.Hour)
	if d1.Unix() > 0 {
		t.Error("Log.Util.getFileModifyTime")
	}
}
