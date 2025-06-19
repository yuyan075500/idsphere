package service

import (
	"errors"
	"fmt"
	"github.com/pquerna/otp/totp"
	"ops-api/config"
	"ops-api/dao"
	"ops-api/global"
	"ops-api/model"
	"ops-api/utils"
	"ops-api/utils/mail"
	"strconv"
	"time"
)

var Account account

type account struct{}

// AccountCreate 账号创建结构体
type AccountCreate struct {
	Name         string `json:"name" binding:"required"`
	LoginAddress string `json:"login_address"`
	LoginMethod  string `json:"login_method"`
	Username     string `json:"username"`
	Password     string `json:"password" binding:"required"`
	Note         string `json:"note"`
}

// BatchAccountCreate 批量账号创建结构体
type BatchAccountCreate struct {
	Accounts []model.Account `json:"accounts" binding:"required,dive"`
}

// CodeVerification 获取密码结构体
type CodeVerification struct {
	ValidateType uint   `json:"validate_type" binding:"required"` // 验证类型：1：短信验证码，2：MFA验证码，3：邮箱验证码
	Code         string `json:"code" binding:"required"`
}

// GetVerification 获取验证码结构体
type GetVerification struct {
	ValidateType uint `json:"validate_type" binding:"required"` // 验证类型：1：短信验证码，2：MFA验证码，3：邮箱验证码
}

// AddAccount 创建账号
func (a *account) AddAccount(data *AccountCreate, userId uint) (res *model.Account, err error) {

	account := &model.Account{
		Name:         data.Name,
		Username:     data.Username,
		Password:     data.Password,
		LoginAddress: data.LoginAddress,
		Note:         data.Note,
		LoginMethod:  data.LoginMethod,
		OwnerUserID:  userId,
	}

	return dao.Account.AddAccount(account)
}

// AddAccounts 批量创建账号
func (a *account) AddAccounts(accounts []model.Account, userId uint) (res []model.Account, err error) {

	// 设置每个账号的所有者
	for i := range accounts {
		accounts[i].OwnerUserID = userId
	}

	return dao.Account.AddAccounts(accounts)
}

