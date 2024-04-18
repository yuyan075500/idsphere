package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"net/http"
	"ops-api/config"
	"ops-api/dao"
	"ops-api/global"
	"ops-api/middleware"
	"ops-api/model"
	"ops-api/service"
	"ops-api/utils"
	"path/filepath"
	"strings"
	"time"
)

var User user

type user struct{}

// Login 用户登录
func (u *user) Login(c *gin.Context) {
	var (
		user model.AuthUser
		err  error
	)

	params := new(struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	})

	if err = c.Bind(params); err != nil {
		logger.Error("无效的请求参数：" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 4000,
			"msg":  "无效的请求参数",
		})
		return
	}

	// 根据用户名查询用户
	if err := global.MySQLClient.Where("username = ?", params.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 4404,
			"msg":  "用户不存在",
		})
		return
	}

	// 判断用户是否禁用
	if user.IsActive == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 4403,
			"msg":  "用户未激活，请联系管理员",
		})
		return
	}

	// 检查密码
	if user.CheckPassword(params.Password) == false {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 4404,
			"msg":  "用户密码错误",
		})
		return
	}

	token, err := middleware.GenerateJWT(user.ID, user.Name, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 4000,
			"msg":  "生成Token错误",
		})
		return
	}

	// 记录用户最后登录时间
	global.MySQLClient.Model(&user).Where("id = ?", user.ID).Update("last_login_at", time.Now())

	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"msg":   "认证成功",
		"token": token,
	})
}

// Logout 用户注销
func (u *user) Logout(c *gin.Context) {
	// 获取Token
	token := c.Request.Header.Get("Authorization")
	parts := strings.SplitN(token, " ", 2)

	// 将Token存入Redis缓存
	err := global.RedisClient.Set(parts[1], true, time.Duration(config.Conf.JWT.Expires)*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 4000,
			"msg":  "用户注销失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "注销成功",
	})
}

// UploadAvatar 用户头像上传
func (u *user) UploadAvatar(c *gin.Context) {
	// 获取上传的头像
	avatar, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 4000,
			"msg":  "无效的请求参数",
		})
		return
	}

	// 打开上传头像
	src, err := avatar.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 4000,
			"msg":  "文件打开失败",
		})
		return
	}

	// 上传头像
	// 获取当前登录用户的用户名
	username, _ := c.Get("username")
	// 拼接头像存储的路径和文件名：avatar/<用户名><文件后缀>
	avatarName := fmt.Sprintf("avatar/%v%v", username, filepath.Ext(avatar.Filename))
	err = utils.FileUpload(avatarName, avatar.Header.Get("Content-Type"), src, avatar.Size)
	if err != nil {
		logger.Error("文件上传失败：" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 4000,
			"msg":  "文件上传失败",
		})
		return
	}

	// 将头像地址存储到数据库
	var user model.AuthUser
	global.MySQLClient.Model(&user).Where("username = ?", username).Update("avatar", avatarName)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "头像更新成功",
	})
}

// GetUser 获取用户信息
func (u *user) GetUser(c *gin.Context) {
	params := new(struct {
		Token string `form:"token" binding:"required"`
	})
	if err := c.Bind(params); err != nil {
		logger.Error("无效的请求参数：" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 4000,
			"msg":  "无效的请求参数",
		})
		return
	}

	// 从Token中获取用户ID
	mc, err := middleware.ParseToken(params.Token)
	if err != nil {
		logger.Error("无效的Token：", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 90401,
			"msg":  "无效的Token",
		})
		return
	}

	// 根据ID获取用户信息
	data, err := dao.User.GetUser(mc.ID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 90404,
			"msg":  "获取用户信息失败",
		})
		return
	}

	// 返回用户信息
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "获取用户信息成功",
		"data": data,
	})
}

// GetUserList 获取用户列表
func (u *user) GetUserList(c *gin.Context) {
	params := new(struct {
		Name  string `form:"name"`
		Page  int    `form:"page" binding:"required"`
		Limit int    `form:"limit" binding:"required"`
	})
	if err := c.Bind(params); err != nil {
		logger.Error("无效的请求参数：" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 4000,
			"msg":  "无效的请求参数",
		})
		return
	}

	data, err := service.User.GetUserList(params.Name, params.Page, params.Limit)
	if err != nil {
		logger.Error("获取用户列表失败：" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 4000,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": data,
	})
}

// AddUser 创建用户
func (u *user) AddUser(c *gin.Context) {
	var (
		user = &service.UserCreate{}
		err  error
	)

	if err = c.ShouldBind(user); err != nil {
		logger.Error("无效的请求参数：" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 4000,
			"msg":  err.Error(),
		})
		return
	}

	if err = service.User.AddUser(user); err != nil {
		logger.Error("创建用户失败：" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 4000,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "创建用户成功",
	})
}