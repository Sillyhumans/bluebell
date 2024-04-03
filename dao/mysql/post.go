package mysql

import (
	"bluebell/models"
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
	sqlStr := `select post_id, title, content, author_id, community_id, status, create_time from post limit ?,?`
	posts = make([]*models.Post, 0, 2)
	err = db.Select(&posts, sqlStr, (offset-1)*limit, limit)
	return
}
