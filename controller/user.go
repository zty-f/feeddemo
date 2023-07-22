package controller

import (
	"feeddemo/common"
	"feeddemo/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

var userService = service.NewUserServiceInstance()

// Register 注册
func Register(c *gin.Context) {
	userName := c.Query("username")
	password := c.Query("password")
	//调用Service层
	var err = userService.DoRegister(userName, password)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 500,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, common.Response{
		StatusCode: 200,
		StatusMsg:  "注册新用户成功！",
	})
	return
}
