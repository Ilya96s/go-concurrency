package storage

import "sync"

type MemoryEngine struct {
	mu   sync.RWMutex
	Data map[string]string
}

func NewMemoryEngine() *MemoryEngine {
	return &MemoryEngine{
		Data: make(map[string]string),
	}
}

func (memoryEngine *MemoryEngine) Set(key string, value string) {
	memoryEngine.mu.Lock()
	defer memoryEngine.mu.Unlock()
	memoryEngine.Data[key] = value
}

func (memoryEngine *MemoryEngine) Get(key string) (string, bool) {
	memoryEngine.mu.RLock()
	defer memoryEngine.mu.RUnlock()
	value, ok := memoryEngine.Data[key]
	return value, ok
}

func (memoryEngine *MemoryEngine) Delete(key string) bool {
	memoryEngine.mu.Lock()
	defer memoryEngine.mu.Unlock()
	_, ok := memoryEngine.Data[key]
	if ok {
		delete(memoryEngine.Data, key)
		return ok
	}
	return ok
}
