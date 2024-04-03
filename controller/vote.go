package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"fmt"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

func PostVoteController(c *gin.Context) {
	// 参数校验
	p := new(models.VoteData)
	err := c.ShouldBindJSON(p)
	if err != nil {
		// 请求参数有误
		zap.L().Error("Post with invalid params", zap.Error(err))
		// 判断validate是否支持翻译该错误
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	// 获取当前用户id
	userID, err := GetCurrentUser(c)
	fmt.Println(userID)
	p.UserID = userID
	if err != nil {
		// 请求参数有误
		zap.L().Error("GetCurrentUser() failed", zap.Error(err))
		// 判断validate是否支持翻译该错误
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeServerBusy)
			return
		}
		ResponseErrorWithMsg(c, CodeServerBusy, removeTopStruct(errs.Translate(trans)))
		return
	}

	// 点赞信息存入
	err = logic.PostVote(p)
	if err != nil {
		// 请求参数有误
		zap.L().Error("PostVote() failed", zap.Error(err))
		// 判断validate是否支持翻译该错误
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeServerBusy)
			return
		}
		ResponseErrorWithMsg(c, CodeServerBusy, removeTopStruct(errs.Translate(trans)))
		return
	}
	// 返回
	ResponseSuccess(c, nil)
}
