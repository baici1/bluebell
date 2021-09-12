package logic

import (
	"bluebell/dao/redis"
	"bluebell/models"
	"strconv"

	"go.uber.org/zap"
)

//分析投票功能
/*
投票的几种情况
direction=1时 有两种情况
	1.之前没有投过票，现在投赞成票 -->更新分数和投票记录
	2.之前投反对票，现在改投赞成票 -->更新分数和投票记录
direction=0时，有两种情况
	1.之前投过赞成票，现在取消 	  -->更新分数和投票记录
	2.之前投反对票，现在取消     -->更新分数和投票记录
direction=-1，有两种情况
	1.之前没有投过票，现在投反对票 -->更新分数和更新记录
	2.之前投赞成票，现在改投反对票 -->更新分数和更新记录

投票的限制：
每个帖子自发表之日起一个星期之内，允许用户投票，超过一个星期就不允许再投票了
	1.到期之后将redis中保存的赞成票数及反对票存储到mysql表中
*/

// VoteForPost 为帖子投票的函数
func VoteForPost(userID int64, p *models.ParamVoteData) error {
	zap.L().Debug("VoteForPost", zap.Int64("userID", userID), zap.Int64("postID", p.PostID), zap.Int8("direction", p.Direction))
	return redis.VoteForPost(strconv.Itoa(int(userID)), strconv.Itoa(int(p.PostID)), float64(p.Direction))
}
