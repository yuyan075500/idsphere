package dao

import (
	"errors"
	"github.com/wonderivan/logger"
	"ops-api/config"
	"ops-api/global"
	"ops-api/model"
	"ops-api/utils"
	"time"
)

var User user

type user struct{}

// UserList 返回给前端的结构体
type UserList struct {
	Items []*UserInfo `json:"items"`
	Total int64       `json:"total"`
}

// UserInfo 返回的用户字段信息
type UserInfo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Username    string `json:"username"`
	PhoneNumber string `json:"phone_number"`
	IsActive    int    `json:"is_active"`
	Email       string `json:"email"`
	Avatar      string `json:"avatar"`
}

// GetUserList 获取用户列表
func (u *user) GetUserList(name string, page, limit int) (data *UserList, err error) {
	// 定义数据的起始位置
	startSet := (page - 1) * limit

	// 定义返回的内容
	var (
		userList []*UserInfo
		total    int64
	)

	// 获取用户列表
	tx := global.MySQLClient.Model(&model.AuthUser{}).
		Where("name like ?", "%"+name+"%"). // 实现过滤
		Count(&total).                      // 获取总数
		Limit(limit).
		Offset(startSet).
		Find(&userList)
	if tx.Error != nil {
		logger.Error("获取用户列表失败：", tx.Error)
		return nil, errors.New("获取用户列表失败：" + tx.Error.Error())
	}

	return &UserList{
		Items: userList,
		Total: total,
	}, nil
}

// GetUser 获取用户信息
func (u *user) GetUser(userid uint) (user *UserInfo, err error) {

	var userInfo *UserInfo

	tx := global.MySQLClient.Model(&model.AuthUser{}).Where("id = ?", userid).Find(&userInfo)
	if tx.Error != nil {
		logger.Error("获取用户信息失败：", tx.Error)
		return nil, errors.New("获取用户信息失败：" + tx.Error.Error())
	}

	// 从OSS中获取头像临时访问URL，临时URL的过期时间与用户Token过期时间保持一致
	avatarURL, err := utils.GetPresignedURL(userInfo.Avatar, time.Duration(config.Conf.JWT.Expires)*time.Hour)
	if err != nil {
		logger.Error("获取用户头像失败：", err.Error())
		userInfo.Avatar = ""
		return userInfo, nil
	}

	userInfo.Avatar = avatarURL.String()
	return userInfo, nil
}

// AddUser 新增用户
func (u *user) AddUser(data *model.AuthUser) (err error) {
	tx := global.MySQLClient.Create(&data)
	if tx.Error != nil {
		logger.Error("新增用户失败：", tx.Error)
		return errors.New("新增用户失败：" + tx.Error.Error())
	}
	return nil
}

// UpdateUser 更新用户信息
//func (u *user) UpdateUser(data map[string]interface{}) (err error) {
//	tx := db.GORM.Model(&model.AuthUser{}).Updates(data)
//	if tx.Error != nil {
//		logger.Error("更新用户信息失败：", tx.Error)
//		return errors.New("更新用户信息失败：" + tx.Error.Error())
//	}
//	return nil
//}