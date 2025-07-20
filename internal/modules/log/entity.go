package log

type Level string

const (
	LevelDebug Level = "DEBUG"
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
	LevelFatal Level = "FATAL"
)

type Log struct {
	ID          string  `db:"id"`
	ProjectID   string  `db:"project_id"`
	Fingerprint string  `db:"fingerprint"`
	Level       Level   `db:"level"`
	Message     string  `db:"message"`
	Context     *string `db:"context"`
	Time        int64   `db:"time"`
	CreatedAt   int64   `db:"created_at"`
	UpdatedAt   int64   `db:"updated_at"`
}
