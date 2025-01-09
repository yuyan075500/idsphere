package service

import (
	"fmt"
	"github.com/wonderivan/logger"
	"mime/multipart"
	"ops-api/dao"
	"ops-api/db"
	"ops-api/global"
	"ops-api/model"
	"ops-api/utils"
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
