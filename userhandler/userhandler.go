package userhandler

import (
	"user_system/models"
	"user_system/repositories"

	"github.com/gin-gonic/gin"
)

func Init() error {
	//初始化数据库连接
	err := repositories.NewDBHandler()
	return err
}

func SendJsonResponse(c *gin.Context, status int, message string) {
	c.Set("message", message)
	c.JSON(status, gin.H{"message": message})
}

func RegisterUser(c *gin.Context) {
	var userInfo models.UserInfo
	err := c.ShouldBindJSON(&userInfo)
	if err != nil {
		SendJsonResponse(c, 400, err.Error())
		return
	}
	response := repositories.CreateUser(&userInfo)
	SendJsonResponse(c, response.Type, response.Message)
}

func LoginUser(c *gin.Context) {
	var userInfo models.UserInfo
	err := c.ShouldBindJSON(&userInfo)
	if err != nil {
		SendJsonResponse(c, 400, err.Error())
		return
	}
	response := repositories.UserLogin(&userInfo)
	SendJsonResponse(c, response.Type, response.Message)
}

func DeleteUser(c *gin.Context) {
	var userInfo models.UserInfo
	err := c.ShouldBindJSON(&userInfo)
	if err != nil {
		SendJsonResponse(c, 400, err.Error())
		return
	}
	response := repositories.DeleteUser(&userInfo)
	SendJsonResponse(c, response.Type, response.Message)
}

func ChangePassword(c *gin.Context) {
	var userInfo models.UserInfo
	err := c.ShouldBindJSON(&userInfo)
	if err != nil {
		SendJsonResponse(c, 400, err.Error())
		return
	}
	response := repositories.ChangePassword(&userInfo)
	SendJsonResponse(c, response.Type, response.Message)
}
