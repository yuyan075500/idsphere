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
	message "ops-api/utils/sms"
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
	Certificate                string `json:"certificate"`
	PublicKey                  string `json:"publicKey"`
	PrivateKey                 string `json:"privateKey"`
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
	settings, err := dao.Settings.GetAllSettings()
	if err != nil {
		return nil, err
	}

	// 创建结果字典
	result := make(map[string]interface{})

	for _, setting := range settings {
		if err := setting.ParseValue(); err != nil {
			return nil, err
		}
		result[setting.Key] = setting.ParsedValue
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

	settingsToUpdate := map[string]interface{}{
		"externalUrl":                data.ExternalUrl,
		"mfa":                        data.Mfa,
		"issuer":                     data.Issuer,
		"secret":                     data.Secret,
		"ldapAddress":                data.LdapAddress,
		"ldapBindDn":                 data.LdapBindDn,
		"ldapBindPassword":           data.LdapBindPassword,
		"ldapSearchDn":               data.LdapSearchDn,
		"ldapFilterAttribute":        data.LdapFilterAttribute,
		"ldapUserPasswordExpireDays": data.LdapUserPasswordExpireDays,
		"passwordExpireDays":         data.PasswordExpireDays,
		"passwordLength":             data.PasswordLength,
		"passwordComplexity":         data.PasswordComplexity,
		"passwordExpiryReminderDays": data.PasswordExpiryReminderDays,
		"certificate":                data.Certificate,
		"publicKey":                  data.PublicKey,
		"privateKey":                 data.PrivateKey,
		"mailAddress":                data.MailAddress,
		"mailPort":                   data.MailPort,
		"mailForm":                   data.MailForm,
		"mailPassword":               data.MailPassword,
		"smsAppSecret":               data.SmsAppSecret,
		"smsAppKey":                  data.SmsAppKey,
		"smsProvider":                data.SmsProvider,
		"smsSignature":               data.SmsSignature,
		"smsEndpoint":                data.SmsEndpoint,
		"smsSender":                  data.SmsSender,
		"smsCallbackUrl":             data.SmsCallbackUrl,
		"smsTemplateId":              data.SmsTemplateId,
		"dingdingAppKey":             data.DingdingAppKey,
		"dingdingAppSecret":          data.DingdingAppSecret,
		"feishuAppId":                data.FeishuAppId,
		"feishuAppSecret":            data.FeishuAppSecret,
		"wechatCorpId":               data.WechatCorpId,
		"wechatAgentId":              data.WechatAgentId,
		"wechatSecret":               data.WechatSecret,
		"TokenExpiresTime":           data.TokenExpiresTime,
		"Swagger":                    data.Swagger,
	}

	// 开启事务
	tx := global.MySQLClient.Begin()

	for key, value := range settingsToUpdate {
		if value == "" || value == nil {
			// 删除空值
			delete(settingsToUpdate, key)
		} else {
			strValue := fmt.Sprintf("%v", value)

			// 敏感信息加密
			if key == "ldapBindPassword" || key == "mailPassword" || key == "smsAppSecret" || key == "dingdingAppSecret" || key == "feishuAppSecret" || key == "wechatSecret" {
				// 对密码进行加密
				cipherText, err := utils.Encrypt(strValue)
				if err != nil {
					return nil, err
				}
				strValue = cipherText
			}

			settingsToUpdate[key] = strValue
		}
	}

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
	data := &message.SendData{
		Note:        "测试",
		PhoneNumber: user.PhoneNumber,
		Username:    user.Username,
	}
	if _, err := SMS.SMSSend(data); err != nil {
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
