package redis

// redis key尽量使用命名空间方式，方便查询和拆分

const (
	KeyPrefix              = "bluebell:"
	KeyPostTimeZSet        = "post:time"   //zset;帖子及发帖时间
	KeyPostScoreZSet       = "post:score"  //zset帖子及投票分数
	KeyPostVotedZSetPrefix = "post:voted:" //zset;记录用户及投票类型 参数是帖子--post_id
	KeyCommunitySetPrefix  = "community:"  //set:保存每个分区帖子的id
)

//为key加上前缀
func getRedisKey(key string) string {
	return KeyPrefix + key
}
