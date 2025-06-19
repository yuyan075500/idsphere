package service

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	"mime/multipart"
	"ops-api/config"
	"ops-api/dao"
	"ops-api/db"
	"ops-api/global"
	"ops-api/model"
	"ops-api/utils"
	"ops-api/utils/mail"
	"time"
)

var Settings settings

type settings struct{}

// SettingsUpdate 修改配置结构体
type SettingsUpdate struct {
	ExternalUrl                string `json:"externalUrl"`
	Mfa                        string `json:"mfa"`
	Issuer                     string `json:"issuer"`
	Secret                     string `json:"secret"`
	LdapAddress                string `json:"ldapAddress"`
	LdapBindDn                 string `json:"ldapBindDn"`
	LdapBindPassword           string `json:"ldapBindPassword"`
	LdapSearchDn               string `json:"ldapSearchDn"`
	LdapFilterAttribute        string `json:"ldapFilterAttribute"`
	LdapUserPasswordExpireDays string `json:"ldapUserPasswordExpireDays"`
	PasswordExpireDays         string `json:"passwordExpireDays"`
	PasswordLength             string `json:"passwordLength"`
	PasswordComplexity         string `json:"passwordComplexity"`
	PasswordExpiryReminderDays string `json:"passwordExpiryReminderDays"`
	MailAddress                string `json:"mailAddress"`
	MailPort                   string `json:"mailPort"`
	MailForm                   string `json:"mailForm"`
	MailPassword               string `json:"mailPassword"`
	SmsAppSecret               string `json:"smsAppSecret"`
	SmsAppKey                  string `json:"smsAppKey"`
	SmsProvider                string `json:"smsProvider"`
	SmsSignature               string `json:"smsSignature"`
	SmsEndpoint                string `json:"smsEndpoint"`
	SmsSender                  string `json:"smsSender"`
	SmsCallbackUrl             string `json:"smsCallbackUrl"`
	SmsTemplateId              string `json:"smsTemplateId"`
	DingdingAppKey             string `json:"dingdingAppKey"`
	DingdingAppSecret          string `json:"dingdingAppSecret"`
	FeishuAppId                string `json:"feishuAppId"`
	FeishuAppSecret            string `json:"feishuAppSecret"`
	WechatCorpId               string `json:"wechatCorpId"`
	WechatAgentId              string `json:"wechatAgentId"`
	WechatSecret               string `json:"wechatSecret"`
	TokenExpiresTime           string `json:"tokenExpiresTime"`
	Swagger                    string `json:"swagger"`
	PasswordMailResetOff       string `json:"passwordMailResetOff"`
}

type MailTest struct {
	Receiver string `json:"receiver" binding:"required"`
}

type SmsTest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type LoginTest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CertTest struct {
	Certificate string `json:"certificate" binding:"required"`
	PublicKey   string `json:"publicKey" binding:"required"`
	PrivateKey  string `json:"privateKey" binding:"required"`
}

// GetAllSettingsWithParsedValues 获取所有配置
func (s *settings) GetAllSettingsWithParsedValues() (map[string]interface{}, error) {
	configs, err := dao.Settings.GetAllSettings()
	if err != nil {
		return nil, err
	}

	// 创建结果字典
	result := make(map[string]interface{})

	for _, conf := range configs {
		if err := conf.ParseValue(); err != nil {
			return nil, err
		}
		result[conf.Key] = conf.ParsedValue
	}
	return result, nil
}

// UploadLogo 上传 Logo
func (s *settings) UploadLogo(path string, logo *multipart.FileHeader) (url string, err error) {

	// 打开上传的图片
	src, err := logo.Open()
	if err != nil {
		return "", err
	}

	// 检查对象是否存在，err不为空是则表示对象已存在
	_, err = utils.StatObject(path)
	if err == nil {
		return "", err
	}

	// 上传
	err = utils.FileUpload(path, logo.Header.Get("Content-Type"), src, logo.Size)
	if err != nil {
		return "", err
	}

	// 获取预览URL
	logoPreview, err := utils.GetPresignedURL(path, 6*time.Hour)

	return logoPreview.String(), nil
}

// GetLogo 获取 Logo
func (s *settings) GetLogo() (string, error) {
	logo, err := dao.Settings.GetSettingByKey("logo")
	if err != nil {
		return "", err
	}

	if logo.Value == nil {
		return "", nil
	}

	// 获取预览URL
	logoPreview, err := utils.GetPresignedURL(*logo.Value, 6*time.Hour)

	return logoPreview.String(), nil
}

// GetSettingByKeyWithParsedValue 获取单个配置
func (s *settings) GetSettingByKeyWithParsedValue(key string) (*model.Settings, error) {
	setting, err := dao.Settings.GetSettingByKey(key)
	if err != nil {
		return nil, err
	}

	// 值类型转换
	if err := setting.ParseValue(); err != nil {
		return nil, err
	}
	return setting, nil
}

