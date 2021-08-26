package mysql

import (
	"bluebell/models"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

//把每一步数据库封装成函数
//给到logic层使用

//业务中不要随便出现字符串，尽量都用变量去代替
//区分错误
var (
	ErrorUserExist       = errors.New("用户已存在！")
	ErrorUserNotExist    = errors.New("用户不存在")
	ErrorInvalidPassword = errors.New("密码错误")
)

// CheckUserExist 检查用户是否在数据库中
func CheckUserExist(username string) error {
	sqlStr := "select count(user_id) from user where username=?"
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
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
		return oPassword
	}
	return string(hashPassword)
}

// CheckPassword 验证加密后的密码是否一致
func CheckPassword(hashPassword, oPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(oPassword))
	if err != nil {
		return false
	}
	return true
}

// Login 用户登录功能
func Login(user *models.User) error {
	oPassword := user.Password
	//判断用户是否存在，获取密码
	sqlStr := `select user_id,username,password from user where username=?`
	err := db.Get(user, sqlStr, user.Username)
	//判断用户是否存在
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}
	if err != nil {
		//查询数据库出错
		return err
	}
	//判断密码是否正确
	if !CheckPassword(user.Password, oPassword) {
		return ErrorInvalidPassword
	}
	return nil
}
