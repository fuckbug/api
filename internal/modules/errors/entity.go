package errors

type Error struct {
	ID          string  `db:"id"`
	ProjectID   string  `db:"project_id"`
	Fingerprint string  `db:"fingerprint"`
	Message     string  `db:"message"`
	Stacktrace  string  `db:"stacktrace"`
	File        string  `db:"file"`
	Line        int     `db:"line"`
	Context     *string `db:"context"`
	IP          *string `db:"ip"`
	URL         *string `db:"url"`
	Method      *string `db:"method"`
	Headers     *string `db:"headers"`
	QueryParams *string `db:"query_params"`
	BodyParams  *string `db:"body_params"`
	Cookies     *string `db:"cookies"`
	Session     *string `db:"session"`
	Files       *string `db:"files"`
	Env         *string `db:"env"`
	Time        int64   `db:"time"`
	CreatedAt   int64   `db:"created_at"`
	UpdatedAt   int64   `db:"updated_at"`
}