// UpdateSettingValue 更新单个配置
func (s *settings) UpdateSettingValue(key, value string) error {
	return dao.Settings.UpdateSetting(key, value)
}

// UpdateSettingValues 更新多个配置
func (s *settings) UpdateSettingValues(data *SettingsUpdate) (map[string]interface{}, error) {

	settingsToUpdate := map[string]interface{}{}

	// 站点基本配置
	if data.ExternalUrl != "" {
		settingsToUpdate["externalUrl"] = data.ExternalUrl
	}
	if data.Swagger != "" {
		settingsToUpdate["Swagger"] = data.Swagger
	}

	// 安全设置
	if data.Mfa != "" {
		settingsToUpdate["mfa"] = data.Mfa
	}
	if data.Issuer != "" {
		settingsToUpdate["issuer"] = data.Issuer
	}
	if data.Secret != "" {
		settingsToUpdate["secret"] = data.Secret
	}
	if data.TokenExpiresTime != "" {
		settingsToUpdate["TokenExpiresTime"] = data.TokenExpiresTime
	}

	// LDAP 设置
	if data.LdapAddress != "" {
		settingsToUpdate["ldapAddress"] = data.LdapAddress
	}
	if data.LdapBindDn != "" {
		settingsToUpdate["ldapBindDn"] = data.LdapBindDn
	}
	if data.LdapBindPassword != "" {
		cipherText, _ := utils.Encrypt(data.LdapBindPassword)
		settingsToUpdate["ldapBindPassword"] = cipherText
	}
	if data.LdapSearchDn != "" {
		settingsToUpdate["ldapSearchDn"] = data.LdapSearchDn
	}
	if data.LdapFilterAttribute != "" {
		settingsToUpdate["ldapFilterAttribute"] = data.LdapFilterAttribute
	}
	if data.LdapUserPasswordExpireDays != "" {
		settingsToUpdate["ldapUserPasswordExpireDays"] = data.LdapUserPasswordExpireDays
	}

	// 用户密码策略
	if data.PasswordExpireDays != "" {
		settingsToUpdate["passwordExpireDays"] = data.PasswordExpireDays
	}
	if data.PasswordLength != "" {
		settingsToUpdate["passwordLength"] = data.PasswordLength
	}
	if data.PasswordComplexity != "" {
		settingsToUpdate["passwordComplexity"] = data.PasswordComplexity
	}
	if data.PasswordExpiryReminderDays != "" {
		settingsToUpdate["passwordExpiryReminderDays"] = data.PasswordExpiryReminderDays
	}
	if data.PasswordMailResetOff != "" {
		settingsToUpdate["passwordMailResetOff"] = data.PasswordMailResetOff
	}

	// 邮件配置
	if data.MailAddress != "" {
		settingsToUpdate["mailAddress"] = data.MailAddress
	}
	if data.MailPort != "" {
		settingsToUpdate["mailPort"] = data.MailPort
	}
	if data.MailForm != "" {
		settingsToUpdate["mailForm"] = data.MailForm
	}
	if data.MailPassword != "" {
		cipherText, _ := utils.Encrypt(data.MailPassword)
		settingsToUpdate["mailPassword"] = cipherText
	}

	// 短信配置
	if data.SmsAppSecret != "" {
		cipherText, _ := utils.Encrypt(data.SmsAppSecret)
		settingsToUpdate["smsAppSecret"] = cipherText
	}
	if data.SmsAppKey != "" {
		settingsToUpdate["smsAppKey"] = data.SmsAppKey
	}
	if data.SmsProvider != "" {
		settingsToUpdate["smsProvider"] = data.SmsProvider
	}
	if data.SmsSignature != "" {
		settingsToUpdate["smsSignature"] = data.SmsSignature
	}
	if data.SmsEndpoint != "" {
		settingsToUpdate["smsEndpoint"] = data.SmsEndpoint
	}
	if data.SmsSender != "" {
		if data.SmsProvider == "aliyun" {
			settingsToUpdate["smsSender"] = nil
		} else {
			settingsToUpdate["smsSender"] = data.SmsSender
		}
	}
	if data.SmsCallbackUrl != "" {
		if data.SmsProvider == "aliyun" {
			settingsToUpdate["smsCallbackUrl"] = nil
		} else {
			settingsToUpdate["smsCallbackUrl"] = data.SmsCallbackUrl
		}
	}
	if data.SmsTemplateId != "" {
		settingsToUpdate["smsTemplateId"] = data.SmsTemplateId
	}

	// 钉钉配置
	if data.DingdingAppKey != "" {
		settingsToUpdate["dingdingAppKey"] = data.DingdingAppKey
	}
	if data.DingdingAppSecret != "" {
		cipherText, _ := utils.Encrypt(data.DingdingAppSecret)
		settingsToUpdate["dingdingAppSecret"] = cipherText
	}

	// 飞书配置
	if data.FeishuAppId != "" {
		settingsToUpdate["feishuAppId"] = data.FeishuAppId
	}
	if data.FeishuAppSecret != "" {
		cipherText, _ := utils.Encrypt(data.FeishuAppSecret)
		settingsToUpdate["feishuAppSecret"] = cipherText
	}

	// 企业微信配置
	if data.WechatAgentId != "" {
		settingsToUpdate["wechatAgentId"] = data.WechatAgentId
	}
	if data.WechatSecret != "" {
		cipherText, _ := utils.Encrypt(data.WechatSecret)
		settingsToUpdate["wechatSecret"] = cipherText
	}
	if data.WechatCorpId != "" {
		settingsToUpdate["wechatCorpId"] = data.WechatCorpId
	}

	// 开启事务
	tx := global.MySQLClient.Begin()

	// 批量更新
	result, err := dao.Settings.UpdateSettings(tx, settingsToUpdate)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 重新加载配置
	if err := db.InitConfig(global.MySQLClient); err != nil {
		logger.Warn("配置加载失败：" + err.Error())
	}

	return result, nil
}

