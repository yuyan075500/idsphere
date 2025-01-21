package dao

import (
	"gorm.io/gorm"
	"ops-api/config"
	"ops-api/global"
	"ops-api/model"
	"ops-api/utils"
	"time"
)

var User user

type user struct{}

// UserList 返回给前端表格的数据结构体
type UserList struct {
	Items []*UserInfo `json:"items"`
	Total int64       `json:"total"`
}

// UserInfoWithMenu 用户信息结构体，用于用户登录后获取用户信息
type UserInfoWithMenu struct {
	UserInfo
	Menus       []*MenuItem `json:"menus"`
	Roles       []string    `json:"roles"`
	Permissions []string    `json:"permissions"`
}

// UserInfo 用户信息结构体
type UserInfo struct {
	ID                int        `json:"id"`
	Name              string     `json:"name"`
	Username          string     `json:"username"`
	WwId              string     `json:"ww_id"`
	CtyunId           string     `json:"ctyun_id"`
	PhoneNumber       string     `json:"phone_number"`
	IsActive          bool       `json:"is_active"`
	Email             string     `json:"email"`
	Avatar            string     `json:"avatar"`
	LastLoginAt       *time.Time `json:"last_login_at"`
	PasswordExpiredAt *time.Time `json:"password_expired_at"`
	UserFrom          string     `json:"user_from"`
}

// UserListAll 返回给前端下拉框或穿梭框的数据结构体
type UserListAll struct {
	Users []*UserBasicInfo `json:"users"`
}

// UserBasicInfo 用户基本信息结构体
type UserBasicInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// PasswordExpiredUserList 密码过期的用户结构体
type PasswordExpiredUserList struct {
	Name              string    `json:"name"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	PasswordExpiredAt time.Time `json:"password_expired_at"`
}

// UserUpdate 更新构体，定义更新时的字段信息
type UserUpdate struct {
	ID          uint    `json:"id" binding:"required"`
	WwId        *string `json:"ww_id"`
	CtyunId     *string `json:"ctyun_id"`
	PhoneNumber string  `json:"phone_number" validate:"omitempty,phone"`
	Email       string  `json:"email" validate:"omitempty,email"`
	IsActive    bool    `json:"is_active" validate:"omitempty"`
}

// UserPasswordUpdate 更改密码结构体
type UserPasswordUpdate struct {
	ID         uint   `json:"id" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required"`
}

// UserCreate 创建结构体，定义新增时的字段信息
type UserCreate struct {
	Name        string `json:"name" binding:"required"`
	Username    string `json:"username" gorm:"unique" binding:"required"`
	Password    string `json:"password" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required" validate:"phone"`
	Email       string `json:"email" binding:"required" validate:"email"`
	UserFrom    string `json:"user_from"`
}

// GetUserListAll 获取所有用户
func (u *user) GetUserListAll() (data *UserListAll, err error) {

	var userBasicInfo []*UserBasicInfo

	// 获取用户列表
	if err := global.MySQLClient.Model(&model.AuthUser{}).
		Select("id, name").
		Find(&userBasicInfo).Error; err != nil {
		return nil, err
	}

	return &UserListAll{
		Users: userBasicInfo,
	}, nil
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
		Where("name like ? OR username like ? OR phone_number like ? OR email like ?", "%"+name+"%", "%"+name+"%", "%"+name+"%", "%"+name+"%"). // 实现过滤
		Count(&total).                                                                                                                          // 获取总数
		Limit(limit).
		Offset(startSet).
		Find(&userList)
	if tx.Error != nil {
		return nil, err
	}

	return &UserList{
		Items: userList,
		Total: total,
	}, nil
}

