package compute

import (
	"errors"
	"strings"
)

// Необходим для правильной обработки команд пользователя.
// Часть compute слоя, которая переводит строку "SET key value" в структуру {Type: SET, Arhs: ["Key", "Value"]}

type CommandType int

// Возможные типы команд
const (
	Unknown CommandType = iota
	Set
	Get
	Del
)

// Структура разобранной команды
type Command struct {
	Type CommandType
	Args []string
}

// Парсит строку запроса в Command
func Parse(input string) (*Command, error) {
	fields := strings.Fields(input)
	if len(fields) == 0 {
		return nil, errors.New("invalid or empty command")
	}

	switch fields[0] {
	case "SET":
		if len(fields) != 3 {
			return nil, errors.New("usage: SET key value")
		}
		return &Command{Type: Set, Args: fields[1:]}, nil
	case "GET":
		if len(fields) != 2 {
			return nil, errors.New("usage: GET key")
		}
		return &Command{Type: Get, Args: fields[1:]}, nil
	case "DEL":
		if len(fields) != 2 {
			return nil, errors.New("usage: DEL key")
		}
		return &Command{Type: Del, Args: fields[1:]}, nil
	default:
		return nil, errors.New("unknown command")
	}
}
