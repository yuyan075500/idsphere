package service

import (
	"encoding/json"
	"gorm.io/gorm"
	"ops-api/config"
	"ops-api/dao"
	"ops-api/model"
	messages "ops-api/utils/sms"
)

var Audit audit

type audit struct{}

type Result struct {
	Total      int    `json:"total"`
	OriginTo   string `json:"originTo"`
	CreateTime string `json:"createTime"`
	From       string `json:"from"`
	SmsMsgId   string `json:"smsMsgId"`
	CountryId  string `json:"countryId"`
	Status     string `json:"status"`
}

// AliyunSMSReceipt 阿里云短信回执
type AliyunSMSReceipt struct {
	Body       ResponseBody      `json:"body"`
	Headers    map[string]string `json:"headers"`
	StatusCode int               `json:"statusCode"`
}
type SmsSendDetailDTOs struct {
	SmsSendDetailDTO []SmsSendDetailDTO `json:"SmsSendDetailDTO"`
}
type SmsSendDetailDTO struct {
	Content      string `json:"Content"`
	ErrCode      string `json:"ErrCode"`
	PhoneNum     string `json:"PhoneNum"`
	ReceiveDate  string `json:"ReceiveDate"`
	SendDate     string `json:"SendDate"`
	SendStatus   int    `json:"SendStatus"`
	TemplateCode string `json:"TemplateCode"`
}
type ResponseBody struct {
	Code              string            `json:"Code"`
	Message           string            `json:"Message"`
	RequestId         string            `json:"RequestId"`
	SmsSendDetailDTOs SmsSendDetailDTOs `json:"SmsSendDetailDTOs"`
	TotalCount        int               `json:"TotalCount"`
}

// GetSMSReceipt 获取短信回执
func (a *audit) GetSMSReceipt(smsId int) (err error) {

	smsProvider := config.Conf.Settings["smsProvider"].(string)

	// 华为云不需要
	if smsProvider != "aliyun" {
		return nil
	}

	// 定义匹配条件
	conditions := map[string]interface{}{
		"id": smsId,
	}

	// 查找短信记录
	smsRecord, err := dao.Audit.GetSendDetail(conditions)
	if err != nil {
		return err
	}

	// 获取短信回执
	date := smsRecord.CreatedAt
	receipt, err := messages.GetSMSReceipt(smsRecord.Receiver, smsRecord.SmsMsgId, date.Format("20060102"))
	if err != nil {
		return err
	}

	// 对数据进行解析
	var response AliyunSMSReceipt
	err = json.Unmarshal([]byte(*receipt), &response)
	if err != nil {
		return
	}

	// 处理回执信息
	callback := &dao.Callback{
		SmsMsgId:  smsRecord.SmsMsgId,
		ErrorCode: "",
	}
	if response.Body.Code == "OK" {
		// 获取短信回执内容
		for _, detail := range response.Body.SmsSendDetailDTOs.SmsSendDetailDTO {
			if detail.SendStatus == 1 {
				callback.Status = "等待回执"
			}

			if detail.SendStatus == 2 {
				callback.Status = "发送失败"
			}

			if detail.SendStatus == 3 {
				callback.Status = "接收成功"
			}

		}
	} else {
		callback.Status = "发送失败"
		callback.ErrorCode = response.Body.Code
	}

	// 将回调数据写入数据库
	return dao.Audit.SMSCallback(callback)
}

// GetSMSRecordList 获取短信发送记录
func (a *audit) GetSMSRecordList(receiver string, page, limit int) (data *dao.SMSRecordList, err error) {
	data, err = dao.Audit.GetSMSRecordList(receiver, page, limit)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetLoginRecordList 获取系统登录记录
func (a *audit) GetLoginRecordList(name string, page, limit int) (data *dao.LoginRecordList, err error) {
	data, err = dao.Audit.GetLoginRecordList(name, page, limit)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetOplogList 获取系统操作
func (a *audit) GetOplogList(name string, page, limit int) (data *dao.OplogList, err error) {
	data, err = dao.Audit.GetOplogList(name, page, limit)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// AddLoginFailedRecord 新增系统登录失败记录
func (a *audit) AddLoginFailedRecord(tx *gorm.DB, username, userAgent, clientIP, loginMethod, application string, failedReason error) (err error) {

	loginRecord := &model.LogLogin{
		Username:     username,
		SourceIP:     clientIP,
		UserAgent:    userAgent,
		AuthMethod:   loginMethod,
		Application:  application,
		Status:       2,                    // Status=2，表示失败
		FailedReason: failedReason.Error(), // 记录错误原因
	}

	// 记录登录客户端信息
	return dao.Audit.AddLoginRecord(tx, loginRecord)
}

// AddLoginSuccessRecord 新增系统登录成功记录
func (a *audit) AddLoginSuccessRecord(tx *gorm.DB, username, userAgent, clientIP, loginMethod, application string) (err error) {

	// Status=1表示成功
	loginRecord := &model.LogLogin{
		Username:    username,
		SourceIP:    clientIP,
		UserAgent:   userAgent,
		AuthMethod:  loginMethod,
		Application: application,
		Status:      1,
	}

	// 记录登录客户端信息
	return dao.Audit.AddLoginRecord(tx, loginRecord)
}
