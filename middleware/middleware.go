package middleware

import (
	"log"
	"strings"
	"time"
	"user_system/utils"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//在处理请求前打印日志
		method := c.Request.Method
		path := c.Request.URL.Path
		IP := c.ClientIP()
		ReceivedTime := time.Now()
		log.Printf("%s | Received  | %s | %s | %s", ReceivedTime.Format("2006-01-02 15:04:05"), method, path, IP)
		//处理请求
		c.Next()
		//在处理请求后打印日志
		status := c.Writer.Status()
		message := c.GetString("message") //获取响应消息
		ResponedTime := time.Now()
		log.Printf("%s | Responded | %s | %s | %s | %s | %d | %s", ResponedTime.Format("2006-01-02 15:04:05"), method, path, IP, ResponedTime.Sub(ReceivedTime), status, message)
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//在处理请求前检查Authorization头
		authHeader := c.GetHeader("Authorization")

		if len(authHeader) < 7 || strings.ToLower(authHeader[:7]) != "bearer " {
			c.Set("message", "Unauthorized: Missing Authorization header")
			c.JSON(401, gin.H{"message": "Unauthorized: Missing Authorization header"})
			c.Abort()
		}
		token := authHeader[7:]
		info, err := utils.GetInfobyToken(token)
		if err != nil {
			c.Set("message", err.Error())
			c.JSON(401, gin.H{"message": err.Error()})
			c.Abort()
		} else if info.ExpiredAt.Before(time.Now()) {
			c.Set("message", "Unauthorized: Token expired")
			c.JSON(401, gin.H{"message": "Unauthorized: Invalid token"})
			c.Abort()
		}
		c.Set("info", info)
		c.Next()
	}
}

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recovered from panic: %v", err)
				c.JSON(500, gin.H{"message": "Internal Server Error"})
				c.Abort()
			}
		}()
		c.Next()
	}
}
