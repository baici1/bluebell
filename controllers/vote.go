package controllers

import (
	"bluebell/logic"
	"bluebell/models"
	"bluebell/pkg/translate"
	"fmt"

	"go.uber.org/zap"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

// PostVoteController 投票功能
func PostVoteController(c *gin.Context) {
	//参数校验
	p := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("请求参数失败", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		errData := translate.RemoveTopStruct(errs.Translate(translate.Trans)) //翻译并去除掉错误的结构体标识
		ResponseErrorWithMsg(c, CodeInvalidParam, errData)
		return
	}
	fmt.Println(p)
	userID, err := GetCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	//业务逻辑
	if err := logic.VoteForPost(userID, p); err != nil {
		zap.L().Error("logic.VoteForPost failed", zap.Error(err))
		ResponseError(c, CodeOperation)
		return
	}
	//返回响应
	ResponseSuccess(c, nil)
}
