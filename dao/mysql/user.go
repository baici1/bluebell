package mysql

import (
	"bluebell/models"
	"errors"

	"go.uber.org/zap"

	"golang.org/x/crypto/bcrypt"
)

//把每一步数据库封装成函数
//给到logic层使用

// CheckUserExist 检查用户是否在数据库中
func CheckUserExist(username string) error {
	sqlStr := "select count(user_id) from user where username=?"
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return errors.New("用户已存在！")
	}
	return nil
}

// InsertUser 向数据库中插入用户信息
func InsertUser(user *models.User) (err error) {
	//给密码加密
	user.Password = encryptPassword(user.Password)
	//执行语句将user信息入库
	sqlStr := "insert into user(user_id,username,password) values(?,?,?);"
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	return
}

// encryptPassword Hash加密密码
func encryptPassword(oPassword string) string {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(oPassword), bcrypt.MinCost)
	if err != nil {
		zap.L().Error("加密密码失败", zap.Any("密码", hashPassword))
		return oPassword
	}
	return string(hashPassword)
}
