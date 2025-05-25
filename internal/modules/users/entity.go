package users

type User struct {
	ID        string `db:"id"`
	Email     string `db:"email"`
	Password  string `db:"password"`
	Role      string `db:"role"`
	CreatedAt int64  `db:"created_at"`
	UpdatedAt int64  `db:"updated_at"`
}
