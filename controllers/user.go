package controllers

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
	"bluebell/pkg/translate"
	"errors"

	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// SignUpHandler 用户注册功能
func SignUpHandler(c *gin.Context) {
	//1.获取参数和参数校验
	var p models.ParamSignUp
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误，直接返回响应
		zap.L().Error("[ERROR]SignUp with invalid param", zap.Error(err))
		// 获取validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, translate.RemoveTopStruct(errs.Translate(translate.Trans)))
		return
	}
	//手动对请求参数进行详细业务规则校验
	//if len(p.Username) == 0 || len(p.Password) == 0 || len(p.RePassword) == 0 || p.RePassword != p.Password {
	//	zap.L().Error("[ERRER]SignUp with invalid param")
	//	c.JSON(http.StatusOK, gin.H{
	//		"msg": "请求参数有误",
	//	})
	//	return
	//}
	//自动 利用第三方库进行字段校验validator库
	//fmt.Println(p)
	//2.业务处理
	if err := logic.SignUp(&p); err != nil {
		zap.L().Error("login.SignUp failed", zap.Error(err))
		//c.JSON(http.StatusOK, gin.H{
		//	"msg": "注册失败",
		//})
		//判断错误是否是用户存在
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		//系统错误，没有数据保存进去
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	//c.JSON(http.StatusOK, gin.H{
	//	"msg": "ok",
	//})
	ResponseSuccess(c, nil)
}

func LoginHandler(c *gin.Context) {
	//获取参数校验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误，直接返回响应
		zap.L().Error("[ERROR]Login with invalid param", zap.Error(err))
		// 获取validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			//c.JSON(http.StatusOK, gin.H{
			//	"msg": err.Error(),
			//})
			ResponseError(c, CodeInvalidParam)
			return
		}
		//翻译错误
		//c.JSON(http.StatusOK, gin.H{
		//	"msg": translate.RemoveTopStruct(errs.Translate(translate.Trans)),
		//})
		ResponseErrorWithMsg(c, CodeInvalidParam, translate.RemoveTopStruct(errs.Translate(translate.Trans)))
		return
	}
	//业务逻辑处理
	token, err := logic.Login(p)
	if err != nil {
		zap.L().Error("[ERROR]logic.Login failed", zap.String("username", p.Username), zap.Error(err))
		//c.JSON(http.StatusOK, gin.H{
		//	"msg": "用户名或者密码错误",
		//})
		if errors.Is(err, mysql.ErrorInvalidPassword) {
			ResponseError(c, CodeInvalidPassword)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	//返回响应
	//c.JSON(http.StatusOK, gin.H{
	//	"msg": "登录成功",
	//})
	ResponseSuccess(c, token)
}
