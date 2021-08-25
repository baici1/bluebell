package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/snowflake"
)

//存业务逻辑的代码

func SignUp(p *models.ParamSignUp) (err error) {
	//判断用户是否存在
	err = mysql.CheckUserExist(p.Username)
	//数据库查询出错和用户已存在的错误
	if err != nil {
		return err
	}
	//生成userID
	userID := snowflake.GenID()
	//构造user实例
	u := models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	//保存进数据库
	err = mysql.InsertUser(&u)
	return
}

func Login(p *models.ParamLogin) (err error) {
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	return mysql.Login(user)
}
