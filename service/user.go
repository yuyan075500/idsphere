package service

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"
	"ops-api/config"
	"ops-api/dao"
	"ops-api/global"
	"ops-api/middleware"
	"ops-api/model"
	"ops-api/utils"
	"ops-api/utils/check"
	"ops-api/utils/mail"
	messages "ops-api/utils/sms"
	"time"
)

var User user

type user struct{}

// AuthorizeParam 通用的单点登录授权参数接口
type AuthorizeParam interface {
	GetResponseType() string
	GetClientId() string
	GetRedirectURI() string
	GetService() string
	GetSAMLRequest() string
	GetRelayState() string
	GetSigAlg() string
	GetSignature() string
	GetScope() string
	GetState() string
	GetNonce() string
}

// UserLogin 用户登录结构体（支持CAS3.0、OAuth2.0、OIDC和SAML2）
type UserLogin struct {
	Username         string `json:"username" binding:"required"`
	Password         string `json:"password" binding:"required"`
	ResponseType     string `json:"response_type"`      // OAuth2.0客户端：授权类型，固定值：code
	ClientId         string `json:"client_id"`          // OAuth2.0客户端：客户端ID
	RedirectURI      string `json:"redirect_uri"`       // OAuth2.0客户端：重定向URL
	State            string `json:"state"`              // OAuth2.0客户端：客户端状态码
	Scope            string `json:"scope"`              // OAuth2.0客户端：申请权限范围
	Nonce            string `json:"nonce"`              // OIDC客户端：随机码
	Service          string `json:"service"`            // CAS3.0客户端：回调地址
	SAMLRequest      string `json:"SAMLRequest"`        // SAML2客户端：SAMLRequest
	RelayState       string `json:"RelayState"`         // SAML2客户端：客户端状态码
	SigAlg           string `json:"SigAlg"`             // SAML2客户端：签名算法
	Signature        string `json:"Signature"`          // SAML2客户端：签名
	NginxRedirectURI string `json:"nginx_redirect_uri"` // Nginx代理客户端：回调地址
}

// DingTalkLogin 钉钉扫码登录结构体（支持CAS3.0、OAuth2.0、OIDC和SAML2）
type DingTalkLogin struct {
	AuthCode     string `json:"authCode" binding:"required"`
	ResponseType string `json:"response_type"` // OAuth2.0客户端：授权类型，固定值：code
	ClientId     string `json:"client_id"`     // OAuth2.0客户端：客户端ID
	RedirectURI  string `json:"redirect_uri"`  // OAuth2.0客户端：重定向URL
	State        string `json:"state"`         // OAuth2.0客户端：客户端状态码
	Scope        string `json:"scope"`         // OAuth2.0客户端：申请权限范围
	Nonce        string `json:"nonce"`         // OIDC客户端：随机码
	Service      string `json:"service"`       // CAS3.0客户端：回调地址
	SAMLRequest  string `json:"SAMLRequest"`   // SAML2客户端：SAMLRequest
	RelayState   string `json:"RelayState"`    // SAML2客户端：客户端状态码
	SigAlg       string `json:"SigAlg"`        // SAML2客户端：签名算法
	Signature    string `json:"Signature"`     // SAML2客户端：签名
}

// WeChatLogin 企业微信扫码登录结构体（支持CAS3.0、OAuth2.0、OIDC和SAML2）
type WeChatLogin struct {
	Code         string `json:"code" binding:"required"`
	Appid        string `json:"appid" binding:"required"`
	ResponseType string `json:"response_type"` // OAuth2.0客户端：授权类型，固定值：code
	ClientId     string `json:"client_id"`     // OAuth2.0客户端：客户端ID
	RedirectURI  string `json:"redirect_uri"`  // OAuth2.0客户端：重定向URL
	State        string `json:"state"`         // OAuth2.0客户端：客户端状态码
	Scope        string `json:"scope"`         // OAuth2.0客户端：申请权限范围
	Nonce        string `json:"nonce"`         // OIDC客户端：随机码
	Service      string `json:"service"`       // CAS3.0客户端：回调地址
	SAMLRequest  string `json:"SAMLRequest"`   // SAML2客户端：SAMLRequest
	RelayState   string `json:"RelayState"`    // SAML2客户端：客户端状态码
	SigAlg       string `json:"SigAlg"`        // SAML2客户端：签名算法
	Signature    string `json:"Signature"`     // SAML2客户端：签名
}

