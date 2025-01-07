package dao

import (
	"gorm.io/gorm"
	"ops-api/global"
	"ops-api/model"
)

var Settings settings

type settings struct{}

// GetAllSettings 获取所有配置
func (s *settings) GetAllSettings() ([]model.Settings, error) {
	var settings []model.Settings
	// 获取所有配置，排队敏感信息
	if err := global.MySQLClient.Not("`key` IN ?", []string{"ldapBindPassword", "mailPassword"}).Find(&settings).Error; err != nil {
		return nil, err
	}
	return settings, nil
}

// GetSettingByKey 获取单个配置
func (s *settings) GetSettingByKey(key string) (*model.Settings, error) {
	var setting model.Settings
	if err := global.MySQLClient.Where("`key` = ?", key).First(&setting).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}

// UpdateSetting 更新单个配置
func (s *settings) UpdateSetting(key, value string) error {
	return global.MySQLClient.Model(&model.Settings{}).Where("`key` = ?", key).Update("value", value).Error
}

// UpdateSettings 批量更新多个配置项
func (s *settings) UpdateSettings(db *gorm.DB, settings map[string]interface{}) (map[string]interface{}, error) {
	// 批量更新
	for key, value := range settings {
		if err := db.Model(&model.Settings{}).Where("`key` = ?", key).Update("value", value).Error; err != nil {
			return nil, err
		}
	}

	// 查询更新后的配置值
	updatedSettings := make(map[string]interface{})
	for key := range settings {
		var setting model.Settings
		if err := db.Where("`key` = ?", key).First(&setting).Error; err != nil {
			return nil, err
		}
		updatedSettings[key] = setting.Value
	}

	return updatedSettings, nil
}

func getKeys(settings map[string]interface{}) []string {
	keys := make([]string, 0, len(settings))
	for key := range settings {
		keys = append(keys, key)
	}
	return keys
}
