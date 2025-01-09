package model

import (
	"encoding/json"
	"ops-api/config"
	"ops-api/utils"
	"strconv"
	"time"
)

type Settings struct {
	Id          uint        `json:"id" gorm:"primaryKey;autoIncrement"`
	Key         string      `json:"key" gorm:"unique"`
	Value       *string     `json:"value"`
	ValueType   string      `json:"value_type"`
	ParsedValue interface{} `gorm:"-" json:"parsed_value"`
}

// ParseValue 值转换
func (s *Settings) ParseValue() error {

	if s.Value == nil {
		s.ParsedValue = nil
		return nil
	}

	switch s.ValueType {
	case "list":
		var list []string
		if err := json.Unmarshal([]byte(*s.Value), &list); err != nil {
			return err
		}
		s.ParsedValue = list
	case "boolean":
		s.ParsedValue = *s.Value == "true"
	case "int":
		intValue, err := strconv.Atoi(*s.Value)
		if err != nil {
			return err
		}
		s.ParsedValue = intValue
	default:
		if s.Key == "logo" || s.Key == "icon" {
			logoPreview, err := utils.GetPresignedURL(*s.Value, time.Duration(config.Conf.JWT.Expires)*time.Hour)
			if err != nil {
				return err
			}
			s.ParsedValue = logoPreview.String()
			return nil
		} else {
			s.ParsedValue = *s.Value
		}
	}
	return nil
}
