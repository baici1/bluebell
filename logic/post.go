package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/snowflake"
)

func CreatePost(p *models.Post) (err error) {
	//生成post id
	p.ID = snowflake.GenID()
	//保存到数据库
	//返回错误
	return mysql.CreatePost(p)

}
func GetPostById(id int64) (data *models.Post, err error) {

	return mysql.GetPostById(id)

}
