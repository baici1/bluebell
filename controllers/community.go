package controllers

import (
	"bluebell/logic"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//跟社区相关

// CommunityHandler 社区分类信息
func CommunityHandler(c *gin.Context) {
	//查询到所有社区 （community_id,community_name）以列表的形式
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) //不轻易把服务端错误返回到前端
		return
	}
	ResponseSuccess(c, data)
}

// CommunityDetailHandler 社区分类详情
func CommunityDetailHandler(c *gin.Context) {
	//1.获取参数ID
	communityID := c.Param("id")
	//处理参数
	id, err := strconv.ParseInt(communityID, 10, 64)
	if err != nil {
		zap.L().Error("strconv.ParseInt failed", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//进去业务处理
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityDetail failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//返回正确信息
	ResponseSuccess(c, data)
}
