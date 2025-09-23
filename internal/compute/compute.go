package compute

import (
	"go.uber.org/zap"
	"in-memory-kv/internal/storage"
)

// Compute — связывает parser и storage
type Compute struct {
	engine storage.Engine
	log    *zap.Logger
}

// NewCompute создает новый обработчик
func NewCompute(engine storage.Engine, log *zap.Logger) *Compute {
	return &Compute{engine: engine,
		log: log}
}

// Handle обрабатывает строку и возвращает результат
func (c *Compute) Handle(line string) string {
	cmd, err := Parse(line)
	if err != nil {
		c.log.Warn("parse error", zap.String("input", line), zap.Error(err))
		return "ERROR: " + err.Error()
	}

	switch cmd.Type {
	case Set:
		c.engine.Set(cmd.Args[0], cmd.Args[1])
		c.log.Info("SET executed",
			zap.String("key", cmd.Args[0]),
			zap.String("value", cmd.Args[1]))
		return "OK"

	case Get:
		if val, ok := c.engine.Get(cmd.Args[0]); ok {
			c.log.Info("GET found",
				zap.String("key", cmd.Args[0]),
				zap.String("value", val))
			return val
		}
		c.log.Info("GET not found",
			zap.String("key", cmd.Args[0]))
		return "NOT_FOUND"

	case Del:
		if c.engine.Delete(cmd.Args[0]) {
			c.log.Info("DEL success",
				zap.String("key", cmd.Args[0]))
			return "OK"
		}
		c.log.Info("DEL not found",
			zap.String("key", cmd.Args[0]))
		return "NOT_FOUND"

	default:
		c.log.Warn("unknown command",
			zap.String("input", line))
		return "ERROR: unknown command"
	}

}
