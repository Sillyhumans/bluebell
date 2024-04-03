package models

type Vote struct {
	PostID int64 `json:"post_id,string" db:"post_id"`
	UserID int64 `json:"user_id,string" db:"user_id"`
	Vote   int8  `json:"vote" db:"vote"`
}
