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

// GetPostListHandler2 根据前端传来的参数动态获取帖子列表接口
//按创建时间或者点赞分数
func GetPostListHandler2(c *gin.Context) {
	//获取flag（获取时间排序的帖子还是点赞分数）
	//初始化结构体指定初始参数
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler2 with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//获取帖子的数据
	data, err := logic.GetPostListNew(p) //更新：两种（查全部或者按社区）查询帖子列表数据合二为一
	if err != nil {
		zap.L().Error("logic.GetPostListNew", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//返回信息
	ResponseSuccess(c, data)

}

// GetCommunityPostListHandler 根据社区查询帖子列表
//func GetCommunityPostListHandler(c *gin.Context) {
//	//获取flag（获取时间排序的帖子还是点赞分数）
//	//初始化结构体指定初始参数
//	p := &models.ParamCommunityPostList{
//		ParamPostList: &models.ParamPostList{
//			Page:  1,
//			Size:  10,
//			Order: models.OrderTime,
//		},
//		CommunityID: 1,
//	}
//	//获取参数
//	if err := c.ShouldBindQuery(p); err != nil {
//		zap.L().Error("ParamCommunityPostList with invalid params", zap.Error(err))
//		ResponseError(c, CodeInvalidParam)
//		return
//	}
//	//获取帖子的数据
//	data, err := logic.GetCommunityPostList2(p)
//	if err != nil {
//		zap.L().Error("logic.GetCommunityPostList2", zap.Error(err))
//		ResponseError(c, CodeServerBusy)
//		return
//	}
//	//返回信息
//	ResponseSuccess(c, data)
//}
