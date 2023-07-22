package service

import (
	"errors"
	"feeddemo/common"
	"feeddemo/repository"
	"feeddemo/utils"
	"fmt"
	"strconv"
	"time"
)

var blogDaoInstance = repository.NewBlogDaoInstance()

type BlogService struct {
}

// NewBlogServiceInstance 返回一个博文流服务类的指针变量，可以方便调用该结构体的方法
func NewBlogServiceInstance() *BlogService {
	return &BlogService{}
}

// DoPublishBlog 发布博文
func (b *BlogService) DoPublishBlog(title string, content string, userId int64) error {
	if len(title) > 100 || len(content) > 500 {
		return errors.New("内容不能超过指定长度，请重新设计~")
	}
	blogId, err := blogDaoInstance.CreateBlogRecord(userId, title, content)
	if err != nil {
		return err
	}
	now := time.Now().UnixNano()
	// 推送到个人的缓存博客列表
	key := "personal_feeds:" + strconv.FormatInt(userId, 10)
	utils.ZAdd(key, now, blogId)
	// 推送到粉丝的列表中 -- 后续这里应该采用异步消息队列解决性能问题
	// 先获取粉丝 id 集合
	followers, err2 := relationDaoInstance.QueryFollowerIdsByUserId(userId)
	if err2 != nil {
		return err2
	}
	// 推送 Feed
	for _, follower := range followers {
		key1 := "following_feeds:" + strconv.FormatInt(follower, 10)
		utils.ZAdd(key1, now, blogId)
	}
	return nil
}

// DoDelBlog 删除博文
func (b *BlogService) DoDelBlog(userId, blogId int64) error {
	err := blogDaoInstance.DelBlogRecord(blogId)
	if err != nil {
		return err
	}
	err1 := commentDaoInstance.DeleteBlogComment(blogId)
	if err1 != nil {
		return err1
	}
	err2 := favoriteDaoInstance.DelBlogFavorite(blogId)
	if err2 != nil {
		return err2
	}
	// 将内容从个人的集合中删除
	key := "personal_feeds:" + strconv.FormatInt(userId, 10)
	utils.ZRem(key, blogId)
	// 将内容从粉丝的集合中删除 -- 异步消息队列优化
	// 先获取我的粉丝
	followers, err2 := relationDaoInstance.QueryFollowerIdsByUserId(userId)
	if err2 != nil {
		return err2
	}
	// 移除 Feed
	for _, follower := range followers {
		key1 := "following_feeds:" + strconv.FormatInt(follower, 10)
		utils.ZRem(key1, blogId)
	}
	return nil
}

// GetBlogListByUserID 获取指定用户ID的博文列表
func (b *BlogService) GetBlogListByUserID(lastMinTime float64, offset, userId, loginUserId int64) ([]common.BlogVo, float64, int64, error) {
	// 走缓存 每次五条 滚动分页
	key := "personal_feeds:" + strconv.FormatInt(loginUserId, 10)
	result := utils.ZRevByScoreWithScores(key, 0, lastMinTime, offset, 5)
	if result == nil || len(result) <= 0 {
		return nil, 0, 0, nil
	}
	feedIds := make([]int64, len(result))
	//转int64
	var minTime float64 = 0
	var os int64 = 1
	for i, v := range result {
		id, err := strconv.ParseInt(v.Member.(string), 10, 64)
		if err != nil {
			return nil, 0, 0, err
		}
		if v.Score == minTime {
			os++
		} else {
			minTime = v.Score
			os = 1
		}
		feedIds[i] = id
	}
	fmt.Println(feedIds)
	blogs, err := blogDaoInstance.QueryByIdsWithTop(feedIds)
	if err != nil {
		return nil, 0, 0, nil
	}
	user, err1 := userDaoInstance.QueryUserById(userId)
	if err1 != nil {
		return nil, 0, 0, err1
	}
	isFollow, err3 := relationDaoInstance.QueryIsFollowByUserIdAndToUserId(loginUserId, userId)
	if err3 != nil {
		return nil, 0, 0, err3
	}
	curUser := &common.UserVo{
		Id:            user.Id,
		UserName:      user.UserName,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		IsFollow:      isFollow,
	}
	var BlogList = make([]common.BlogVo, len(blogs))
	for i, v := range BlogList {
		var isFavorite bool
		actionType, err2 := favoriteDaoInstance.QueryActionTypeByUserIdAndBlogId(loginUserId, v.Id)
		if err2 != nil {
			return nil, 0, 0, err2
		}
		if actionType == 1 {
			isFavorite = true
		} else {
			isFavorite = false
		}
		BlogList[i] = common.BlogVo{
			Id:            v.Id,
			Author:        *curUser,
			Title:         v.Title,
			Content:       v.Content,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			Top:           v.Top,
			IsFavorite:    isFavorite,
			CreateTime:    v.CreateTime,
			UpdateTime:    v.UpdateTime,
		}
	}
	return BlogList, minTime, os, nil
}

