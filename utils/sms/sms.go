package sms

import (
	"encoding/json"
	"errors"
	"ops-api/config"
)

// Sender 发送短信接口
type Sender interface {
	SendSMS(data *SendData, code string) (string, error)
	ProcessResponse(resp string) (smsMsgId string, err error)
}

// SendData 发送短信数据结构体
type SendData struct {
	Username    string `json:"username" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Note        string `json:"note"`
}

// SendDetail 发送详情
type SendDetail struct {
	Date        string
	PhoneNumber string
	BizId       string
}

// Response 短信返回的数据
type Response struct {
	Result      []Result `json:"result"`      // 华为云
	Code        string   `json:"code"`        // 华为云/阿里云短信回执
	Description string   `json:"description"` // 华为云
	Body        Body     `json:"body"`        // 阿里云
	StatusCode  int      `json:"statusCode"`  // 阿里云
}
type Body struct {
	BizId     string `json:"BizId"`
	Code      string `json:"Code"`
	Message   string `json:"Message"`
	RequestId string `json:"RequestId"`
}
type Result struct {
	Total      int    `json:"total"`
	OriginTo   string `json:"originTo"`
	CreateTime string `json:"createTime"`
	From       string `json:"from"`
	SmsMsgId   string `json:"smsMsgId"`
	CountryId  string `json:"countryId"`
	Status     string `json:"status"`
}

// HuaweiSMSSender 华为云短信发送器
type HuaweiSMSSender struct{}

// AliyunSMSSender 阿里云短信发送器
type AliyunSMSSender struct{}

// SendSMS 华为云短信发送
func (s *HuaweiSMSSender) SendSMS(data *SendData, code string) (string, error) {

	var (
		smsSender      = config.Conf.Settings["smsSender"].(string)
		smsTemplateId  = config.Conf.Settings["smsTemplateId"].(string)
		smsCallbackUrl = config.Conf.Settings["smsCallbackUrl"].(string)
		smsSignature   = config.Conf.Settings["smsSignature"].(string)
	)

	return HuaweiSend(
		smsSender,
		smsTemplateId,
		smsCallbackUrl,
		smsSignature,
		data.PhoneNumber,
		code,
	)
}

// SendSMS 阿里云短信发送
func (s *AliyunSMSSender) SendSMS(data *SendData, code string) (string, error) {
	resp, err := AliyunSend(data.PhoneNumber, code)
	if err != nil {
		return "", err
	}
	return *resp, nil
}

// ProcessResponse 华为云响应处理
func (s *HuaweiSMSSender) ProcessResponse(resp string) (string, error) {
	var response Response
	if err := json.Unmarshal([]byte(resp), &response); err != nil {
		return "", err
	}
	if response.Code != "000000" {
		return "", errors.New("短信发送失败，错误码：" + response.Code)
	}

	// SmsMsgId短信唯一标识，在接收短信回调时会使用
	return response.Result[0].SmsMsgId, nil
}

// ProcessResponse 阿里云响应处理
func (s *AliyunSMSSender) ProcessResponse(resp string) (string, error) {
	var response Response
	if err := json.Unmarshal([]byte(resp), &response); err != nil {
		return "", err
	}
	if response.Body.Code != "OK" {
		return "", errors.New("短信发送失败，错误码：" + response.Body.Code)
	}

	// BizId短信唯一标识，在后续可以使用此获取短信发送状态
	return response.Body.BizId, nil
}

// GetSMSSender 获取短信发送器
func GetSMSSender() Sender {

	smsProvider := config.Conf.Settings["smsProvider"].(string)

	switch smsProvider {
	case "huawei":
		return &HuaweiSMSSender{}
	case "aliyun":
		return &AliyunSMSSender{}
	default:
		return nil
	}
}
