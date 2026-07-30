package tsc

import (
	"testing"
	"time"
)

func TestCompileTimesAddNestedEmitTime(t *testing.T) {
	t.Parallel()

	test := func(checkTime time.Duration, nested time.Duration, wantCheck time.Duration, wantEmit time.Duration) {
		t.Helper()
		times := &CompileTimes{checkTime: checkTime, emitTime: 7 * time.Second}
		times.addNestedEmitTime(nested)
		if times.checkTime != wantCheck || times.emitTime != wantEmit {
			t.Errorf("times = {check: %v, emit: %v}, want {check: %v, emit: %v}", times.checkTime, times.emitTime, wantCheck, wantEmit)
		}
	}

	test(5*time.Second, 3*time.Second, 2*time.Second, 10*time.Second)
	test(5*time.Second, 8*time.Second, 0, 15*time.Second)
}
