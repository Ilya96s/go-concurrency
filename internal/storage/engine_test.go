package storage

import "testing"

func TestMemoryEngine_SetAndGet(t *testing.T) {
	engine := NewMemoryEngine()
	engine.Set("foo", "bar")

	val, ok := engine.Get("foo")
	if !ok {
		t.Fatal("expected key to exist")
	}
	if val != "bar" {
		t.Errorf("expected 'bar', got %v", val)
	}
}

func TestMemoryEngine_GetNotFound(t *testing.T) {
	engine := NewMemoryEngine()

	_, ok := engine.Get("missing")
	if ok {
		t.Error("expected missing key to not exist")
	}
}

func TestMemoryEngine_Del(t *testing.T) {
	engine := NewMemoryEngine()
	engine.Set("foo", "bar")

	ok := engine.Delete("foo")
	if !ok {
		t.Error("expected key to be deleted")
	}

	_, exists := engine.Get("foo")
	if exists {
		t.Error("expected key to be gone after delete")
	}
}

func TestMemoryEngine_DelNotFound(t *testing.T) {
	engine := NewMemoryEngine()

	ok := engine.Delete("missing")
	if ok {
		t.Error("expected false for deleting non-existent key")
	}
}
