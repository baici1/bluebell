package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/snowflake"

	"go.uber.org/zap"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post) (err error) {
	//生成post id
	p.ID = snowflake.GenID()
	//保存到数据库
	//返回错误
	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}
	err = redis.CreatePost(p.ID)
	return
}

// GetPostById 根据id查询帖子数据
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
		//获取用户详情数据
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

func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	//redis查询id列表
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn(" redis.GetPostIDsInOrder get ids return 0 data")
	}
	zap.L().Debug("GetPostList2", zap.Any("ids", ids))
	//根据id取mysql查询post列表数据
	//按照给定的顺序返回
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	zap.L().Debug("mysql.GetPostListByIDs", zap.Any("posts", posts))
	//提前查好每个帖子投票的数据
	votes, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}
	//根据post获取作者和社区信息填充到帖子中
	for index, post := range posts {
		//获取用户详情数据
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
			VoteNum:         votes[index],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}

	return
}
