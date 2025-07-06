package log

import (
	"encoding/json"
	"fmt"
)

type Level string

const (
	LevelDebug Level = "DEBUG"
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
)

type Log struct {
	ID          string       `db:"id"`
	ProjectID   string       `db:"project_id"`
	Fingerprint string       `db:"fingerprint"`
	Level       Level        `db:"level"`
	Message     string       `db:"message"`
	Context     *interface{} `db:"context"`
	Time        int64        `db:"time"`
	CreatedAt   int64        `db:"created_at"`
	UpdatedAt   int64        `db:"updated_at"`
}

func (l *Log) PrepareForDB() error {
	if l.Context != nil {
		jsonData, err := json.Marshal(*l.Context)
		if err != nil {
			return fmt.Errorf("failed to marshal context: %w", err)
		}

		var contextValue interface{} = string(jsonData)
		l.Context = &contextValue
	}
	return nil
}
