package mysql

import (
	"bluebell/models"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

func CreatePost(p *models.Post) (err error) {
	sqlStr := "insert into post(post_id, title, content, author_id, community_id) values (?, ?, ?, ?, ?)"
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

func GetPostByID(pid int64) (data *models.Post, err error) {
	data = new(models.Post)
	sqlStr := `select post_id, title, content, author_id, community_id, status, create_time
				from post where post_id = ?`
	err = db.Get(data, sqlStr, pid)
	return
}

func GetPostList(offset, limit int64) (posts []*models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, status, create_time from post order by create_time ASC limit ?,?`
	posts = make([]*models.Post, 0, 2)
	err = db.Select(&posts, sqlStr, (offset-1)*limit, limit)
	return
}

func GetPostCreateTime(postID int64) (t time.Time, err error) {
	sqlStr := "select create_time from post where post_id = ?"
	row, err := db.Query(sqlStr, postID)
	for row.Next() {
		err = row.Scan(&t)
		if err != nil {
			return
		}
	}
	return
}

func GetPostListByIDs(ids []string) (posts []*models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time 
				from post where id in (?) order by FIND_IN_SET(post_id, ?)`
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	query = db.Rebind(query)
	err = db.Select(&posts, query, args)
	return
}
