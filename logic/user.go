package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/jwt"
	mySnowflake "bluebell/pkg/snowflake"
	"database/sql"
	"errors"
	"fmt"
)

func SignUp(p *models.ParamSignUp) (err error) {
	// 判断用户是否存在
	var exist bool
	exist, err = mysql.CheckUserByName(p.Username)
	// 数据库查找错误
	if err != nil {
		return err
	}
	// 用户已经存在
	if !exist {
		return errors.New("用户名已存在")
	}
	// 生成uid
	uid := mySnowflake.GetID()
	// 存入数据库
	user := models.User{
		UserID:   uid,
		UserName: p.Username,
		Password: p.Password,
	}
	err = mysql.InsertUser(&user)
	return
}

func Login(p *models.ParamLogin) (aToken, rToken string, err error) {
	// 判断用户是否存在
	var user *models.User
	user, err = mysql.QueryUserByName(p.Username)
	fmt.Println(p.Username)
	// 数据库查询错误
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", errors.New("用户不存在")
		}
		return "", "", err
	}
	password := mysql.EncryptPassword(p.Password)
	if password != user.Password {
		return "", "", errors.New("用户名或密码错误")
	}
	// 生成JWT的token
	return jwt.GenToken(int64(user.UserID), user.UserName)
}