// GetNewFeed 获取平台最新5条博客流
func (b *BlogService) GetNewFeed(loginUserId, lastMaxTime int64) ([]common.BlogVo, int64, error) {
	timeStr := time.Unix(0, lastMaxTime).Format("2006-01-02 15:04:05")
	fmt.Println(timeStr)
	blogList, err := blogDaoInstance.QueryNewFeedFlow(timeStr)
	if err != nil || len(blogList) <= 0 {
		return nil, 0, err
	}
	blogListResp := make([]common.BlogVo, len(blogList))
	fmt.Println("获取博文流成功！")
	for i, v := range blogList {
		var isFavorite bool
		user, err1 := userDaoInstance.QueryUserById(v.UserId)
		if err1 != nil {
			return nil, 0, err1
		}
		actionType, err2 := favoriteDaoInstance.QueryActionTypeByUserIdAndBlogId(loginUserId, v.Id)
		if err2 != nil {
			return nil, 0, err2
		}
		if actionType == 1 {
			isFavorite = true
		} else {
			isFavorite = false
		}
		favoriteCount, err3 := favoriteDaoInstance.QueryFavoriteCountByUserId(v.UserId)
		if err3 != nil {
			return nil, 0, err3
		}
		var totalFavorited = int64(0)
		if count, err4 := blogDaoInstance.QueryPublishCountByUserId(v.UserId); err4 != nil {
			return nil, 0, err4
		} else if count > 0 {
			totalFavorited, err = blogDaoInstance.QueryTotalFavoriteCountByUserId(v.UserId)
		}
		if err != nil {
			return nil, 0, err
		}
		isFollow, err5 := relationDaoInstance.QueryIsFollowByUserIdAndToUserId(loginUserId, v.UserId)
		if err5 != nil {
			return nil, 0, err5
		}
		tmpUser := &common.UserVo{
			Id:             user.Id,
			UserName:       user.UserName,
			FollowCount:    user.FollowCount,
			FollowerCount:  user.FollowerCount,
			IsFollow:       isFollow,
			Avatar:         user.Avatar,
			TotalFavorited: totalFavorited,
			FavoriteCount:  favoriteCount,
		}
		blogListResp[i] = common.BlogVo{
			Id:            v.Id,
			Author:        *tmpUser,
			Title:         v.Title,
			Content:       v.Content,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			Top:           v.Top,
			IsFavorite:    isFavorite,
			CreateTime:    v.CreateTime,
			UpdateTime:    v.UpdateTime,
		}
	}
	return blogListResp, blogListResp[0].CreateTime.UnixNano(), err
}

