package controllers

import (
	"bluebell/logic"
	"bluebell/models"
	"bluebell/pkg/translate"

	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func CreatePostHandler(c *gin.Context) {
	//获取参数及参数校验
	var p models.Post
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("CreatePostHandler with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			//返回非validator错误
			ResponseError(c, CodeInvalidParam)
			return
		}
		//返回validator错误
		ResponseErrorWithMsg(c, CodeInvalidParam, translate.RemoveTopStruct(errs.Translate(translate.Trans)))
		return
	}
	//取到当前发起请求的用户ID
	userID, err := GetCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
	}
	p.AuthorID = userID
	//创建帖子
	if err := logic.CreatePost(&p); err != nil {
		zap.L().Error("logic.CreatePost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//返回响应
	ResponseSuccess(c, nil)
}
