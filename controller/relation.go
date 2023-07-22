package controller

import (
	"feeddemo/common"
	"feeddemo/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type RelationListResponse struct {
	common.Response
	UserList []common.UserVo `json:"user_list,omitempty"`
}

var relationService = service.NewRelationServiceInstance()

// RelationAction 关注
func RelationAction(c *gin.Context) {
	loginUserId, err0 := strconv.ParseInt(c.Query("login_user_id"), 10, 64)
	toUserId, err1 := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	actionType, err2 := strconv.ParseInt(c.Query("action_type"), 10, 32)
	if err0 != nil || err1 != nil || err2 != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 500,
			StatusMsg:  "服务端错误，评论操作失败！",
		})
		return
	}
	// 调用service层
	err := relationService.DoRelationAction(loginUserId, toUserId, actionType)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 500,
			StatusMsg:  "服务端错误，修改关注状态失败！",
		})
		return
	}
	c.JSON(http.StatusOK, common.Response{
		StatusCode: 200,
		StatusMsg:  "修改关注状态成功！",
	})
	return
}

// RelationFollowList 获取关注列表
func RelationFollowList(c *gin.Context) {
	userId, err1 := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err1 != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 500,
			StatusMsg:  "服务端错误，评论操作失败！",
		})
		return
	}
	//todo 当前登录用户ID，现在先默认值
	loginUserId := int64(1)
	// 调用service层
	userList, err := relationService.DoRelationFollowList(userId, loginUserId)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 500,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, RelationListResponse{
		Response: common.Response{StatusMsg: "获取关注列表成功！"},
		UserList: userList,
	})
	return
}

// RelationFollowerList 获取粉丝列表
func RelationFollowerList(c *gin.Context) {
	userId, err1 := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err1 != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 500,
			StatusMsg:  "服务端错误，评论操作失败！",
		})
		return
	}
	//todo 当前登录用户ID，现在先默认值
	loginUserId := int64(1)
	// 调用service层
	userList, err := relationService.DoRelationFollowerList(userId, loginUserId)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 500,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, RelationListResponse{
		Response: common.Response{StatusMsg: "获取粉丝列表成功！"},
		UserList: userList,
	})
	return
}
