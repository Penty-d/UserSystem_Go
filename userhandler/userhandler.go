package userhandler

import (
	"fmt"
	"strconv"
	"user_system/models"
	"user_system/repositories"
	"user_system/utils"

	"github.com/gin-gonic/gin"
)

func Init() error {
	//初始化数据库连接
	err := repositories.NewDBHandler()
	if err != nil {
		return err
	}
	return utils.NewAuthDBHandler()
}

func SendResponse(c *gin.Context, status int, message string) { //规范返回
	if status == 200 {
		c.Set("message", message)
		c.JSON(status, gin.H{"message": message})
		return
	}
	c.Set("message", message)
	c.JSON(status, gin.H{"error": "Failed to execute request"})
}

func RegisterUser(c *gin.Context) {
	var userInfo models.CreateUserRequest
	err := c.ShouldBindJSON(&userInfo)
	if err != nil {
		SendResponse(c, 400, err.Error())
		return
	}
	response := repositories.CreateUser(&userInfo)
	SendResponse(c, response.Type, response.Message)
}

func LoginUser(c *gin.Context) {
	var userInfo models.LoginRequest
	err := c.ShouldBindJSON(&userInfo)
	if err != nil {
		SendResponse(c, 400, err.Error())
		return
	}
	response, token := repositories.UserLogin(&userInfo)
	if token == "" {
		SendResponse(c, response.Type, response.Message)
		return
	}
	c.Set("message", response.Message)
	if response.Type == 200 {
		c.JSON(response.Type, gin.H{"message": response.Message, "token": token})
		return
	}
	c.JSON(response.Type, gin.H{"error": "Failed to execute request"})
}

func DeleteUser(c *gin.Context) {
	info, exist := c.Get("info")
	if !exist {
		SendResponse(c, 400, "Failed to get info by token")
		return
	}
	if info.(*utils.TokenInfo).Role != "admin" {
		SendResponse(c, 400, fmt.Sprintf("Failed to delete user,%s", info.(*utils.TokenInfo).Role))
		return
	}
	var userInfo models.UpdateUserRequest
	err := c.ShouldBindJSON(&userInfo)
	if err != nil {
		SendResponse(c, 400, err.Error())
		return
	}
	status := "deleted"
	userInfo.Status = &status
	response := repositories.UpdateUser(&userInfo)
	SendResponse(c, response.Type, response.Message)
}

func ChangePassword(c *gin.Context) {
	var userInfo models.UpdateUserRequest
	err := c.ShouldBindJSON(&userInfo)
	if err != nil {
		SendResponse(c, 400, err.Error())
		return
	}
	info, exist := c.Get("info")
	if !exist {
		SendResponse(c, 400, "Failed to get info by token")
		return
	}
	if info.(*utils.TokenInfo).Username != userInfo.Username || info.(*utils.TokenInfo).Role != "admin" {
		SendResponse(c, 400, "Failed to change password")
		return
	}
	response := repositories.UpdateUser(&userInfo)
	SendResponse(c, response.Type, response.Message)
}

func GetUser(c *gin.Context) {
	info, exist := c.Get("info")
	if !exist {
		SendResponse(c, 400, "Failed to get info by token")
		return
	}
	if info.(*utils.TokenInfo).Role != "admin" {
		SendResponse(c, 400, fmt.Sprintf("Failed to get user,%s", info.(*utils.TokenInfo).Role))
		return
	}
	Username := c.Query("username")

	if Username != "" {
		userInfo, response := repositories.GetUserByUsername(Username)
		c.Set("message", response.Message)
		c.JSON(response.Type, gin.H{"message": response.Message, "user": userInfo})
		return
	}
	ID, err := strconv.ParseUint(c.Query("id"), 10, 64)
	if err != nil {
		SendResponse(c, 400, "Failed to get user")
		return
	}
	if ID != 0 {
		ID := uint(ID)
		userInfo, response := repositories.GetUserInfoByID(ID)
		c.Set("message", response.Message)
		c.JSON(response.Type, gin.H{"message": response.Message, "user": userInfo})
		return
	}
	SendResponse(c, 400, "Failed to get user")
}
