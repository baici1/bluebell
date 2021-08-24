package logic

import (
	"bluebell/dao/mysql"
	"bluebell/pkg/snowflake"
)

//存业务逻辑的代码

func SignUp() {
	//判断用户是否存在
	mysql.QueryUserByUsername()
	//生成userID
	snowflake.GenID()
	//密码加密
	//保存进数据库
	mysql.InsertUser()
}
