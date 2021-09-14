package redis

import (
	"errors"
	"math"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/go-redis/redis"
)

const (
	oneWeekInSeconds = 7 * 24 * 3600
	//每一票占的分数
	scorePerVote = 432
)

//用户投一票加432分，  86400/200 ->200张赞成票可以给帖子续一天 -> 《redis实战》

// ErrorVoteTimeExpire 错误码
var (
	ErrorVoteTimeExpire = errors.New("投票时间已过")
	ErrorVoteRepeated   = errors.New("重复投票")
)

//分析投票功能
/*
投票的几种情况
direction=1时 有两种情况
	1.之前没有投过票，现在投赞成票 -->更新分数和投票记录 差值绝对值：1 +432
	2.之前投反对票，现在改投赞成票 -->更新分数和投票记录 差值绝对值：2 +432*2
direction=0时，有两种情况
	2.之前投反对票，现在取消     -->更新分数和投票记录 差值绝对值：1 +432
	1.之前投过赞成票，现在取消 	  -->更新分数和投票记录 差值绝对值：1  -432

direction=-1，有两种情况
	1.之前没有投过票，现在投反对票 -->更新分数和更新记录 差值绝对值：1 -432
	2.之前投赞成票，现在改投反对票 -->更新分数和更新记录 差值绝对值：2 -432*2

当当前的投票值比之前的大，那么就时加，反之时减
投票的限制：
每个帖子自发表之日起一个星期之内，允许用户投票，超过一个星期就不允许再投票了
	1.到期之后将redis中保存的赞成票数及反对票存储到mysql表中
*/

func VoteForPost(userID, postID string, direction float64) error {
	//1.判断投票限制
	//在redis取帖子发布时间
	postTime := rdb.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	//如果帖子发布时间超过7天则无法点赞
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrorVoteTimeExpire
	}
	//2和3需要放到一个事务执行

	//2.更新帖子分数
	//先查询先前投票记录
	odir := rdb.ZScore(getRedisKey(KeyPostVotedZSetPrefix+postID), userID).Val()
	//如果这一次投票与之前保存的一致，就报错（不允许重复投票）
	if direction == odir {
		return ErrorVoteRepeated
	}
	var pn float64
	//根据新投票记录和旧地投票记录相比 高的就是需要加分 低地需要减分

	if direction > odir {
		pn = 1
	} else {
		pn = -1
	}
	//计算新旧数据的差值
	diff := math.Abs(odir - direction)
	pipeline := rdb.TxPipeline()
	//更新分数
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), diff*pn*scorePerVote, postID)
	//3.记录用户为该帖子投过票
	if direction == 0 {
		pipeline.ZRem(getRedisKey(KeyPostVotedZSetPrefix+postID), userID)
	} else {
		pipeline.ZAdd(getRedisKey(KeyPostVotedZSetPrefix+postID), redis.Z{
			Score:  direction, //赞成还是反对票
			Member: userID,    //用户ID
		})
	}
	_, err := pipeline.Exec()
	if err != nil {
		zap.L().Error("更新分数或者记录投票 失败", zap.Error(err))
		return err
	}
	return err
}

// CreatePost 实现redis 添加文章的时间和分数
func CreatePost(postID, communityID int64) error {
	//事务执行 为文章添加时间，为文章添加分数
	pipeline := rdb.TxPipeline()
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  0,
		Member: postID,
	})
	//补充：把帖子id加到社区set
	ckey := getRedisKey(KeyCommunitySetPrefix + strconv.Itoa(int(communityID)))
	pipeline.SAdd(ckey, postID)
	_, err := pipeline.Exec()
	return err
}
