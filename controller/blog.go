package controller

import (
	"feeddemo/common"
	"feeddemo/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type FeedResponse struct {
	common.Response
	BlogList    []common.BlogVo `json:"blog_list,omitempty"`
	LastMinTime float64         `json:"last_min_time,omitempty"`
	LastMaxTime int64           `json:"last_max_time,omitempty"`
	Offset      int64           `json:"offset,omitempty"`
}

var blogService = service.NewBlogServiceInstance()

// PublishBlog 发布博文
func PublishBlog(c *gin.Context) {
	// todo 应该为登录用户id
	userId, err1 := strconv.ParseInt(c.PostForm("user_id"), 10, 64)
	if err1 != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 500,
			StatusMsg:  "获取用户ID失败！",
		})
		return
	}
	title := c.PostForm("title")
	content := c.PostForm("content")
	//调用service层
	err := blogService.DoPublishBlog(title, content, userId)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 500,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, common.Response{
		StatusCode: 200,
		StatusMsg:  "博文上传成功！！！！",
	})
	return
}

// DeleteBlog 删除博文
func DeleteBlog(c *gin.Context) {
	// todo 应该为登录用户id
	userId, err1 := strconv.ParseInt(c.Query("user_id"), 10, 64)
	blogId, err2 := strconv.ParseInt(c.Query("blog_id"), 10, 64)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 500,
			StatusMsg:  "数据转换失败！",
		})
		return
	}
	//调用service层
	err := blogService.DoDelBlog(userId, blogId)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 500,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, common.Response{
		StatusCode: 200,
		StatusMsg:  "博文删除成功！！！！",
	})
	return
}

// GetNewFeeds 获取平台最新30条博客流
func GetNewFeeds(c *gin.Context) {
	lastMaxTime, err := strconv.ParseInt(c.Query("last_max_time"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 500,
			StatusMsg:  "数据转换错误！",
		})
		return
	}
	//todo 当前登录用户ID，这里先固定值
	loginUserId := int64(1)
	// 调用service层
	blogListResp, maxTime, err := blogService.GetNewFeed(loginUserId, lastMaxTime)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 500,
			StatusMsg:  "服务端错误！",
		})
		return
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response:    common.Response{StatusCode: 200, StatusMsg: "获取最新博文流成功！"},
		BlogList:    blogListResp,
		LastMaxTime: maxTime,
	})
	return
}

// GetHotFeeds 获取平台最火热30条博客流
func GetHotFeeds(c *gin.Context) {
	//todo 当前登录用户ID，这里先固定值
	loginUserId := int64(1)
	// 调用service层
	blogListResp, err := blogService.GetHotFeed(loginUserId)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 500,
			StatusMsg:  "服务端错误！",
		})
		return
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response: common.Response{StatusCode: 200, StatusMsg: "获取最热博文流成功！"},
		BlogList: blogListResp,
	})
	return
}

// GetUserPublishList 获取指定用户发布博文列表 置顶+发布时间排序
func GetUserPublishList(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	lastMinTime, err1 := strconv.ParseFloat(c.Query("last_min_time"), 64)
	offset, err2 := strconv.ParseInt(c.Query("offset"), 10, 64)
	if err != nil || err1 != nil || err2 != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 500,
			StatusMsg:  "数据转换错误！",
		})
		return
	}
	//todo 当前登录用户ID，这里先固定值
	var loginUserId int64 = 1
	//调用service层
	PublishedList, minTime, os, err1 := blogService.GetBlogListByUserID(lastMinTime, offset, userId, loginUserId)
	if err1 != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 500,
			StatusMsg:  "服务端错误！",
		})
		return
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response:    common.Response{StatusCode: 200, StatusMsg: "获取当前用户发布的博文列表成功！"},
		BlogList:    PublishedList,
		LastMinTime: minTime,
		Offset:      os,
	})
	return
}

// TopBlog 置顶个人的一篇博客
func TopBlog(c *gin.Context) {
	blogId, err1 := strconv.ParseInt(c.Query("blog_id"), 10, 64)
	actionType, err2 := strconv.ParseInt(c.Query("action_type"), 10, 64)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 500,
			StatusMsg:  "数据转换错误！",
		})
		return
	}
	//todo 当前登录用户ID，这里先固定值
	var loginUserId int64 = 1
	err := blogService.TopBlog(blogId, loginUserId, actionType)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 500,
			StatusMsg:  "服务端错误！",
		})
		return
	}
	c.JSON(http.StatusOK, common.Response{
		StatusCode: 200,
		StatusMsg:  "修改博文置顶状态成功！",
	})
	return
}

// GetRelationFeedList 获取关注用户最新博客流(分页查询) //重复问题
func GetRelationFeedList(c *gin.Context) {
	page, err1 := strconv.ParseInt(c.Query("page"), 10, 64)
	pageSize, err2 := strconv.ParseInt(c.Query("page_size"), 10, 64)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 500,
			StatusMsg:  "数据转换错误！",
		})
		return
	}
	//todo 当前登录用户ID，这里先固定值
	loginUserId := int64(1)
	// 调用service层
	blogListResp, err := blogService.GetRelationNewFeed(page, pageSize, loginUserId)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 500,
			StatusMsg:  "服务端错误！",
		})
		return
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response: common.Response{StatusCode: 200, StatusMsg: "获取关注用户最新博文流成功！"},
		BlogList: blogListResp,
	})
	return
}

// GetRelationFeedList1 获取关注用户最新博客流(滚动分页查询)
func GetRelationFeedList1(c *gin.Context) {
	lastMinTime, err1 := strconv.ParseFloat(c.Query("last_min_time"), 64)
	offset, err2 := strconv.ParseInt(c.Query("offset"), 10, 64)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 500,
			StatusMsg:  "数据转换错误！",
		})
		return
	}
	//todo 当前登录用户ID，这里先固定值
	loginUserId := int64(1)
	// 调用service层
	blogListResp, minTime, offset, err := blogService.GetRelationNewFeed1(lastMinTime, offset, loginUserId)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 500,
			StatusMsg:  "服务端错误！",
		})
		return
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response:    common.Response{StatusCode: 200, StatusMsg: "获取关注用户最新博文流成功！"},
		BlogList:    blogListResp,
		LastMinTime: minTime,
		Offset:      offset,
	})
	return
}
