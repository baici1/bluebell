package redis

import (
	"bluebell/models"
)

// GetPostIDsInOrder 从redis查询id记录
func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	//从redis获取id
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	start := (p.Page - 1) * p.Size
	stop := start + p.Size + 1
	//按照分数或者时间查询指定数量的id记录
	return rdb.ZRevRange(key, start, stop).Result()
}
