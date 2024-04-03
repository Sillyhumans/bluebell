package mysql

import (
	"bluebell/models"
	"database/sql"
	"errors"

	"go.uber.org/zap"
)

func GetCommunityList() (communityList []*models.Community, err error) {
	sqlStr := "select community_id, community_name from community"
	err = db.Select(&communityList, sqlStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			zap.L().Warn("there is no community in db")
			err = nil
		}
	}
	return
}

func GetCommunityDetailByID(id int64) (community *models.CommunityDetail, err error) {
	community = new(models.CommunityDetail)
	sqlStr := "select community_id, community_name, introduction, create_time from community where community_id=?"
	err = db.Get(community, sqlStr, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = errors.New("there is invalid id")
		}
	}
	return
}
