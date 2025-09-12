package storage

type MemoryEngine struct {
	Data map[string]string
}

func NewMemoryEngine() *MemoryEngine {
	return &MemoryEngine{
		Data: make(map[string]string),
	}
}

func (memoryEngine *MemoryEngine) Set(key string, value string) {
	memoryEngine.Data[key] = value
}

func (memoryEngine *MemoryEngine) Get(key string) (string, bool) {
	value, ok := memoryEngine.Data[key]
	return value, ok
}

func (memoryEngine *MemoryEngine) Delete(key string) bool {
	_, ok := memoryEngine.Data[key]
	if ok {
		delete(memoryEngine.Data, key)
		return ok
	}
	return ok
}
