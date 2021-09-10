package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/snowflake"

	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	//生成post id
	p.ID = snowflake.GenID()
	//保存到数据库
	//返回错误
	return mysql.CreatePost(p)

}
func GetPostById(id int64) (data *models.ApiPostDetail, err error) {
	//查询数据并组合我们接口想用的数据
	//获取帖子详情数据
	data = new(models.ApiPostDetail)
	post, err := mysql.GetPostById(id)
	if err != nil {
		zap.L().Error("mysql.GetPostById failed", zap.Int64("id", id), zap.Error(err))
		return
	}
	//获取用户详情数据
	user, err := mysql.GetUserById(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetuserById failed", zap.Int64("AuthorID", post.AuthorID), zap.Error(err))
		return
	}
	//获取社区详情数据
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID", zap.Int64("CommunityID", post.CommunityID), zap.Error(err))
		return
	}
	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
	}
	return

}

// GetPostList 查询帖子列表数据
func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		zap.L().Error(" mysql.GetPostList failed", zap.Error(err))
		return
	}
	for _, post := range posts {
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetuserById failed", zap.Int64("AuthorID", post.AuthorID), zap.Error(err))
			continue
		}
		//获取社区详情数据
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID", zap.Int64("CommunityID", post.CommunityID), zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}
