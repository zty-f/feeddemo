package service

import (
	"errors"
	"feeddemo/common"
	"feeddemo/repository"
	"fmt"
)

var favoriteDaoInstance = repository.NewFavoriteDaoInstance()

type FavoriteService struct {
}

// NewFavoriteServiceInstance 返回一个点赞服务类的指针变量，可以方便调用该结构体的方法
func NewFavoriteServiceInstance() *FavoriteService {
	return &FavoriteService{}
}

// DoFavoriteAction 点赞博文
func (f *FavoriteService) DoFavoriteAction(userId, blogId int64, actionType int32) error {
	fmt.Printf("点赞userId：%d==blogId：%d==actionType:%d\n", userId, blogId, actionType)
	flag, err := favoriteDaoInstance.QueryFavoriteByUserIdAndBlogId(userId, blogId)
	if err != nil {
		return err
	}
	if flag {
		err = favoriteDaoInstance.UpdateFavorite(userId, blogId, actionType)
		if err != nil {
			return err
		}
	} else {
		err = favoriteDaoInstance.CreateFavorite(userId, blogId, actionType)
		if err != nil {
			return err
		}
	}
	return nil
}

// DoFavoriteList 获取点赞博文列表
func (f *FavoriteService) DoFavoriteList(userId, loginUserId int64) ([]common.BlogVo, error) {
	fmt.Printf("获取点赞博文列表userId：%d\n", userId)
	ids, err := favoriteDaoInstance.QueryBlogsIdByUserId(userId)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, errors.New("该用户未点赞任何博文！")
	}
	blogList, err1 := blogDaoInstance.QueryByIds(ids)
	if err1 != nil {
		return nil, err1
	}
	blogListResp := make([]common.BlogVo, len(blogList))
	fmt.Println("获取点赞博文列表成功！")
	for i, _ := range blogList {
		var isFavorite bool
		user, err2 := userDaoInstance.QueryUserById(blogList[i].UserId)
		if err2 != nil {
			return nil, err2
		}
		actionType, err3 := favoriteDaoInstance.QueryActionTypeByUserIdAndBlogId(loginUserId, blogList[i].Id)
		if err3 != nil {
			return nil, err3
		}
		isFollow, err4 := relationDaoInstance.QueryIsFollowByUserIdAndToUserId(loginUserId, userId)
		if err4 != nil {
			return nil, err4
		}
		if actionType == 1 {
			isFavorite = true
		} else {
			isFavorite = false
		}
		tmpUser := &common.UserVo{
			Id:            user.Id,
			UserName:      user.UserName,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FollowerCount,
			IsFollow:      isFollow,
		}
		blogListResp[i] = common.BlogVo{
			Id:            blogList[i].Id,
			Author:        *tmpUser,
			Title:         blogList[i].Title,
			Content:       blogList[i].Content,
			FavoriteCount: blogList[i].FavoriteCount,
			CommentCount:  blogList[i].CommentCount,
			IsFavorite:    isFavorite,
			CreateTime:    blogList[i].CreateTime,
			UpdateTime:    blogList[i].UpdateTime,
		}
	}
	return blogListResp, nil
}
