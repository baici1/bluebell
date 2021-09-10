package controllers

import (
	"bluebell/logic"
	"bluebell/models"
	"bluebell/pkg/translate"
	"strconv"

	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// CreatePostHandler 创建帖子的函数
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

// GetPostDetailHandler 获取帖子详情的函数
func GetPostDetailHandler(c *gin.Context) {
	//获取帖子的id
	pidStr := c.Param("id")
	id, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//根据id取出数据
	data, err := logic.GetPostById(id)
	if err != nil {
		zap.L().Error("logic.GetPostById", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler 获取帖子列表函数
func GetPostListHandler(c *gin.Context) {
	//获取分页参数
	//pageStr := c.Query("offset")
	//sizeStr := c.Query("limit")
	//var (
	//	page int64
	//	size int64
	//	err  error
	//)
	//page, err = strconv.ParseInt(pageStr, 10, 64)
	//if err != nil {
	//	page = 0
	//}
	//size, err = strconv.ParseInt(sizeStr, 10, 64)
	//if err != nil {
	//	size = 10
	//}
	//整合到一个函数里面
	page, size := getPageInfo(c)

	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}