// FeishuLogin 飞书扫码登录结构体（支持CAS3.0、OAuth2.0、OIDC和SAML2）
type FeishuLogin struct {
	Code         string `json:"code" binding:"required"`
	Byte         string `json:"byte" binding:"required"` // 自定义参数
	ResponseType string `json:"response_type"`           // OAuth2.0客户端：授权类型，固定值：code
	ClientId     string `json:"client_id"`               // OAuth2.0客户端：客户端ID
	RedirectURI  string `json:"redirect_uri"`            // OAuth2.0客户端：重定向URL
	State        string `json:"state"`                   // OAuth2.0客户端：客户端状态码
	Scope        string `json:"scope"`                   // OAuth2.0客户端：申请权限范围
	Nonce        string `json:"nonce"`                   // OIDC客户端：随机码
	Service      string `json:"service"`                 // CAS3.0客户端：回调地址
	SAMLRequest  string `json:"SAMLRequest"`             // SAML2客户端：SAMLRequest
	RelayState   string `json:"RelayState"`              // SAML2客户端：客户端状态码
	SigAlg       string `json:"SigAlg"`                  // SAML2客户端：签名算法
	Signature    string `json:"Signature"`               // SAML2客户端：签名
}

func (f FeishuLogin) GetResponseType() string { return f.ResponseType }
func (f FeishuLogin) GetClientId() string     { return f.ClientId }
func (f FeishuLogin) GetRedirectURI() string  { return f.RedirectURI }
func (f FeishuLogin) GetService() string      { return f.Service }
func (f FeishuLogin) GetSAMLRequest() string  { return f.SAMLRequest }
func (f FeishuLogin) GetRelayState() string   { return f.RelayState }
func (f FeishuLogin) GetSigAlg() string       { return f.SigAlg }
func (f FeishuLogin) GetSignature() string    { return f.Signature }
func (f FeishuLogin) GetScope() string        { return f.Scope }
func (f FeishuLogin) GetState() string        { return f.State }
func (f FeishuLogin) GetNonce() string        { return f.Nonce }

func (d DingTalkLogin) GetResponseType() string { return d.ResponseType }
func (d DingTalkLogin) GetClientId() string     { return d.ClientId }
func (d DingTalkLogin) GetRedirectURI() string  { return d.RedirectURI }
func (d DingTalkLogin) GetService() string      { return d.Service }
func (d DingTalkLogin) GetSAMLRequest() string  { return d.SAMLRequest }
func (d DingTalkLogin) GetRelayState() string   { return d.RelayState }
func (d DingTalkLogin) GetSigAlg() string       { return d.SigAlg }
func (d DingTalkLogin) GetSignature() string    { return d.Signature }
func (d DingTalkLogin) GetScope() string        { return d.Scope }
func (d DingTalkLogin) GetState() string        { return d.State }
func (d DingTalkLogin) GetNonce() string        { return d.Nonce }

func (w WeChatLogin) GetResponseType() string { return w.ResponseType }
func (w WeChatLogin) GetClientId() string     { return w.ClientId }
func (w WeChatLogin) GetRedirectURI() string  { return w.RedirectURI }
func (w WeChatLogin) GetService() string      { return w.Service }
func (w WeChatLogin) GetSAMLRequest() string  { return w.SAMLRequest }
func (w WeChatLogin) GetRelayState() string   { return w.RelayState }
func (w WeChatLogin) GetSigAlg() string       { return w.SigAlg }
func (w WeChatLogin) GetSignature() string    { return w.Signature }
func (w WeChatLogin) GetScope() string        { return w.Scope }
func (w WeChatLogin) GetState() string        { return w.State }
func (w WeChatLogin) GetNonce() string        { return w.Nonce }

