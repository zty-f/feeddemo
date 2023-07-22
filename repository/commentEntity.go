package repository

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Comment struct {
	ID         int64
	UserID     int64
	BlogID     int64
	Content    string
	CreateTime time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

type CommentDao struct {
}

// NewCommentDaoInstance 返回一个评论表实体类的指针变量，可以方便调用该结构体的方法
func NewCommentDaoInstance() *CommentDao {
	return &CommentDao{}
}

// CreateComment 新增评论返回评论id
func (c *CommentDao) CreateComment(userId, blogId int64, commentText string) (int64, error) {
	blog := &Blog{}
	comment := &Comment{
		UserID:  userId,
		BlogID:  blogId,
		Content: commentText,
	}
	tx := db.Begin()
	if err := tx.Select("user_id", "blog_id", "content").Create(&comment).Error; err != nil {
		tx.Rollback()
		return comment.ID, err
	}
	// 对应博文评论数+1
	if err := tx.Model(blog).Where("id = ? ", blogId).Update("comment_count", gorm.Expr("comment_count+ ?", 1)).Error; err != nil {
		tx.Rollback()
		return comment.ID, err
	}
	tx.Commit()
	// 返回comment对象只包含传入字段以及主键id
	return comment.ID, nil
}

// DeleteComment 删除评论
func (c *CommentDao) DeleteComment(commentId, blogId int64) error {
	Blog := &Blog{}
	comment := &Comment{}
	tx := db.Begin()
	err := tx.Delete(comment, commentId).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	// 对应博文评论数-1
	if err1 := db.Model(Blog).Where("id = ? ", blogId).Update("comment_count", gorm.Expr("comment_count- ?", 1)).Error; err1 != nil {
		tx.Rollback()
		return err1
	}
	tx.Commit()
	return nil
}

// DeleteBlogComment 根据ID删除博文全部评论
func (c *CommentDao) DeleteBlogComment(blogId int64) error {
	comment := &Comment{}
	if err := db.Where("blog_id = ?", blogId).Delete(comment).Error; err != nil {
		return err
	}
	return nil
}

// QueryCommentById 通过id查询评论
func (c *CommentDao) QueryCommentById(commentId int64) (*Comment, error) {
	comment := &Comment{}
	err := db.First(comment, commentId).Error
	if err != nil {
		return comment, err
	}
	return comment, nil
}

// QueryCommentsByBlogId 通过博文id查询该博文所有评论
func (c *CommentDao) QueryCommentsByBlogId(blogId int64) ([]Comment, error) {
	var comments []Comment
	fmt.Println("通过blogId查询所有评论")
	err := db.Order("create_time desc").Where("blog_id = ?", blogId).Find(&comments).Error
	if err != nil {
		return comments, err
	}
	return comments, nil
}
