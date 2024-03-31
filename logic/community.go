package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
)

func GetCommunityList() (data []*models.Community, err error) {
	// 查找到所有的community，并返回
	return mysql.GetCommunityList()
}

func GetCommunityDetailByID(id int64) (c *models.CommunityDetail, err error) {
	return mysql.GetCommunityDetailByID(id)
}
