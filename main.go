package main

import (
	"feeddemo/config"
	"feeddemo/repository"
	"feeddemo/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	r := gin.Default()
	config.InitRouter(r)
	if err := repository.Init(); err != nil {
		panic(err)
	}
	if err := utils.RedisInit(); err != nil {
		fmt.Printf("redis连接失败! err : %v\n", err)
		return
	}
	fmt.Println("redis连接成功！")
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	//可以添加多个搜索路径  第一个找不到会找后面的
	viper.AddConfigPath("./config")
	//读取内容
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	port := viper.GetString("server.port")
	fmt.Println(port)
	r.Run(":" + port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
