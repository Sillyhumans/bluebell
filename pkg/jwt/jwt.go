package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const TokenExpireDuration = time.Hour * 2

var mySecret = []byte("sillyhumans")

type MyClaims struct {
	UserID   int64  `json:"user_id"`
	UserName string `json:"username"`
	jwt.StandardClaims
}

// GenToken 获取token
func GenToken(userID int64, username string) (aToken, rToken string, err error) {
	// 创建声明 access token
	c := MyClaims{
		userID,
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			Issuer:    "bubble",
		},
	}
	// 使用指定的签名方法创建签名对象
	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(mySecret)

	// 创建声明 refresh token 不需要存任何自定义数据
	rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
		Issuer:    "bubble",
	}).SignedString(mySecret)
	// 使用指定secret签名并返回
	return
}

// ParseToken 解析token
func ParseToken(tokenString string) (claims *MyClaims, err error) {
	var token *jwt.Token
	claims = new(MyClaims)
	token, err = jwt.ParseWithClaims(tokenString, claims, myKey)
	if err != nil {
		return
	}
	if !token.Valid {
		err = errors.New("invalid token")
	}
	return
}

func RefreshToken(aToken, rToken string) (newAToken, newRToken string, err error) {
	// refresh Token无效直接返回
	if _, err = jwt.Parse(rToken, myKey); err != nil {
		return
	}
	// 从旧access token中解析出claim数据
	var claims MyClaims
	_, err = jwt.ParseWithClaims(aToken, &claims, myKey)
	v, _ := err.(*jwt.ValidationError)

	// 当access token是过期错误并且 refresh token没有过期时就创建一个新的access token
	if v.Errors == jwt.ValidationErrorExpired {
		return GenToken(claims.UserID, claims.UserName)
	}
	return
}

func myKey(*jwt.Token) (i interface{}, err error) {
	return mySecret, nil
}
