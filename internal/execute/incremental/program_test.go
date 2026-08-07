package incremental

import (
	"testing"
	"time"
)

func TestNestedEmitTime(t *testing.T) {
	t.Parallel()

	now := time.Unix(0, 0)
	program := &Program{
		nestedEmitNow: func() time.Time {
			return now
		},
	}

	doneOuter := program.beginNestedEmit()
	now = now.Add(time.Second)
	doneInner := program.beginNestedEmit()
	now = now.Add(2 * time.Second)
	doneInner()
	now = now.Add(3 * time.Second)
	doneOuter()

	if duration := program.TakeNestedEmitTime(); duration != 6*time.Second {
		t.Fatalf("nested emit time = %v, want %v", duration, 6*time.Second)
	}
	if duration := program.TakeNestedEmitTime(); duration != 0 {
		t.Fatalf("nested emit time after take = %v, want 0", duration)
	}
}
