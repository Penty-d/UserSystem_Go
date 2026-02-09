package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//在处理请求前打印日志
		method := c.Request.Method
		path := c.Request.URL.Path
		IP := c.ClientIP()
		log.Printf("Received %s request for %s from %s", method, path, IP)
		//处理请求
		c.Next()
		//在处理请求后打印日志
		status := c.Writer.Status()
		message := c.GetString("message") //获取响应消息
		log.Printf("Responded %s : %s with status %d for %s %s", IP, message, status, method, path)
	}
}

/* 咕，有考虑写token认证的中间件，但懒得写了，等有时间再写吧
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//在处理请求前检查Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"message": "Unauthorized: Missing Authorization header"})
			c.Abort()
			return
		}
*/
