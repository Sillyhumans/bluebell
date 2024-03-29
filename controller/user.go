package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"fmt"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

// printErr 写入日志
func printErr(msg string, err error, c *gin.Context) {

}

// SignUpHandler 获取参数请求函数
func SignUpHandler(c *gin.Context) {
	// 获取参数和参数校验
	p := new(models.ParamSignUp)
	err := c.ShouldBindJSON(&p)
	fmt.Println(p)
	if err != nil {
		// 请求参数有误
		zap.L().Error("SignUp with invalid params", zap.Error(err))
		// 判断validate是否支持翻译该错误
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}

	// 业务处理
	err = logic.SignUp(p)
	if err != nil {
		// 数据库错误或用户已经存在
		zap.L().Error("SignUp failed", zap.Error(err))
		// 判断validate是否支持翻译该错误
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseErrorWithMsg(c, CodeUserExist, removeTopStruct(errs.Translate(trans)))
		return
	}
	// 返回响应
	ResponseSuccess(c, nil)
}

func LoginHandler(c *gin.Context) {
	// 获取参数
	var p models.ParamLogin
	err := c.ShouldBindJSON(&p)
	if err != nil {
		// 请求参数有误
		zap.L().Error("Login with invalid params", zap.Error(err))
		// 判断validate是否支持翻译该错误
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	// 查询数据
	aToken, rToken, err := logic.Login(&p)
	if err != nil {
		// 数据库错误
		zap.L().Error("Login failed", zap.Error(err))
		// 判断validate是否支持翻译该错误
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseErrorWithMsg(c, CodeUserNotExist, err.Error())
			return
		}
		ResponseErrorWithMsg(c, CodeUserNotExist, removeTopStruct(errs.Translate(trans)))
		return
	}
	// 返回响应
	data := map[string]string{"aToken": aToken, "rToken": rToken}
	ResponseSuccess(c, data)
}
