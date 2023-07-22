package controller

import (
	"feeddemo/common"
	"feeddemo/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var commentService = service.NewCommentServiceInstance()

type CommentResponse struct {
	common.Response
	Comment common.CommentVo `json:"comment"`
}

type CommentListResponse struct {
	common.Response
	CommentList []common.CommentVo `json:"comment_list,omitempty"`
}

// CommentAction 评论功能
func CommentAction(c *gin.Context) {
	blogId, err1 := strconv.ParseInt(c.Query("blog_id"), 10, 64)
	actionType, err2 := strconv.ParseInt(c.Query("action_type"), 10, 32)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 500,
			StatusMsg:  "服务端错误，评论操作失败！",
		})
		return
	}
	//todo 当前登录用户ID，先写固定值
	var loginUserId int64 = 1
	fmt.Printf("评论userId：%d==blogId：%d==actionType:%d\n", loginUserId, blogId, actionType)
	if actionType == 1 {
		//新增评论
		commentText := c.Query("comment_text")
		//调用service层
		commentVo, err := commentService.DoAddCommentAction(loginUserId, blogId, commentText)
		if err != nil {
			c.JSON(http.StatusOK, common.Response{
				StatusCode: 500,
				StatusMsg:  err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, CommentResponse{
			Response: common.Response{
				StatusCode: 200,
				StatusMsg:  "新增评论成功！",
			},
			Comment: *commentVo,
		})
	} else {
		//删除评论
		commentId, err4 := strconv.ParseInt(c.Query("comment_id"), 10, 64)
		if err4 != nil {
			c.JSON(http.StatusOK, common.Response{
				StatusCode: 500,
				StatusMsg:  "服务端错误，评论操作失败！",
			})
			return
		}
		//调用service层
		err5 := commentService.DoDelCommentAction(loginUserId, blogId, commentId)
		if err5 != nil {
			c.JSON(http.StatusOK, common.Response{
				StatusCode: 500,
				StatusMsg:  err5.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 200,
			StatusMsg:  "删除评论成功！",
		})
	}
	return
}

// CommentList 获取评论列表
func CommentList(c *gin.Context) {
	blogId, err1 := strconv.ParseInt(c.Query("blog_id"), 10, 64)
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
	commentList, err := commentService.DoGetCommentList(loginUserId, blogId)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 500,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response: common.Response{
			StatusCode: 200,
			StatusMsg:  "获取博文评论列表成功！",
		},
		CommentList: commentList,
	})
	return
}
