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

type GroupsStatus string

const (
	GroupsStatusUnresolved GroupsStatus = "unresolved"
	GroupsStatusResolved   GroupsStatus = "resolved"
	GroupsStatusIgnored    GroupsStatus = "ignored"
)

type Group struct {
	Fingerprint string       `db:"fingerprint"`
	ProjectID   string       `db:"project_id"`
	Type        string       `db:"type"`
	File        string       `db:"file"`
	Line        int          `db:"line"`
	Message     string       `db:"message"`
	FirstSeenAt int64        `db:"first_seen_at"`
	LastSeenAt  int64        `db:"last_seen_at"`
	Counter     int64        `db:"counter"`
	Status      GroupsStatus `db:"status"`
}
