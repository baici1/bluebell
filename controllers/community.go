package controllers

import (
	"bluebell/logic"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//跟社区相关

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
