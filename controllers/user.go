package controllers

import (
	"bluebell/logic"
	"bluebell/models"
	"bluebell/pkg/translate"
	"fmt"
	"net/http"

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
		zap.L().Error("[ERRER]SignUp with invalid param", zap.Error(err))
		// 获取validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"msg": translate.RemoveTopStruct(errs.Translate(translate.Trans)),
		})
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
	fmt.Println(p)
	//2.业务处理
	if err := logic.SignUp(&p); err != nil {
		zap.L().Error("login.SignUp failed", zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"msg": "注册失败",
		})
		return
	}
	//3.返回响应
	c.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

func LoginHandler(c *gin.Context) {
	//获取参数校验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误，直接返回响应
		zap.L().Error("[ERRER]Login with invalid param", zap.Error(err))
		// 获取validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
			return
		}
		//翻译错误
		c.JSON(http.StatusOK, gin.H{
			"msg": translate.RemoveTopStruct(errs.Translate(translate.Trans)),
		})
		return
	}
	//业务逻辑处理
	if err := logic.Login(p); err != nil {
		zap.L().Error("[ERRER]logic.Login failed", zap.String("username", p.Username), zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"msg": "用户名或者密码错误",
		})
		return
	}
	//返回响应
	c.JSON(http.StatusOK, gin.H{
		"msg": "登录成功",
	})
}