func (u UserLogin) GetResponseType() string { return u.ResponseType }
func (u UserLogin) GetClientId() string     { return u.ClientId }
func (u UserLogin) GetRedirectURI() string  { return u.RedirectURI }
func (u UserLogin) GetService() string      { return u.Service }
func (u UserLogin) GetSAMLRequest() string  { return u.SAMLRequest }
func (u UserLogin) GetRelayState() string   { return u.RelayState }
func (u UserLogin) GetSigAlg() string       { return u.SigAlg }
func (u UserLogin) GetSignature() string    { return u.Signature }
func (u UserLogin) GetScope() string        { return u.Scope }
func (u UserLogin) GetState() string        { return u.State }
func (u UserLogin) GetNonce() string        { return u.Nonce }

// RestPassword 重置密码时用户信息绑定的结构体
type RestPassword struct {
	Username     string `json:"username" binding:"required"`
	PhoneNumber  string `json:"phone_number"`
	ValidateType uint   `json:"validate_type" binding:"required"` // 验证类型：1：短信验证码，2：MFA验证码
	Code         string `json:"code" binding:"required"`
	Password     string `json:"password" binding:"required"`
	RePassword   string `json:"re_password" binding:"required"`
}

// GetUserListAll 获取用户列表（下拉框、穿梭框）
func (u *user) GetUserListAll() (data *dao.UserListAll, err error) {
	data, err = dao.User.GetUserListAll()
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetUserList 获取用户列表（表格）
func (u *user) GetUserList(name string, page, limit int) (data *dao.UserList, err error) {
	data, err = dao.User.GetUserList(name, page, limit)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetUser 获取用户信息
func (u *user) GetUser(userid uint) (user *dao.UserInfoWithMenu, err error) {
	user, err = dao.User.GetUserInfo(userid)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// AddUser 创建用户
func (u *user) AddUser(data *dao.UserCreate) (authUser *model.AuthUser, err error) {

	// 字段校验
	validate := validator.New()
	// 注册自定义检验方法
	if err := validate.RegisterValidation("phone", check.PhoneNumberCheck); err != nil {
		return nil, err
	}
	if err := validate.Struct(data); err != nil {
		return nil, err.(validator.ValidationErrors)
	}

	// 检查密码是否复合要求
	if err := check.PasswordCheck(data.Password); err != nil {
		return nil, err
	}

	fmt.Printf("Name: %v, Type: %T\n", data.PasswordExpiredAt, data.PasswordExpiredAt)

	user := &model.AuthUser{
		Name:              data.Name,
		Username:          data.Username,
		Password:          data.Password,
		PhoneNumber:       data.PhoneNumber,
		IsActive:          true,
		Email:             data.Email,
		UserFrom:          data.UserFrom,
		PasswordExpiredAt: data.PasswordExpiredAt,
	}

	return dao.User.AddUser(user)
}

// DeleteUser 删除
func (u *user) DeleteUser(id int) (err error) {

	// 定义用户匹配条件
	conditions := map[string]interface{}{
		"id": id,
	}

	// 在本地数据库中查找匹配的用户
	user, err := dao.User.GetUser(conditions)
	if err != nil {
		return err
	}

	if user.ID == 1 || user.Username == "admin" {
		return errors.New("超级管理员不允许删除")
	}

	// 开始事务
	tx := global.MySQLClient.Begin()

	// 删除用户关联分组
	if err := tx.Model(&user).Association("Groups").Clear(); err != nil {
		tx.Rollback()
		return err
	}

	// 删除用户
	if err = dao.User.DeleteUser(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// UpdateUser 更新
func (u *user) UpdateUser(data *dao.UserUpdate) (*model.AuthUser, error) {

	// 字段校验
	validate := validator.New()
	// 注册自定义检验方法
	if err := validate.RegisterValidation("phone", check.PhoneNumberCheck); err != nil {
		return nil, err
	}
	if err := validate.Struct(data); err != nil {
		return nil, err.(validator.ValidationErrors)
	}

	// 查询要修改的用户
	user := &model.AuthUser{}
	if err := global.MySQLClient.First(user, data.ID).Error; err != nil {
		return nil, err
	}

	return dao.User.UpdateUser(user, data)
}

// UpdateUserPassword 更改密码
func (u *user) UpdateUserPassword(data *dao.UserPasswordUpdate) (err error) {
	// 密码校验
	if data.Password != data.RePassword {
		return errors.New("两次输入的密码不匹配")
	}
	if err := check.PasswordCheck(data.Password); err != nil {
		return err
	}

	// 查询要修改的用户
	user := &model.AuthUser{}
	if err := global.MySQLClient.First(user, data.ID).Error; err != nil {
		return err
	}

	// 如果是LDAP账号则访问LDAP进行用户密码重置
	if user.UserFrom == "LDAP" {
		if err := AD.LDAPUserResetPassword(user.Username, data.Password); err != nil {
			return err
		}
	}

	return dao.User.UpdateUserPassword(user, data)
}

// ResetUserMFA 重置MFA
func (u *user) ResetUserMFA(id int) error {

	// 查询要重置的用户
	user := &model.AuthUser{}
	if err := global.MySQLClient.First(user, id).Error; err != nil {
		return err
	}

	return dao.User.ResetUserMFA(user)
}

// GetVerificationCode 获取重置密码短信验证码
func (u *user) GetVerificationCode(data *messages.SendData, expirationTime int) (err error) {

	var (
		keyName = fmt.Sprintf("%s_rest_password_verification_code", data.Username)
	)

	// 判断Redis缓存中指定的Key是否存在
	val, err := global.RedisClient.Exists(keyName).Result()
	if err != nil {
		return err
	}
	// 已存在
	if val >= 1 {
		// 判断Key的有效期
		ttl, err := global.RedisClient.TTL(keyName).Result()
		if err != nil {
			return err
		}
		// 如果Key的有效期大于4分钟，则提示用户请勿频繁发送校验码
		if ttl.Minutes() > 4 {
			return errors.New("请勿频繁发送校验码")
		}
	}

	// 发送短信验证码
	data.Note = "密码重置"
	code, err := SMS.SMSSend(data)
	if err != nil {
		return err
	}

	// 将验证码写入Redis缓存，如果已存在则会更新Key的值并刷新TTL
	return global.RedisClient.Set(keyName, code, time.Duration(expirationTime)*time.Minute).Err()
}

// UpdateSelfPassword 用户重置密码
func (u *user) UpdateSelfPassword(data *RestPassword) (err error) {

	var (
		user    model.AuthUser
		keyName = fmt.Sprintf("%s_rest_password_verification_code", data.Username)
	)

	// 获取用户信息
	if err := global.MySQLClient.First(&user, "username", data.Username).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 短信验证码校验
	if data.ValidateType == 1 {

		if data.PhoneNumber != user.PhoneNumber {
			return errors.New("手机号与用户不匹配")
		}

		// 从缓存中获取验证码
		result, err := global.RedisClient.Get(keyName).Result()
		if err != nil {
			return err
		}
		// 判断是否正确
		if result != data.Code {
			return errors.New("校验码错误")
		}
	}

	// MFA验证码校验
	if data.ValidateType == 2 {

		// 获取Secret
		if user.MFACode == nil {
			return errors.New("您还未绑定MFA")
		}

		// 校验MFA
		valid := totp.Validate(data.Code, *user.MFACode)
		if !valid {
			return errors.New("校验码错误")
		}
	}

	// 执行密码重置
	userInfo := &dao.UserPasswordUpdate{
		ID:         user.ID,
		Password:   data.Password,
		RePassword: data.RePassword,
	}
	return u.UpdateUserPassword(userInfo)
}

// DingTalkLogin 钉钉扫码认证
func (u *user) DingTalkLogin(params *DingTalkLogin) (token, redirectUri, username, application string, err error) {

	// 初始化钉钉客户端
	dingClient, err := NewDingTalkClient()
	if err != nil {
		return "", "", "", "", err
	}

	// 获取用户Token
	userAccessToken, err := dingClient.GetUserAccessToken(params.AuthCode)
	if err != nil {
		return "", "", "", "", err
	}

	// 获取用户信息
	userInfo, err := dingClient.GetUserInfo(userAccessToken)
	if err != nil {
		return "", "", "", "", err
	}

	// 定义用户匹配条件
	conditions := map[string]interface{}{
		"name":         userInfo.Body.Nick,
		"phone_number": userInfo.Body.Mobile,
	}

	// 在本地数据库中查找匹配的用户
	user, err := dao.User.GetUser(conditions)
	if err != nil {
		return "", "", "", "", err
	}

	// 判断用户是否禁用
	if !user.IsActive {
		return "", "", "", "", errors.New("拒绝登录，请联系管理员")
	}

	// 判断密码是否过期
	if user.PasswordExpiredAt != nil {
		now := time.Now()
		if user.PasswordExpiredAt.Before(now) {
			return "", "", "", "", errors.New("密码过期")
		}
	}

	// 生成用户Token
	token, err = middleware.GenerateJWT(user.ID, user.Name, user.Username)
	if err != nil {
		return "", "", "", "", err
	}

	// 处理单点登录请求
	if params.SAMLRequest != "" || params.Service != "" || params.ClientId != "" {
		callbackData, siteName, err := SSO.Login(params, *user)
		if err != nil {
			return "", "", "", siteName, err
		}
		// 这里的callbackData，如果是SAML2认证则为html，如果是其它认证则为回调地址
		return token, callbackData, user.Username, siteName, nil
	}

	return token, "", user.Username, "", nil
}

// FeishuLogin 飞书扫码认证
func (u *user) FeishuLogin(params *FeishuLogin) (token, redirectUri, username, application string, err error) {

	client, err := NewFeishuClient()
	if err != nil {
		return "", "", "", "", err
	}

	// 获取用户访问Token
	resp, err := client.GetUserAccessToken(params.Code)
	if err != nil {
		return "", "", "", "", err
	}

	// 获取用户信息
	userinfo, err := client.GetUserInfo(*resp.Data.AccessToken)

	// 定义用户匹配条件
	mobile := *userinfo.Data.Mobile
	conditions := map[string]interface{}{
		"name":         *userinfo.Data.Name,
		"phone_number": mobile[3:],
		"email":        *userinfo.Data.Email,
	}

	// 在本地数据库中查找匹配的用户
	user, err := dao.User.GetUser(conditions)
	if err != nil {
		return "", "", "", "", err
	}

	// 判断用户是否禁用
	if !user.IsActive {
		return "", "", "", "", errors.New("拒绝登录，请联系管理员")
	}

	// 判断密码是否过期
	if user.PasswordExpiredAt != nil {
		now := time.Now()
		if user.PasswordExpiredAt.Before(now) {
			return "", "", "", "", errors.New("密码过期")
		}
	}

	// 生成用户Token
	token, err = middleware.GenerateJWT(user.ID, user.Name, user.Username)
	if err != nil {
		return "", "", "", "", err
	}

	// 处理单点登录请求
	if params.SAMLRequest != "" || params.Service != "" || params.ClientId != "" {
		callbackData, siteName, err := SSO.Login(params, *user)
		if err != nil {
			return "", "", "", siteName, err
		}
		// 这里的callbackData，如果是SAML2认证则为html，如果是其它认证则为回调地址
		return token, callbackData, user.Username, siteName, nil
	}

	return token, "", user.Username, "", nil
}

// WeChatLogin 企业微信扫码认证
func (u *user) WeChatLogin(params *WeChatLogin) (token, redirectUri, username, application string, err error) {

	wechatClient, err := NewWeChatClient()
	if err != nil {
		return "", "", "", "", err
	}

	// 获取用户ID
	userId, err := wechatClient.GetUserId(params.Code)
	if err != nil {
		return "", "", "", "", err
	}

	// 定义用户匹配条件
	conditions := map[string]interface{}{
		"ww_id": userId,
	}

	// 在本地数据库中查找匹配的用户
	user, err := dao.User.GetUser(conditions)
	if err != nil {
		return "", "", "", "", err
	}

	// 判断用户是否禁用
	if !user.IsActive {
		return "", "", "", "", errors.New("拒绝登录，请联系管理员")
	}

	// 判断密码是否过期
	if user.PasswordExpiredAt != nil {
		now := time.Now()
		if user.PasswordExpiredAt.Before(now) {
			return "", "", "", "", errors.New("密码过期")
		}
	}

	// 生成用户Token
	token, err = middleware.GenerateJWT(user.ID, user.Name, user.Username)
	if err != nil {
		return "", "", "", "", err
	}

	// 处理单点登录请求
	if params.SAMLRequest != "" || params.Service != "" || params.ClientId != "" {
		callbackData, siteName, err := SSO.Login(params, *user)
		if err != nil {
			return "", "", "", siteName, err
		}
		// 这里的callbackData，如果是SAML2认证则为html，如果是其它认证则为回调地址
		return token, callbackData, user.Username, siteName, nil
	}

	return token, "", user.Username, "", nil
}

// Login 用户登录（支持CAS3.0、OAuth2.0、OIDC和SAML2）
func (u *user) Login(params *UserLogin) (token, redirectUri, application string, mfaPage *string, err error) {

	var user model.AuthUser

	// 用户认证
	if err := u.AuthenticateUser(params, &user); err != nil {
		return "", "", "", nil, err
	}

	// 判断用户是否禁用
	if !user.IsActive {
		return "", "", "", nil, errors.New("拒绝登录，请联系管理员")
	}

	// 判断密码是否过期
	if user.PasswordExpiredAt != nil {
		now := time.Now()
		if user.PasswordExpiredAt.Before(now) {
			return "", "", "", nil, errors.New("密码已过期")
		}
	}

	// 判断系统是否启用MFA认证
	mfaEnable := config.Conf.Settings["mfa"].(bool)
	if mfaEnable {
		token, nextPage, err := handleMFA(user)
		if err != nil {
			return "", "", "", nil, err
		}
		if nextPage != nil {
			return token, "", "", nextPage, nil
		}
	}

	// 生成用户Token
	token, err = middleware.GenerateJWT(user.ID, user.Name, user.Username)
	if err != nil {
		return "", "", "", nil, err
	}

	// 处理单点登录请求
	if params.SAMLRequest != "" || params.Service != "" || params.ClientId != "" {
		callbackData, siteName, err := SSO.Login(params, user)
		if err != nil {
			return "", "", siteName, nil, err
		}
		// 这里的callbackData，如果是SAML2认证则为html，如果是其它认证则为回调地址
		return token, callbackData, siteName, nil, nil
	}

	return token, "", "", nil, nil
}

// UserSync LDAP用户同步
func (u *user) UserSync() error {
	return AD.LDAPUserSync()
}

// UpdateUserLoginTime 更新用户最后登录时间
func (u *user) UpdateUserLoginTime(tx *gorm.DB, userId uint) (err error) {

	return tx.Model(&model.AuthUser{}).Where("id = ?", userId).Update("last_login_at", time.Now()).Error
}

// AuthenticateUser 用户认证
func (u *user) AuthenticateUser(params *UserLogin, user *model.AuthUser) error {

	// 查找用户
	userQuery := global.MySQLClient.Where("username = ?", params.Username)

	// 没有找到对应的用户
	if err := userQuery.First(&user).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 用户密码检查
	if user.UserFrom == "LDAP" {
		// LDAP用户认证
		if _, err := AD.LDAPUserAuthentication(params.Username, params.Password); err != nil {
			return err
		}
	} else if !user.CheckPassword(params.Password) {
		return errors.New("用户密码错误")
	}

	return nil
}

// RecordLoginInfo 记录用户登录信息
func (u *user) RecordLoginInfo(loginMethod, username, userAgent, clientIP, application string, failedReason error) error {

	// 开启事务
	tx := global.MySQLClient.Begin()

	// 更新用户登录时间
	conditions := map[string]interface{}{
		"username": username,
	}
	user, err := dao.User.GetUser(conditions)
	if err == nil {

		if err := u.UpdateUserLoginTime(tx, user.ID); err != nil {
			tx.Rollback()
			return err
		}
	}

	// failedReason为空，则表示用户登录成功
	if failedReason == nil {
		// 记录登录成功信息
		if err := Audit.AddLoginSuccessRecord(tx, username, userAgent, clientIP, loginMethod, application); err != nil {
			tx.Rollback()
			return err
		}
	} else {
		// 记录登录失败信息
		if err := Audit.AddLoginFailedRecord(tx, username, userAgent, clientIP, loginMethod, application, failedReason); err != nil {
			tx.Rollback()
			return err
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// handleMFA 双因素认证返回
func handleMFA(user model.AuthUser) (temToken string, nextPage *string, err error) {

	// 生成一个32位长度的随机字符串作为临时token
	token := utils.GenerateRandomString(32)

	// 将token写入Redis缓存，并设置有效期为2分钟（这里的时间和前端配置的定时器保持一至）
	if err := global.RedisClient.Set(token, user.Username, 2*time.Minute).Err(); err != nil {
		return "", nil, err
	}

	// 判断用户是否已经绑定MFA，为空则未绑定
	// MFA_AUTH和MFA_ENABLE是在前端定义页面名称，认证通过后会根据redirect的值跳转到对应的页面
	redirect := "MFA_AUTH"
	if user.MFACode == nil {
		redirect = "MFA_ENABLE"
		return token, &redirect, nil
	}

	return token, &redirect, nil
}

// PasswordExpiredNotice 用户密码过期通知
func (u *user) PasswordExpiredNotice() (map[string]interface{}, error) {

	users, err := dao.User.GetPasswordExpiredUserList()
	if err != nil {
		return nil, err
	}

	// 初始化结果集合
	result := map[string]interface{}{
		"success": []string{},
		"failed":  []map[string]string{},
	}

	for _, user := range users {

		// 格式化密码过期时间
		expiredAt := user.PasswordExpiredAt.Format("2006-01-02 15:04:05")

		// 生成HTML内容
		htmlBody := PasswordExpiredNoticeHTML(user.Name, user.Username, expiredAt)

		// 发送邮件函数
		if err := mail.Email.SendMsg([]string{user.Email}, nil, nil, "密码过期提醒", htmlBody, "html"); err != nil {
			// 添加到失败列表
			result["failed"] = append(result["failed"].([]map[string]string), map[string]string{
				"email": user.Email,
				"error": err.Error(),
			})
		} else {
			// 添加到成功列表
			result["success"] = append(result["success"].([]string), user.Email)
		}
	}

	return result, nil
}

// PasswordExpiredNoticeHTML 密码过期通知正文
func PasswordExpiredNoticeHTML(name, username, expiredAt string) string {

	var (
		externalUrl = config.Conf.Settings["externalUrl"].(string)
		issuer      = config.Conf.Settings["issuer"].(string)
	)

	url := externalUrl
	if url != "" && url[len(url)-1] != '/' {
		url += "/"
	}

	resetPasswordURL := url + "reset_password"

	return fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>密码到期提醒</title>
		</head>
		<body>
			<p>亲爱的同事：</p>
			<p>您好，目前检测到您的账户（<strong>%s</strong>）密码在 <strong>%s</strong> 过期。</p>
			<p>为了保护您的账号安全，请及时登录【<a href="%s" target="_blank">密码修改自助平台</a>】修改密码，逾期未修改则会导致您的账户无法登录到相关的平台。</p>
			<br>
			<p>此致，<br>%s</p>
			<p style="color: red">此邮件为系统自动发送，请不要回复此邮件。</p>
		</body>
		</html>
	`, username, expiredAt, resetPasswordURL, issuer)
}