// MailTest 发送邮件测试
func (s *settings) MailTest(receiver string) error {

	// 生成HTML内容
	htmlBody := TestHTML()

	// 发送
	if err := mail.Email.SendMsg([]string{receiver}, nil, nil, "配置测试", htmlBody, "html"); err != nil {
		return err
	}

	return nil
}

// SmsTest 发送短信测试
func (s *settings) SmsTest(username string) error {

	// 定义用户匹配条件
	conditions := map[string]interface{}{
		"username": username,
	}

	// 在本地数据库中查找匹配的用户
	user, err := dao.User.GetUser(conditions)
	if err != nil {
		return err
	}

	// 发送短信
	if _, err := SMS.SMSSend(user.PhoneNumber, "测试"); err != nil {
		return err
	}

	return nil
}

// LoginTest LDAP 用户登录测试
func (s *settings) LoginTest(username, password string) error {

	// 认证测试
	if _, err := AD.LDAPUserAuthentication(username, password); err != nil {
		return err
	}

	return nil
}

// CertTest 密钥证书测试
func (s *settings) CertTest(certificate, privateKey, publicKey string) error {

	// 解析私钥
	privateKeyBlock, _ := pem.Decode([]byte(privateKey))
	if privateKeyBlock == nil || privateKeyBlock.Type != "PRIVATE KEY" {
		return errors.New("无效的私钥")
	}
	privInterface, err := x509.ParsePKCS8PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return errors.New(fmt.Sprintf("无效的私钥: %v", err))
	}

	// 确保私钥为 RSA
	priv, ok := privInterface.(*rsa.PrivateKey)
	if !ok {
		return errors.New("私钥不是 RSA 类型")
	}

	// 解析公钥
	publicKeyBlock, _ := pem.Decode([]byte(publicKey))
	if publicKeyBlock == nil {
		return errors.New("无效的公钥")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return errors.New(fmt.Sprintf("无效的公钥: %v", err))
	}

	// 确保公钥为 RSA
	pub, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return errors.New("公钥不是 RSA 类型")
	}

	// 验证私钥和公钥是否匹配
	privPub := priv.Public()
	privPubKey, ok := privPub.(*rsa.PublicKey)
	if !ok {
		return errors.New(fmt.Sprintf("无法从私钥中提取公钥: %v", err))
	}
	if privPubKey.N.Cmp(pub.N) != 0 || privPubKey.E != pub.E {
		return errors.New("私钥和公钥不匹配")
	}

	// 解析证书
	certBlock, _ := pem.Decode([]byte(certificate))
	if certBlock == nil || certBlock.Type != "CERTIFICATE" {
		return errors.New("无效的证书")
	}
	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return errors.New(fmt.Sprintf("无效的证书: %v", err))
	}

	// 验证证书中的公钥是否匹配
	certPub, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return errors.New("证书中提取的公钥类型不是 RSA")
	}
	if certPub.N.Cmp(pub.N) != 0 || certPub.E != pub.E {
		return errors.New("证书中提取的公钥和提供的公钥不匹配")
	}

	return nil
}

