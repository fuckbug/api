package errors

type Error struct {
	ID          string `db:"id"`
	ProjectID   string `db:"project_id"`
	Fingerprint string `db:"fingerprint"`
	Message     string `db:"message"`
	Stacktrace  string `db:"stacktrace"`
	File        string `db:"file"`
	Line        int    `db:"line"`
	Context     string `db:"context"`
	Time        int64  `db:"time"`
	CreatedAt   int64  `db:"created_at"`
	UpdatedAt   int64  `db:"updated_at"`
}
