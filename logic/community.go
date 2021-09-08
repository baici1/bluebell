package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
)

// GetCommunityList 查数据库 查找到所有的社区并返回
func GetCommunityList() (data []*models.Community, err error) {
	return mysql.GetCommunityList()
}

// GetCommunityDetail 根据ID查询分类详情
func GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	return mysql.GetCommunityDetailByID(id)
}