// GetHotFeed 获取平台最热10条博客流
func (b *BlogService) GetHotFeed(loginUserId int64) ([]common.BlogVo, error) {
	blogList, err := blogDaoInstance.QueryHotFeedFlow()
	if err != nil {
		return nil, err
	}
	blogListResp := make([]common.BlogVo, len(blogList))
	fmt.Println("获取博文流成功！")
	for i, v := range blogList {
		var isFavorite bool
		user, err1 := userDaoInstance.QueryUserById(v.UserId)
		if err1 != nil {
			return nil, err1
		}
		actionType, err2 := favoriteDaoInstance.QueryActionTypeByUserIdAndBlogId(loginUserId, v.Id)
		if err2 != nil {
			return nil, err2
		}
		if actionType == 1 {
			isFavorite = true
		} else {
			isFavorite = false
		}
		favoriteCount, err3 := favoriteDaoInstance.QueryFavoriteCountByUserId(v.UserId)
		if err3 != nil {
			return nil, err3
		}
		var totalFavorited = int64(0)
		if count, err4 := blogDaoInstance.QueryPublishCountByUserId(v.UserId); err4 != nil {
			return nil, err4
		} else if count > 0 {
			totalFavorited, err = blogDaoInstance.QueryTotalFavoriteCountByUserId(v.UserId)
		}
		if err != nil {
			return nil, err
		}
		isFollow, err5 := relationDaoInstance.QueryIsFollowByUserIdAndToUserId(loginUserId, v.UserId)
		if err5 != nil {
			return nil, err5
		}
		tmpUser := &common.UserVo{
			Id:             user.Id,
			UserName:       user.UserName,
			FollowCount:    user.FollowCount,
			FollowerCount:  user.FollowerCount,
			IsFollow:       isFollow,
			Avatar:         user.Avatar,
			TotalFavorited: totalFavorited,
			FavoriteCount:  favoriteCount,
		}
		blogListResp[i] = common.BlogVo{
			Id:            v.Id,
			Author:        *tmpUser,
			Title:         v.Title,
			Content:       v.Content,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			Top:           v.Top,
			IsFavorite:    isFavorite,
			CreateTime:    v.CreateTime,
			UpdateTime:    v.UpdateTime,
		}
	}
	return blogListResp, err
}

// TopBlog 置顶个人的一篇博客
func (b *BlogService) TopBlog(blogId, loginUserId, actionType int64) error {
	err := blogDaoInstance.TopBlog(blogId, actionType)
	if err != nil {
		return err
	}
	key := "personal_feeds:" + strconv.FormatInt(loginUserId, 10)
	if actionType == 1 {
		now := time.Now().AddDate(50, 0, 0).UnixNano()
		// 置顶个人的缓存博客列表
		utils.ZAdd(key, now, blogId)
	} else {
		blog, err1 := blogDaoInstance.QueryBlogById(blogId)
		if err1 != nil {
			return err1
		}
		utils.ZAdd(key, blog.CreateTime.UnixNano(), blogId)
	}
	return nil
}

// GetRelationNewFeed 获取关注用户最新博客流(分页查询)  todo 存在重复查询问题
func (b *BlogService) GetRelationNewFeed(page, pageSize, loginUserId int64) ([]common.BlogVo, error) {
	if page < 0 {
		page = 1
	}
	key := "following_feeds:" + strconv.FormatInt(loginUserId, 10)
	start := (page - 1) * pageSize
	end := page*pageSize - 1
	tFeedIds := utils.ZRevRange(key, start, end)
	if tFeedIds == nil || len(tFeedIds) <= 0 {
		return nil, nil
	}
	feedIds := make([]int64, len(tFeedIds))
	//string转int64
	for i, v := range tFeedIds {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		feedIds[i] = id
	}
	fmt.Println(feedIds)
	blogList, err := blogDaoInstance.QueryByIds(feedIds)
	if err != nil {
		return nil, err
	}
	blogListResp := make([]common.BlogVo, len(blogList))
	fmt.Println("获取博文流成功！")
	for i, v := range blogList {
		var isFavorite bool
		user, err1 := userDaoInstance.QueryUserById(v.UserId)
		if err1 != nil {
			return nil, err1
		}
		actionType, err2 := favoriteDaoInstance.QueryActionTypeByUserIdAndBlogId(loginUserId, v.Id)
		if err2 != nil {
			return nil, err2
		}
		if actionType == 1 {
			isFavorite = true
		} else {
			isFavorite = false
		}
		favoriteCount, err3 := favoriteDaoInstance.QueryFavoriteCountByUserId(v.UserId)
		if err3 != nil {
			return nil, err3
		}
		var totalFavorited = int64(0)
		if count, err4 := blogDaoInstance.QueryPublishCountByUserId(v.UserId); err4 != nil {
			return nil, err4
		} else if count > 0 {
			totalFavorited, err = blogDaoInstance.QueryTotalFavoriteCountByUserId(v.UserId)
		}
		if err != nil {
			return nil, err
		}
		isFollow, err5 := relationDaoInstance.QueryIsFollowByUserIdAndToUserId(loginUserId, v.UserId)
		if err5 != nil {
			return nil, err5
		}
		tmpUser := &common.UserVo{
			Id:             user.Id,
			UserName:       user.UserName,
			FollowCount:    user.FollowCount,
			FollowerCount:  user.FollowerCount,
			IsFollow:       isFollow,
			Avatar:         user.Avatar,
			TotalFavorited: totalFavorited,
			FavoriteCount:  favoriteCount,
		}
		blogListResp[i] = common.BlogVo{
			Id:            v.Id,
			Author:        *tmpUser,
			Title:         v.Title,
			Content:       v.Content,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			Top:           v.Top,
			IsFavorite:    isFavorite,
			CreateTime:    v.CreateTime,
			UpdateTime:    v.UpdateTime,
		}
	}
	return blogListResp, err
}

