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

	gin.SetMode(gin.ReleaseMode)

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

	router := gin.New()

	//注册用户相关的路由

	router.Use(middleware.RecoveryMiddleware(), middleware.LoggerMiddleware()) //使用日志与恢复中间件

	public := router.Group("/api") //公开路由组
	{
		public.POST("/register", userhandler.RegisterUser)
		public.POST("/login", userhandler.LoginUser)
	}
	private := router.Group("/api") //私有路由组
	private.Use(middleware.AuthMiddleware())
	{
		private.POST("/delete", userhandler.DeleteUser)
		private.POST("/change_password", userhandler.ChangePassword)
		private.GET("/users", userhandler.GetUser)
	}
	//启动服务器
	if err := router.Run(ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
