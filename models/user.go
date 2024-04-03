package models

import "database/sql"

type User struct {
	UserID   int64          `db:"user_id"`
	Gender   int32          `db:"gender"`
	Email    sql.NullString `db:"email"`
	UserName string         `db:"username"`
	Password string         `db:"password"`
}