// GetUser 获取用户信息（动态查找，返回单个用户信息）
func (u *user) GetUser(conditions interface{}) (*model.AuthUser, error) {

	var user model.AuthUser

	if err := global.MySQLClient.Where(conditions).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserInfo 获取用户信息
func (u *user) GetUserInfo(userid uint) (userinfo *UserInfoWithMenu, err error) {

	var (
		userInfo *UserInfo
		user     model.AuthUser
		roles    []string
	)

	// 开启事务
	tx := global.MySQLClient.Begin()

	// 获取用户信息
	if err := tx.Model(&model.AuthUser{}).Where("id = ?", userid).Find(&userInfo).Error; err != nil {
		return nil, err
	}

	// 从OSS中获取头像临时访问URL，临时URL的过期时间与用户Token过期时间保持一致
	avatarURL, err := utils.GetPresignedURL(userInfo.Avatar, 6*time.Hour)
	if err != nil {
		userInfo.Avatar = ""
	} else {
		userInfo.Avatar = avatarURL.String()
	}

	// 获取用户菜单
	menus, err := Menu.GetUserMenu(tx, userInfo.Username)
	if err != nil {
		return nil, err
	}

	// 获取用户角色
	err = tx.Preload("Groups", "is_role_group = ?", true).Where("username = ?", userInfo.Username).First(&user).Error
	if err != nil {
		return nil, err
	}
	for _, group := range user.Groups {
		roles = append(roles, group.Name)
	}

	// 获取用户接口权限
	permissions, err := CasBin.GetUserPermissions(userInfo.Username)
	if err != nil {
		return nil, err
	}

	userInfoWithMenu := &UserInfoWithMenu{
		UserInfo:    *userInfo,
		Menus:       menus,
		Roles:       roles,
		Permissions: permissions,
	}

	return userInfoWithMenu, nil
}

// AddUser 新增
func (u *user) AddUser(data *model.AuthUser) (user *model.AuthUser, err error) {
	if err := global.MySQLClient.Create(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// SyncUsers 用户同步
func (u *user) SyncUsers(users []*model.AuthUser) (err error) {

	// 使用事务（Transaction）来处理批量插入，这样可以确保数据一致性，要么全部成功，要么全部失败
	if err := global.MySQLClient.Transaction(func(tx *gorm.DB) error {
		// 遍历用户列表，逐个插入，需要对数据库中相同用户名的用户进行特殊处理
		for _, user := range users {
			// 新增用户
			err := tx.Create(user).Error
			if err != nil {
				// 如果用户已存在则更新
				if utils.IsDuplicateEntryError(err) {
					// 仅更新来源为LDAP的用户，则进行用户更新
					if err := tx.Select("email", "phone_number", "password_expired_at", "is_active").Where("username = ? AND user_from = ?", user.Username, user.UserFrom).Updates(user).Error; err != nil {
						return err
					}
				} else {
					return err
				}
			} else {
				continue
			}
		}

		// 获取当前数据库中所有来源为 LDAP 的用户
		var localLDAPUsers []model.AuthUser
		if err := tx.Where("user_from = ?", "LDAP").Find(&localLDAPUsers).Error; err != nil {
			return err
		}

		// 提取传入的 LDAP 用户名集合
		usernamesFromInput := make(map[string]bool)
		for _, user := range users {
			usernamesFromInput[user.Username] = true
		}

		// 禁用本地存在但不在传入列表中的 LDAP 用户
		for _, localUser := range localLDAPUsers {
			if !usernamesFromInput[localUser.Username] {
				if err := tx.Model(&localUser).Update("is_active", false).Error; err != nil {
					return err
				}
			}
		}

		return nil
	}); err != nil {
		return err
	}
	return nil
}

// UpdateUser 修改
func (u *user) UpdateUser(user *model.AuthUser, data *UserUpdate) (*model.AuthUser, error) {

	var userinfo = data

	// 判断是否为空，为空则设置为nil
	if *data.WwId == "" {
		userinfo.WwId = nil
	}
	if *data.CtyunId == "" {
		userinfo.CtyunId = nil
	}

	// 当is_active=0，需要使用Select选中对应字段进行更新，否则无法设置为0
	if err := global.MySQLClient.Model(user).Select("*").Updates(data).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUserPasswordExpiredAt 修改用户密码过期时间
func (u *user) UpdateUserPasswordExpiredAt(userId uint, passwordExpiredAt *time.Time) (err error) {
	return global.MySQLClient.Model(&model.AuthUser{}).Where("id = ?", userId).Update("password_expired_at", passwordExpiredAt).Error
}

// DeleteUser 删除
func (u *user) DeleteUser(tx *gorm.DB, id int) (err error) {
	return tx.Where("id = ?", id).Unscoped().Delete(&model.AuthUser{}).Error
}

// UpdateUserPassword 更改密码
func (u *user) UpdateUserPassword(user *model.AuthUser, data *UserPasswordUpdate) (err error) {

	// 对密码进行加密
	cipherText, err := utils.Encrypt(data.Password)
	if err != nil {
		return err
	}

	// 指定新密码有效期
	currentTime := time.Now()
	passwordExpiredAtDays := config.Conf.Settings["passwordExpireDays"].(int)
	passwordExpiredAt := currentTime.AddDate(0, 0, passwordExpiredAtDays)

	// 更新密码
	return global.MySQLClient.Model(&user).Updates(map[string]interface{}{"password": cipherText, "password_expired_at": passwordExpiredAt}).Error
}

// ResetUserMFA 重置MFA
func (u *user) ResetUserMFA(data *model.AuthUser) (err error) {

	// 将MFA重置为nil
	return global.MySQLClient.Model(&model.AuthUser{}).Where("id = ?", data.ID).Update("mfa_code", nil).Error
}

// GetPasswordExpiredUserList 获取密码过期用户列表
func (u *user) GetPasswordExpiredUserList() (userList []*PasswordExpiredUserList, err error) {

	passwordExpiryReminderDays := config.Conf.Settings["passwordExpiryReminderDays"].(int)

	var (
		results        []*PasswordExpiredUserList
		now            = time.Now()
		sevenDaysLater = now.Add(time.Duration(passwordExpiryReminderDays) * 24 * time.Hour)
	)

	if err := global.MySQLClient.Model(&model.AuthUser{}).Select("name, username, email, password_expired_at").
		Where("is_active = ?", true).
		Where("password_expired_at IS NOT NULL AND (password_expired_at < ? OR password_expired_at BETWEEN ? AND ?)", now, now, sevenDaysLater).
		Where("email IS NOT NULL").
		Find(&results).
		Error; err != nil {
		return nil, err
	}

	return results, nil
}
