package repository

import (
	"fmt"
	"time"
)

const MaxListLength = 5

type Blog struct {
	Id            int64     `gorm:"primaryKey"`
	Title         string    `gorm:"column:title;size:128"`
	Content       string    `gorm:"column:content;size:500"`
	FavoriteCount int64     `gorm:"column:favorite_count"`
	CommentCount  int64     `gorm:"column:comment_count"`
	UserId        int64     `gorm:"column:user_id"`
	Top           string    `gorm:"column:top;size:24"`
	CreateTime    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdateTime    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

type BlogDao struct {
}

// NewBlogDaoInstance 返回一个博文实体类的指针变量，可以方便调用该结构体的方法
func NewBlogDaoInstance() *BlogDao {
	return &BlogDao{}
}

// QueryByOwner 通过用户id查询该用户发布的所有博文
func (d *BlogDao) QueryByOwner(ownerId int64) ([]Blog, error) {
	var blogs []Blog
	if err := db.Order("top desc,id desc").Where("user_id=?", ownerId).Find(&blogs).Error; err != nil {
		return nil, err
	}
	return blogs, nil
}

// QueryPublishCountByUserId 通过用户id查询该用户发布的博文数量
func (d *BlogDao) QueryPublishCountByUserId(userId int64) (int64, error) {
	var count int64
	if err := db.Table("blogs").Where("user_id=?", userId).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// QueryTotalFavoriteCountByUserId 通过用户id查询该用户发布的所有博文总获赞数量
func (d *BlogDao) QueryTotalFavoriteCountByUserId(userId int64) (int64, error) {
	var totalFavoriteCount int64
	fmt.Println("通过userId查询所有已发布博文的总获赞数")
	if err := db.Table("blogs").Select("sum(favorite_count) as total").Where("user_id = ?", userId).Take(&totalFavoriteCount).Error; err != nil {
		return 0, err
	}
	return totalFavoriteCount, nil
}

// QueryByIds 通过一组博文id获取对应的博文列表
func (d *BlogDao) QueryByIds(ids []int64) ([]Blog, error) {
	var blogs []Blog
	if err := db.Order("id desc").Find(&blogs, ids).Error; err != nil {
		return nil, err
	}
	return blogs, nil
}

// QueryBlogById 通过一组博文id获取对应的博文列表
func (d *BlogDao) QueryBlogById(id int64) (*Blog, error) {
	var blog = &Blog{}
	if err := db.First(blog, id).Error; err != nil {
		return nil, err
	}
	return blog, nil
}

// QueryByIdsWithTop 通过一组博文id获取对应的博文列表 含置顶
func (d *BlogDao) QueryByIdsWithTop(ids []int64) ([]Blog, error) {
	var blogs []Blog
	if err := db.Order("top desc,id desc").Find(&blogs, ids).Error; err != nil {
		return nil, err
	}
	return blogs, nil
}

// CreateBlogRecord 通过传入参数创建新的博文记录
func (d *BlogDao) CreateBlogRecord(userId int64, title string, content string) (int64, error) {
	var blog = &Blog{
		Title:         title,
		Content:       content,
		FavoriteCount: 0,
		CommentCount:  0,
		UserId:        userId,
	}
	if err := db.Create(blog).Error; err != nil {
		return 0, err
	}
	return blog.Id, nil
}

// DelBlogRecord 根据id删除博文
func (d *BlogDao) DelBlogRecord(blogID int64) error {
	blog := &Blog{}
	if err := db.Where("id = ?", blogID).Delete(blog).Error; err != nil {
		return err
	}
	return nil
}

// QueryNewFeedFlow 获取平台最新30条博客流
func (d *BlogDao) QueryNewFeedFlow(time string) ([]Blog, error) {
	var blogs []Blog
	if err := db.Where("create_time>?", time).Order("id desc").Limit(MaxListLength).Find(&blogs).Error; err != nil {
		return nil, err
	}
	return blogs, nil
}

// QueryHotFeedFlow 获取平台最火热30条博客流  根据点赞和评论数量总和排序
func (d *BlogDao) QueryHotFeedFlow() ([]Blog, error) {
	var blogs []Blog
	if err := db.Select("*,(favorite_count+comment_count) as sum").Order("sum desc,id desc").Limit(10).Find(&blogs).Error; err != nil {
		return nil, err
	}
	return blogs, nil
}

// TopBlog 置顶个人的一篇博客
func (d *BlogDao) TopBlog(blogId, actionType int64) error {
	blog := &Blog{}
	if actionType == 1 {
		now := time.Now().Format("2006-01-02 15:04:05")
		if err := db.Model(blog).Where("id = ? ", blogId).Update("top", now).Error; err != nil {
			return err
		}
		return nil
	}
	if err := db.Model(blog).Where("id = ? ", blogId).Update("top", 0).Error; err != nil {
		return err
	}
	return nil
}
