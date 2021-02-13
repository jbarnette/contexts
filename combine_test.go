package contexts

import (
	"context"
	"testing"
	"time"
)

func TestCombineDeadline(t *testing.T) {
	d1 := time.Now().Add(5 * time.Minute)
	d2 := time.Now().Add(1 * time.Minute)

	c1, cancel := context.WithDeadline(context.Background(), d1)
	defer cancel()

	c2, cancel := context.WithDeadline(context.Background(), d2)
	defer cancel()

	combined := Combine(c1, c2)
	dm, ok := combined.Deadline()

	if !ok {
		t.Error("missing deadline")
	}

	if dm != d2 {
		t.Errorf("deadline %v != %v", dm, d2)
	}
}

func TestCombineValue(t *testing.T) {
	type key string

	c1 := context.WithValue(context.Background(), key("c1"), "c1-value")
	c2 := context.WithValue(context.Background(), key("c1"), "c2-c1-value")
	c2 = context.WithValue(c2, key("c2"), "c2-value")

	ctx := Combine(c1, c2)

	if v, ok := ctx.Value(key("c1")).(string); !ok || v != "c1-value" {
		t.Errorf("c1 value %q != %q", v, "c2-c1-value")
	}

	if v, ok := ctx.Value(key("c2")).(string); !ok || v != "c2-value" {
		t.Errorf("c2 value %q != %q", v, "c2-value")
	}

	if v, ok := ctx.Value(key("nonexistent")).(string); ok {
		t.Errorf("nonexistent value %q", v)
	}
}

func TestCombineDoneErr(t *testing.T) {
	c1, cancel := context.WithCancel(context.Background())
	defer cancel()

	c2, cancel := context.WithCancel(context.Background())
	ctx := Combine(c1, c2)

	if err := ctx.Err(); err != nil {
		t.Errorf("unexpected Err: %v", err)
	}

	cancel()
	<-ctx.Done()

	if err := ctx.Err(); err != context.Canceled {
		t.Errorf("unexpected Err: %v != context.Canceled", err)
	}
}
