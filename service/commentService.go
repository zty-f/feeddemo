package service

import (
	"errors"
	"feeddemo/common"
	"feeddemo/repository"
	"fmt"
)

var commentDaoInstance = repository.NewCommentDaoInstance()

type CommentService struct {
}

// NewCommentServiceInstance 返回一个评论服务类的指针变量，可以方便调用该结构体的方法
func NewCommentServiceInstance() *CommentService {
	return &CommentService{}
}

// DoAddCommentAction 新增评论功能
func (c *CommentService) DoAddCommentAction(loginUserId, blogId int64, commentText string) (*common.CommentVo, error) {
	//新增评论
	cid, err3 := commentDaoInstance.CreateComment(loginUserId, blogId, commentText)
	if err3 != nil {
		return nil, err3
	}
	comment, err4 := commentDaoInstance.QueryCommentById(cid)
	if err4 != nil {
		return nil, err4
	}
	// 按照2006-01-02 15:04:05这个固定来格式化
	createDate := comment.CreateTime.Format("2006-01-02 15:04:05")
	user, err := userDaoInstance.QueryUserById(loginUserId)
	if err != nil {
		return nil, err
	}
	isFollow, err1 := relationDaoInstance.QueryIsFollowByUserIdAndToUserId(loginUserId, loginUserId)
	if err1 != nil {
		return nil, err1
	}
	userVo := &common.UserVo{
		Id:            user.Id,
		UserName:      user.UserName,
		Avatar:        user.Avatar,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		IsFollow:      isFollow,
	}
	commentVo := &common.CommentVo{
		Id:         comment.ID,
		User:       *userVo,
		Content:    comment.Content,
		CreateDate: createDate,
	}
	return commentVo, nil
}

// DoDelCommentAction 删除评论功能
func (c *CommentService) DoDelCommentAction(loginUserId, blogId, commentId int64) error {
	comment, err4 := commentDaoInstance.QueryCommentById(commentId)
	if err4 != nil {
		return err4
	}
	if comment.UserID != loginUserId {
		return errors.New("不能删除他人的评论，谢谢！")
	}
	//删除评论
	err := commentDaoInstance.DeleteComment(commentId, blogId)
	if err != nil {
		return err
	}
	return nil
}

// DoGetCommentList 获取评论列表
func (c *CommentService) DoGetCommentList(loginUserId, blogId int64) ([]common.CommentVo, error) {
	comments, err2 := commentDaoInstance.QueryCommentsByBlogId(blogId)
	if err2 != nil {
		return nil, err2
	}
	if len(comments) == 0 {
		return nil, errors.New("该博文还没有用户评论，欢迎评论^_^！")
	}
	commentList := make([]common.CommentVo, len(comments))
	fmt.Println("获取该博文评论列表成功！")
	for i, _ := range comments {
		// 按照2006-01-02 15:04:05这个固定来格式化
		createDate := comments[i].CreateTime.Format("2006-01-02 15:04:05")
		user, err := userDaoInstance.QueryUserById(comments[i].UserID)
		if err != nil {
			return nil, err
		}
		isFollow, err4 := relationDaoInstance.QueryIsFollowByUserIdAndToUserId(loginUserId, comments[i].UserID)
		if err4 != nil {
			return nil, err4
		}
		userVo := &common.UserVo{
			Id:            user.Id,
			UserName:      user.UserName,
			Avatar:        user.Avatar,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FollowerCount,
			IsFollow:      isFollow,
		}
		commentList[i] = common.CommentVo{
			Id:         comments[i].ID,
			User:       *userVo,
			Content:    comments[i].Content,
			CreateDate: createDate,
		}
	}
	return commentList, nil
}
