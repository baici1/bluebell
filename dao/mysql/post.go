package mysql

import (
	"bluebell/models"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post) error {
	sqlStr := `insert into post(post_id,title,content,author_id,community_id) values(?,?,?,?,?)`
	_, err := db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return err
}

// GetPostById 根据id查询帖子详情信息
func GetPostById(id int64) (data *models.Post, err error) {
	data = new(models.Post)
	sqlStr := `select post_id,title,content,author_id,community_id,create_time from post where post_id=?`
	err = db.Get(data, sqlStr, id)
	return data, err
}

// GetPostList 获取帖子列表数据
func GetPostList(page, size int64) (data []*models.Post, err error) {
	sqlStr := `select post_id,title,content,author_id,community_id,create_time from post ORDER BY create_time DESC limit ?,?`
	err = db.Select(&data, sqlStr, (page-1)*size, size)
	return
}

// GetPostListByIDs 根据给定的id列表查询帖子数据
func GetPostListByIDs(ids []string) (postlist []*models.Post, err error) {
	sqlStr := `select post_id,title,content,author_id,community_id,create_time 
	from post 
	where post_id in (?) 
	order by FIND_IN_SET(post_id,?)`
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}
	fmt.Println(query)
	fmt.Println(args)
	query = db.Rebind(query)
	err = db.Select(&postlist, query, args...)
	return
}
