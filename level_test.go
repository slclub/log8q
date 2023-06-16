package log8q

import "testing"

func TestLevel(t *testing.T) {
	info := LEVEL_INFO.String()
	if info != "INFO " {
		t.Error("Level.info string", LEVEL_INFO.Int(), LEVEL_INFO.String())
	}

	if TRACE_INFO.String() != "TRACE INFO " {
		t.Error("Level.trace_info string", TRACE_INFO.Int(), TRACE_INFO.String())
	}
}