// GetRelationNewFeed1 获取关注用户最新博客流(使用最后一个时间戳分页查询)
func (b *BlogService) GetRelationNewFeed1(max float64, offset, loginUserId int64) ([]common.BlogVo, float64, int64, error) {
	key := "following_feeds:" + strconv.FormatInt(loginUserId, 10)
	result := utils.ZRevByScoreWithScores(key, 0, max, offset, 5)
	if result == nil || len(result) <= 0 {
		return nil, 0, 0, nil
	}
	feedIds := make([]int64, len(result))
	//转int64
	var minTime float64 = 0
	var os int64 = 1
	for i, v := range result {
		id, err := strconv.ParseInt(v.Member.(string), 10, 64)
		if err != nil {
			return nil, 0, 0, err
		}
		if v.Score == minTime {
			os++
		} else {
			minTime = v.Score
			os = 1
		}
		feedIds[i] = id
	}
	fmt.Println(feedIds)
	blogList, err := blogDaoInstance.QueryByIds(feedIds)
	if err != nil {
		return nil, 0, 0, nil
	}
	blogListResp := make([]common.BlogVo, len(blogList))
	fmt.Println("获取博文流成功！")
	for i, v := range blogList {
		var isFavorite bool
		user, err1 := userDaoInstance.QueryUserById(v.UserId)
		if err1 != nil {
			return nil, 0, 0, err1
		}
		actionType, err2 := favoriteDaoInstance.QueryActionTypeByUserIdAndBlogId(loginUserId, v.Id)
		if err2 != nil {
			return nil, 0, 0, err2
		}
		if actionType == 1 {
			isFavorite = true
		} else {
			isFavorite = false
		}
		favoriteCount, err3 := favoriteDaoInstance.QueryFavoriteCountByUserId(v.UserId)
		if err3 != nil {
			return nil, 0, 0, err3
		}
		var totalFavorited = int64(0)
		if count, err4 := blogDaoInstance.QueryPublishCountByUserId(v.UserId); err4 != nil {
			return nil, 0, 0, err4
		} else if count > 0 {
			totalFavorited, err = blogDaoInstance.QueryTotalFavoriteCountByUserId(v.UserId)
		}
		if err != nil {
			return nil, 0, 0, err
		}
		isFollow, err5 := relationDaoInstance.QueryIsFollowByUserIdAndToUserId(loginUserId, v.UserId)
		if err5 != nil {
			return nil, 0, 0, err5
		}
		tmpUser := &common.UserVo{
			Id:             user.Id,
			UserName:       user.UserName,
			FollowCount:    user.FollowCount,
			FollowerCount:  user.FollowerCount,
			IsFollow:       isFollow,
			Avatar:         user.Avatar,
			TotalFavorited: totalFavorited,
			FavoriteCount:  favoriteCount,
		}
		blogListResp[i] = common.BlogVo{
			Id:            v.Id,
			Author:        *tmpUser,
			Title:         v.Title,
			Content:       v.Content,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			Top:           v.Top,
			IsFavorite:    isFavorite,
			CreateTime:    v.CreateTime,
			UpdateTime:    v.UpdateTime,
		}
	}
	return blogListResp, minTime, os, err
}
