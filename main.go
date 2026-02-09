package main

import (
	"log"
	"user_system/database"
	"user_system/middleware"
	"user_system/userhandler"

	"github.com/gin-gonic/gin"
)

func main() {

	ServerPort := ":8080" //默认端口
	//初始化数据库连接
	err := database.InitDB()
	if err != nil {
		log.Fatalf("%v", err)
		log.Printf("Failed to initialize database: %v", err)
		panic(err)
	}
	//确保在程序结束时关闭数据库连接
	defer func() {
		if err := database.CloseDB(); err != nil {
			log.Fatalf("%v", err)
			panic(err)
		}
	}()
	err = userhandler.Init() //初始化数据库连接
	if err != nil {
		log.Fatalf("%v", err)
		panic(err)
	}
	//创建Gin路由

	router := gin.Default()

	//注册用户相关的路由

	router.Use(middleware.LoggerMiddleware()) //使用日志中间件

	router.POST("/api/register", userhandler.RegisterUser)
	router.POST("/api/login", userhandler.LoginUser)
	//router.POST("/api/delete", userhandler.DeleteUser) 等验证中间件写完了再开放
	//router.POST("/api/change_password", userhandler.ChangePassword) 等验证中间件写完了再开放

	//启动服务器
	if err := router.Run(ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
