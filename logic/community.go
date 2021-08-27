package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
)

// GetCommunityList 查数据库 查找到所有的社区并返回
func GetCommunityList() (data []*models.Community, err error) {
	return mysql.GetCommunityList()
}
