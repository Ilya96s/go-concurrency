package storage

// Интерфейс движка для хранения данных
type Engine interface {
	Set(key, value string)

	Get(key string) (string, bool)

	Delete(key string) bool
}
