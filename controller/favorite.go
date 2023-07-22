package controller

import (
	"feeddemo/common"
	"feeddemo/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type FavoriteResponse struct {
	common.Response
	BlogList []common.BlogVo `json:"blog_list,omitempty"`
}

var favoriteService = service.NewFavoriteServiceInstance()

// FavoriteAction 点赞
func FavoriteAction(c *gin.Context) {
	userId, err1 := strconv.ParseInt(c.Query("user_id"), 10, 64)
	blogId, err2 := strconv.ParseInt(c.Query("blog_id"), 10, 64)
	actionType, err3 := strconv.ParseInt(c.Query("action_type"), 10, 32)
	if err1 != nil || err2 != nil || err3 != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 500,
			StatusMsg:  "服务端错误！",
		})
		return
	}
	//调用service层
	err := favoriteService.DoFavoriteAction(userId, blogId, int32(actionType))
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 500,
			StatusMsg:  "服务端错误！",
		})
		return
	}
	c.JSON(http.StatusOK,
		common.Response{
			StatusCode: 200,
			StatusMsg:  "更新点赞状态成功！",
		})
	return
}

// FavoriteList 获取指定用户点赞列表
func FavoriteList(c *gin.Context) {
	userId, err1 := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err1 != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 500,
			StatusMsg:  "服务端错误！",
		})
		return
	}
	//todo 当前登录用户ID，先写固定值
	var loginUserId int64 = 1
	//调用service层
	blogListResp, err := favoriteService.DoFavoriteList(userId, loginUserId)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 500,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, FavoriteResponse{
		Response: common.Response{StatusMsg: "获取点赞博文列表成功！"},
		BlogList: blogListResp,
	})
	return
}
