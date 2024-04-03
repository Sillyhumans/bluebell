package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	mySnowflake "bluebell/pkg/snowflake"
)

func CreatePost(p *models.Post) error {
	// 生成post_id
	p.ID = int64(mySnowflake.GetID())
	// 保存到数据库
	err := mysql.CreatePost(p)
	return err
}

func GetPostByID(pid int64) (data *models.ApiPostDetail, err error) {
	// 查询并组合接口想用的数据
	p, err := mysql.GetPostByID(pid)
	if err != nil {
		return
	}
	// 根据作者id获取作者信息
	user, err := mysql.GetUserByID(p.AuthorID)
	if err != nil {
		return
	}
	// 根据ID获取社区信息
	community, err := mysql.GetCommunityDetailByID(p.CommunityID)
	if err != nil {
		return
	}
	data = new(models.ApiPostDetail)
	data.User = user
	data.Post = p
	data.CommunityDetail = community
	return
}

// GetPostList 获取所有post以及每个post的作者，社区信息
func GetPostList(offset, limit int64) (data []*models.ApiPostDetail, errs error) {
	// 获取所有post
	posts, errs := mysql.GetPostList(offset, limit)
	if errs != nil {
		return
	}
	l := len(posts)
	data = make([]*models.ApiPostDetail, 0, l)

	// 获取所有post对应user
	for i := 0; i < l; i++ {
		user, err := mysql.GetUserByID(posts[i].AuthorID)
		if err != nil {
			errs = err
			return
		}
		community, err := mysql.GetCommunityDetailByID(posts[i].CommunityID)
		if err != nil {
			errs = err
			return
		}
		postDetail := &models.ApiPostDetail{
			user,
			posts[i],
			community,
		}
		data = append(data, postDetail)
	}
	return
}
