package project

type Project struct {
	ID        string `db:"id"`
	Name      string `db:"name"`
	PublicKey string `db:"public_key"`
	CreatedAt int64  `db:"created_at"`
	UpdatedAt int64  `db:"updated_at"`
	DeletedAt *int64 `db:"deleted_at"`
}
