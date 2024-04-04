package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/vote"
	"database/sql"
	"errors"
	"strconv"
)

// PostVote 为帖子投票
func PostVote(p *models.VoteData) error {
	var v int8
	key := strconv.FormatInt(p.PostID, 10) + strconv.FormatInt(p.UserID, 10)
	// 向redis中查找是否有值
	v, err := redis.GetVote(key)
	if err != nil {
		// 如果没有值 查找数据库
		voteObj, err1 := mysql.QueryVote(p)
		// 用户没投过票 创建
		if errors.Is(err1, sql.ErrNoRows) {
			err2 := mysql.CreateVote(p)
			if err2 != nil {
				return err2
			}
			// 写入缓存
			err2 = redis.SetVote(key, p.Vote)
			if err2 != nil {
				return err2
			}
			// 更新帖子分数
			up, down, err2 := vote.GetUpAndDown(p.PostID)
			if err2 != nil {
				return err2
			}
			err2 = SetPostScore(up, down, p.PostID)
			if err2 != nil {
				return err2
			}
			return err2
		} else if err1 != nil {
			return err1
		}
		v = voteObj.Vote
	}
	// 用户投过票 则进行对比
	// 相同则不进行更新 否则更新并删除redis缓存
	if p.Vote == v {
		return nil
	} else {
		// 更新数据库
		err = mysql.UpdateVote(p)
		if err != nil {
			return err
		}
		// 更新帖子分数
		up, down, err := vote.GetUpAndDown(p.PostID)
		if err != nil {
			return err
		}
		err = SetPostScore(up, down, p.PostID)
		if err != nil {
			return err
		}
		// 删除缓存, 如果缓存不存在就更新缓存
		err = redis.DelSetVote(key, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func SetPostScore(up, down int64, postID int64) error {
	t, err := mysql.GetPostCreateTime(postID)
	if err != nil {
		return err
	}
	score := vote.Hot(up, down, t)
	err = redis.SetPostScore(postID, score)
	return err
}
