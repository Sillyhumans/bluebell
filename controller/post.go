package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

var (
	orderTime = "time"
)

func CreatePostHandler(c *gin.Context) {
	// 1.获取参数及参数校验
	p := new(models.Post)
	err := c.ShouldBindJSON(p)
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
	// 从c中获取当前用户id
	userID, err := GetCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID
	// 2.创建post
	err = logic.CreatePost(p)
	if err != nil {
		// 请求参数有误
		zap.L().Error("CreatePost failed", zap.Error(err))
		// 判断validate是否支持翻译该错误
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeServerBusy)
			return
		}
		ResponseErrorWithMsg(c, CodeServerBusy, removeTopStruct(errs.Translate(trans)))
		return
	}
	//3. 返回响应
	ResponseSuccess(c, nil)
}

func GetPostDetailHandler(c *gin.Context) {
	// 1.获取参数
	id := c.Param("id")

	pid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 2.根据id取数据
	data, err := logic.GetPostByID(pid)
	if err != nil {
		zap.L().Error("logic.GetPostByID(pid) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3.返回
	ResponseSuccess(c, data)
}

func GetPostListHandler(c *gin.Context) {
	// 获取参数
	p := models.ParaPostList{
		Page:  1,
		Size:  10,
		Order: orderTime,
	}

	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 获取数据
	data, err := logic.GetPostList(&p)
	if err != nil {
		// 请求参数有误
		zap.L().Error("GetPostList failed", zap.Error(err))
		// 判断validate是否支持翻译该错误
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeServerBusy)
			return
		}
		ResponseErrorWithMsg(c, CodeServerBusy, removeTopStruct(errs.Translate(trans)))
		return
	}
	ResponseSuccess(c, data)
}
