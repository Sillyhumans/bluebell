package mysql

import (
	"bluebell/models"
	"crypto/md5"
	"encoding/hex"
)

const secret = "sillyhumans.com"

// CheckUserByName 查询用户名是否存在
func CheckUserByName(username string) (bool, error) {
	sqlStr := `select count(user_id) from user where username=?`
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return false, err
	}
	if count != 0 {
		return false, nil
	}
	return true, nil
}

// QueryUserByName 根据用户名获取用户信息
func QueryUserByName(username string) (*models.User, error) {
	sqlStr := `select user_id, username, password from user where username=?`
	user := models.User{}
	err := db.Get(&user, sqlStr, username)
	if err != nil {
		return &user, err
	}
	return &user, nil
}

// InsertUser 向数据库中插入一条新用户
func InsertUser(user *models.User) (err error) {
	// 执行sql语句
	sqlStr := "insert into user(user_id, username, password) values(?,?,?)"
	_, err = db.Exec(sqlStr, user.UserID, user.UserName, EncryptPassword(user.Password))
	return
}

// encryptPassword 用md5对密码进行简单加密
func EncryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}