// CertUpdate 证书及密钥更新
func (s *settings) CertUpdate(certificate, privateKey, publicKey string) (map[string]interface{}, error) {

	var (
		settingsToUpdate = map[string]interface{}{
			"certificate": certificate,
			"publicKey":   publicKey,
			"privateKey":  privateKey,
		}
		authUsers      []model.AuthUser
		accounts       []model.Account
		cfgs           []model.Settings
		domainProvider []model.DomainServiceProvider
		keys           = []string{"ldapBindPassword", "mailPassword", "smsAppSecret", "dingdingAppSecret", "feishuAppSecret", "wechatSecret"}
	)

	// 开启事务
	tx := global.MySQLClient.Begin()

	// 域名提供商密码更新
	if err := global.MySQLClient.Find(&domainProvider).Error; err != nil {
		return nil, err
	}
	for _, provider := range domainProvider {

		if provider.AccessKey == nil && provider.SecretKey == nil && provider.IamPassword == nil {
			continue
		}

		if provider.AccessKey != nil {
			ak := *provider.AccessKey
			// 获取明文AK
			plaintextAK, err := utils.Decrypt(ak)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			// 使用新密钥加密
			newCiphertextAK, err := utils.EncryptWithPublicKey(plaintextAK, publicKey)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			// 保存
			if err := tx.Model(&provider).Where("id = ?", provider.Id).Updates(map[string]interface{}{"access_key": newCiphertextAK}).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		if provider.SecretKey != nil {
			sk := *provider.SecretKey
			// 获取明文SK
			plaintextSK, err := utils.Decrypt(sk)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			// 使用新密钥加密
			newCiphertextSK, err := utils.EncryptWithPublicKey(plaintextSK, publicKey)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			// 保存
			if err := tx.Model(&provider).Where("id = ?", provider.Id).Updates(map[string]interface{}{"secret_key": newCiphertextSK}).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		if provider.IamPassword != nil {
			iamPassword := *provider.IamPassword
			// 获取明文密码
			plaintextPassword, err := utils.Decrypt(iamPassword)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			// 使用新密钥加密
			newCiphertextPassword, err := utils.EncryptWithPublicKey(plaintextPassword, publicKey)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			// 保存
			if err := tx.Model(&provider).Where("id = ?", provider.Id).Updates(map[string]interface{}{"iam_password": newCiphertextPassword}).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	// 用户密码更新
	if err := global.MySQLClient.Find(&authUsers).Error; err != nil {
		return nil, err
	}
	for _, user := range authUsers {
		// 获取明文密码
		plaintext, err := utils.Decrypt(user.Password)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// 使用新密钥加密
		newCiphertext, err := utils.EncryptWithPublicKey(plaintext, publicKey)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// 保存
		if err := tx.Model(&user).Where("id = ?", user.ID).Updates(map[string]interface{}{"password": newCiphertext}).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// 账号资产密码更新
	if err := global.MySQLClient.Find(&accounts).Error; err != nil {
		return nil, err

	}
	for _, account := range accounts {
		// 获取明文密码
		plaintext, err := utils.Decrypt(account.Password)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// 使用新密钥加密
		newCiphertext, err := utils.EncryptWithPublicKey(plaintext, publicKey)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// 保存
		if err := tx.Model(&account).Where("id = ?", account.ID).Updates(map[string]interface{}{"password": newCiphertext}).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// 配置信息加密数据更新
	if err := global.MySQLClient.Where("`key` IN ?", keys).Find(&cfgs).Error; err != nil {
		return nil, err
	}
	for _, setting := range cfgs {

		// 未配置则跳过
		if setting.Value == nil {
			continue
		}

		// 解密当前值
		plaintext, err := utils.Decrypt(*setting.Value)
		if err != nil {
			tx.Rollback()
			return nil, errors.New(fmt.Sprintf("failed to decrypt value for key %v: %v", setting.Key, err))
		}

		// 重新加密
		newCiphertext, err := utils.EncryptWithPublicKey(plaintext, publicKey)
		if err != nil {
			tx.Rollback()
			return nil, errors.New(fmt.Sprintf("failed to decrypt value for key %v: %v", setting.Key, err))
		}

		// 更新值
		setting.Value = &newCiphertext
		if err := tx.Save(&setting).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	for key, value := range settingsToUpdate {
		strValue := fmt.Sprintf("%v", value)
		settingsToUpdate[key] = strValue
	}

	// 证书及密钥配置更新
	result, err := dao.Settings.UpdateSettings(tx, settingsToUpdate)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 重新加载配置
	if err := db.InitConfig(global.MySQLClient); err != nil {
		logger.Warn("配置加载失败：" + err.Error())
	}

	return result, nil
}

func TestHTML() string {

	issuer := config.Conf.Settings["issuer"].(string)

	return fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>测试</title>
		</head>
		<body>
			<p>亲爱的同事：</p>
			<p>恭喜您，收到此邮件表示系统邮件配置正确。</p>
			<br>
			<p>此致，<br>%s</p>
			<p style="color: red">此邮件为测试邮件，请不要回复此邮件。</p>
		</body>
		</html>
	`, issuer)
}
