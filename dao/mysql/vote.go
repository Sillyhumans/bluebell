package mysql

import (
	"bluebell/models"
)

func QueryVote(p *models.VoteData) (vote *models.Vote, err error) {
	vote = new(models.Vote)
	sqlStr := "select post_id, user_id, vote from vote where user_id=? and post_id=?"
	err = db.Get(vote, sqlStr, p.UserID, p.PostID)
	return
}

func CreateVote(p *models.VoteData) (err error) {
	sqlStr := "insert into vote(post_id, user_id, vote) values(?, ?, ?)"
	_, err = db.Exec(sqlStr, p.PostID, p.UserID, p.Vote)
	return
}

func UpdateVote(p *models.VoteData) (err error) {
	sqlStr := "update vote set vote = vote + ? where post_id = ? and user_id = ?"
	_, err = db.Exec(sqlStr, p.Vote, p.PostID, p.UserID)
	return
}

func GetCountUp(postID int64) (up int64, err error) {
	sqlStr := "select count(*) from vote where vote=1 and post_id = ?"
	row, err := db.Query(sqlStr, postID)
	for row.Next() {
		err = row.Scan(&up)
		if err != nil {
			return
		}
	}
	return
}

func GetCountDown(postID int64) (up int64, err error) {
	sqlStr := "select count(*) from vote where vote=-1 and post_id = ?"
	row, err := db.Query(sqlStr, postID)
	for row.Next() {
		err = row.Scan(&up)
		if err != nil {
			return
		}
	}
	return
}
