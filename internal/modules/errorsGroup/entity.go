package errorsgroup

type Status string

const (
	StatusUnresolved Status = "unresolved"
	StatusResolved   Status = "resolved"
	StatusIgnored    Status = "ignored"
)

type Group struct {
	ID          string `db:"id"`
	ProjectID   string `db:"project_id"`
	File        string `db:"file"`
	Line        int    `db:"line"`
	Message     string `db:"message"`
	FirstSeenAt int    `db:"first_seen_at"`
	LastSeenAt  int    `db:"last_seen_at"`
	Counter     int    `db:"counter"`
	Status      Status `db:"status"`
}