// GetAccountList 获取账号列表（表格）
func (a *account) GetAccountList(name string, userID uint, page, limit int) (data *dao.AccountList, err error) {
	data, err = dao.Account.GetAccountList(name, userID, page, limit)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// DeleteAccount 删除账号
func (a *account) DeleteAccount(id, userId int) error {

	// 判断是否有权限操作
	owner, _, err := dao.Account.GetAccountOwnerAndUsers(id)
	if err != nil {
		return err
	}
	if owner.ID != uint(userId) {
		return errors.New("此账号你无权操作")
	}

	return dao.Account.DeleteAccount(id)
}

// UpdateAccount 更新账号
func (a *account) UpdateAccount(data *dao.AccountUpdate, userId uint) (*model.Account, error) {

	// 判断是否有权限操作
	owner, _, err := dao.Account.GetAccountOwnerAndUsers(int(data.ID))
	if err != nil {
		return nil, err
	}
	if owner.ID != userId {
		return nil, errors.New("此账号你无权操作")
	}

	return dao.Account.UpdateAccount(data)
}

// BatchUpdateAccountOwner 批量修改账号所有者
func (a *account) BatchUpdateAccountOwner(accounts []uint, oldOwnerId, newOwnerId uint) ([]*model.Account, error) {

	// 判断是否有权限操作
	for _, accountId := range accounts {
		owner, _, err := dao.Account.GetAccountOwnerAndUsers(int(accountId))
		if err != nil {
			return nil, err
		}
		if owner.ID != oldOwnerId {
			return nil, errors.New("此账号你无权操作")
		}
	}

	return dao.Account.BatchUpdateAccountOwner(accounts, newOwnerId)
}

// UpdateAccountUser 用户分享
func (a *account) UpdateAccountUser(data *dao.AccountUpdateUser, userId uint) (*model.Account, error) {

	// 判断是否有权限操作
	owner, _, err := dao.Account.GetAccountOwnerAndUsers(int(data.ID))
	if err != nil {
		return nil, err
	}
	if owner.ID != userId {
		return nil, errors.New("此账号你无权操作")
	}

	var (
		users   []model.AuthUser
		account = &model.Account{}
	)

	// 查询要修改的用户组
	if err := global.MySQLClient.First(account, data.ID).Error; err != nil {
		return nil, err
	}

	// 查询出要更新的所有用户
	if len(data.Users) == 0 {
		users = []model.AuthUser{}
	} else {
		if err := global.MySQLClient.Find(&users, data.Users).Error; err != nil {
			return nil, err
		}
	}

	return dao.Account.UpdateAccountUser(account, users)
}

// UpdatePassword 更改密码
func (a *account) UpdatePassword(data *dao.AccountUpdatePassword, userId uint) (err error) {

	// 判断是否有权限操作
	owner, _, err := dao.Account.GetAccountOwnerAndUsers(int(data.ID))
	if err != nil {
		return err
	}
	if owner.ID != userId {
		return errors.New("此账号你无权操作")
	}

	// 密码校验
	if data.Password != data.RePassword {
		return errors.New("两次输入的密码不匹配")
	}
	// 检查密码复杂度
	//if err := check.PasswordCheck(data.Password); err != nil {
	//	return err
	//}

	// 查询要修改的账号
	account := &model.Account{}
	if err := global.MySQLClient.First(account, data.ID).Error; err != nil {
		return err
	}

	return dao.Account.UpdatePassword(account, data)
}

// GetAccountPassword 获取账号密码
func (a *account) GetAccountPassword(id uint, username string, userId uint) (password *string, err error) {

	// 判断是否有权限操作
	owner, users, err := dao.Account.GetAccountOwnerAndUsers(int(id))
	if err != nil {
		return nil, err
	}
	if owner.ID != userId {
		hasAccess := false
		for _, user := range users {
			if user.ID == userId {
				hasAccess = true
				break
			}
		}
		if !hasAccess {
			return nil, errors.New("此账号你无权操作")
		}
	}

	// 判断是否需要认证，Redis缓存中指定的Key是否存在，存在则不需要认证，否则需要认证
	var keyName = fmt.Sprintf("%s_get_account_password_enabled", username)
	val, err := global.RedisClient.Exists(keyName).Result()
	if err != nil {
		return nil, err
	}
	// 0 表示不存在
	if val == 0 {
		return nil, err
	}

	// 获取密码
	p, err := dao.Account.GetAccountPassword(id)
	if err != nil {
		return nil, err
	}

	// 密码解密
	pwd := p
	str, err := utils.Decrypt(pwd)
	if err != nil {
		return nil, err
	}

	return &str, nil
}

// GetSMSCode 发送获取账号密码验证码
func (a *account) GetSMSCode(data *GetVerification, userID uint) (err error) {

	// 查找用户
	conditions := map[string]interface{}{
		"id": userID,
	}
	user, err := dao.User.GetUser(conditions)
	if err != nil {
		return err
	}

	var keyName = fmt.Sprintf("%s_get_account_password_verification_code", user.Username)

	// 判断Redis缓存中指定的Key是否存在
	val, err := global.RedisClient.Exists(keyName).Result()
	if err != nil {
		return err
	}
	// 1 表示已存在
	if val == 1 {
		// 判断Key的有效期，如果Key的有效期大于4分钟，表示在1分钟内发送过验证码，提示用户请勿频繁发送校验码
		ttl, err := global.RedisClient.TTL(keyName).Result()
		if err != nil {
			return err
		}
		if ttl.Seconds() > 240 {
			return errors.New("请勿频繁发送校验码")
		}
	}

	// 发送短信验证码
	var code string
	if data.ValidateType == 1 {
		if user.PhoneNumber == "" {
			return errors.New("账号未绑定手机号")
		}
		number, err := SMS.SMSSend(user.PhoneNumber, "密码查询")
		if err != nil {
			return err
		}
		code = number
	}
	if data.ValidateType == 3 {

		// 判断是否开启此功能
		if config.Conf.Settings["passwordMailResetOff"].(bool) == false {
			return errors.New("功能未启用，请联系管理员")
		}

		if user.Email == "" {
			return errors.New("账号未绑定邮箱，请联系管理员")
		}

		// 生成HTML内容
		code = strconv.Itoa(utils.GenerateRandomNumber())
		htmlBody := GetPasswordHTML(code)

		// 发送邮件
		if err := mail.Email.SendMsg([]string{user.Email}, nil, nil, "密码获取", htmlBody, "html"); err != nil {
			return err
		}
	}

	// 将验证码写入Redis缓存，如果已存在则会更新Key的值并刷新TTL
	if err := global.RedisClient.Set(keyName, code, 5*time.Minute).Err(); err != nil {
		return err
	}

	return nil
}

func GetPasswordHTML(number string) string {

	issuer := config.Conf.Settings["issuer"].(string)

	return fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>密码获取</title>
		</head>
		<body>
			<p>您的验证码为：%s，请妥善保管，切勿泄露。</p>
			<br>
			<p>此致，<br>%s</p>
			<p style="color: red">此邮件为系统自动发送，请勿回复此邮件。</p>
		</body>
		</html>
	`, number, issuer)
}

// CodeVerification 校验验证码
func (a *account) CodeVerification(userId uint, data *CodeVerification) (err error) {

	// 查找用户
	conditions := map[string]interface{}{
		"id": userId,
	}
	user, err := dao.User.GetUser(conditions)
	if err != nil {
		return err
	}

	// 验证码校验
	if data.ValidateType == 1 || data.ValidateType == 3 {

		// 从缓存中获取验证码
		result, err := global.RedisClient.Get(fmt.Sprintf("%s_get_account_password_verification_code", user.Username)).Result()
		if err != nil {
			return err
		}
		// 判断是否正确
		if result != data.Code {
			return errors.New("校验码错误")
		}
	} else if data.ValidateType == 2 {
		// MFA验证码校验
		// 获取Secret
		if user.MFACode == nil {
			return errors.New("您还未绑定MFA")
		}

		// 校验MFA
		valid := totp.Validate(data.Code, *user.MFACode)
		if !valid {
			return errors.New("校验码错误")
		}
	} else {
		return errors.New("验证码类型错误")
	}

	// 写入允许用户获取密码的Redis缓存
	return global.RedisClient.Set(fmt.Sprintf("%s_get_account_password_enabled", user.Username), true, 5*time.Minute).Err()
}
