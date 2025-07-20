package loggroup

type Level string

const (
	LevelDebug Level = "DEBUG"
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
	LevelFatal Level = "FATAL"
)

type Status string

const (
	StatusUnresolved Status = "unresolved"
	StatusResolved   Status = "resolved"
	StatusIgnored    Status = "ignored"
)

type Group struct {
	ID          string `db:"id"`
	ProjectID   string `db:"project_id"`
	Level       Level  `db:"level"`
	Message     string `db:"message"`
	FirstSeenAt int64  `db:"first_seen_at"`
	LastSeenAt  int64  `db:"last_seen_at"`
	Counter     int    `db:"counter"`
	Status      Status `db:"status"`
}
