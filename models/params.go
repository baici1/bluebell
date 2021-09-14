package models

//定义请求的参数结构体

const (
	OrderTime  = "time"
	OrderScore = "score"
)

//ParamSignUp 请求注册的参数
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

// ParamLogin 请求登录的参数
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ParamVoteData 投票参数
type ParamVoteData struct {
	//UserID 从请求中获取当前用户
	PostID    int64 `json:"post_id,string" binding:"required"`        //帖子id
	Direction int8  `json:"direction,string" binding:"oneof=1 0 -1" ` //赞成票（1）反对票（-1）取消投票(0)
}

// ParamPostList 获取帖子列表参数
type ParamPostList struct {
	Page  int64  `json:"page" form:"page"`
	Size  int64  `json:"size" form:"size"`
	Order string `json:"order" form:"order"`
}

// ParamCommunityPostList 根据社区获取帖子列表
type ParamCommunityPostList struct {
	*ParamPostList
	CommunityID int64 `json:"community_id" form:"community_id"`
}
