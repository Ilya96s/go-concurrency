package compute

import "in-memory-kv/storage"

// Связывает parser и storage
// Compute — связывает parser и storage
type Compute struct {
	engine storage.Engine
}

// NewCompute создает новый обработчик
func NewCompute(engine storage.Engine) *Compute {
	return &Compute{engine: engine}
}

// Handle обрабатывает строку и возвращает результат
func (c *Compute) Handle(line string) string {
	cmd, err := Parse(line)
	if err != nil {
		return "ERROR: " + err.Error()
	}

	switch cmd.Type {
	case Set:
		c.engine.Set(cmd.Args[0], cmd.Args[1])
		return "OK"
	case Get:
		if val, ok := c.engine.Get(cmd.Args[0]); ok {
			return val
		}
		return "NOT_FOUND"
	case Del:
		if c.engine.Delete(cmd.Args[0]) {
			return "OK"
		}
		return "NOT_FOUND"
	default:
		return "ERROR: unknown command"
	}

}
