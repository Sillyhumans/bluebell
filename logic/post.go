package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	mySnowflake "bluebell/pkg/snowflake"
	"errors"
	"strconv"
)

var (
	orderTime  = "time"
	orderScore = "score"
)

func CreatePost(p *models.Post) error {
	// 生成post_id
	p.ID = int64(mySnowflake.GetID())
	// 保存到数据库
	err := mysql.CreatePost(p)
	if err != nil {
		return err
	}
	// 计算分数 存入缓存
	err = SetPostScore(0, 0, p.ID)
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
func GetPostList(p *models.ParaPostList) (data []*models.ApiPostDetail, err error) {
	//从redis获取 post的id
	//如果是按分数排名
	var (
		ids   []string
		posts []*models.Post
	)
	if p.Order == orderScore {
		ids, err = redis.GetPostOrderScore(p)
		if err != nil {
			return
		}
		posts, err = mysql.GetPostListByIDs(ids)
		if err != nil {
			return
		}
	} else if p.Order == orderTime {
		posts, err = mysql.GetPostList(p.Page, p.Size)
		ids = make([]string, len(posts))
		for i := 0; i < len(posts); i++ {
			ids[i] = strconv.Itoa(int(posts[i].AuthorID))
		}
		if err != nil {
			return
		}
	} else {
		return nil, errors.New("invalid err")
	}

	l := len(posts)
	data = make([]*models.ApiPostDetail, 0, l)
	// 获取所有post对应user
	for i := 0; i < l; i++ {
		user, err := mysql.GetUserByID(posts[i].AuthorID)
		if err != nil {
			return nil, err
		}
		community, err := mysql.GetCommunityDetailByID(posts[i].CommunityID)
		if err != nil {
			return nil, err
		}
		postDetail := &models.ApiPostDetail{
			User:            user,
			Post:            posts[i],
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}
