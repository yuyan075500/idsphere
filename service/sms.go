package service

import (
	"errors"
	"net/url"
	"ops-api/config"
	"ops-api/dao"
	"ops-api/model"
	"ops-api/utils"
	message "ops-api/utils/sms"
	"strconv"
	"strings"
)

var SMS sms

type sms struct{}

// SMSSend 发送短信
func (s *sms) SMSSend(phoneNumber, note string) (string, error) {

	var (
		smsSignature  = config.Conf.Settings["smsSignature"].(string)
		smsTemplateId = config.Conf.Settings["smsTemplateId"].(string)
	)

	// 定义验证码
	var code = strconv.Itoa(utils.GenerateRandomNumber())

	// 获取短信服务商
	smsSender := message.GetSMSSender()
	if smsSender == nil {
		return "", errors.New("不支持的短信服务提供商")
	}

	// 发送短信
	resp, err := smsSender.SendSMS(phoneNumber, code)
	if err != nil {
		return "", err
	}

	// 处理响应并获取短信唯一标识
	smsMsgId, err := smsSender.ProcessResponse(resp)
	if err != nil {
		return "", err
	}

	// 记录短信发送日志
	smsLog := &model.LogSMS{
		Note:       note,
		Signature:  smsSignature,
		TemplateId: smsTemplateId,
		Receiver:   phoneNumber,
		Status:     "API请求成功",
		SmsMsgId:   smsMsgId,
	}

	if err := dao.Audit.AddSMSRecord(smsLog); err != nil {
		return "", err
	}

	return code, nil
}

// SMSCallback 短信回调
func (s *sms) SMSCallback(data string) error {

	// 对回调返回的数据进行处理
	ss, _ := url.QueryUnescape(data)
	params := strings.Split(ss, "&")
	keyValues := make(map[string]string)
	for i := range params {
		temp := strings.Split(params[i], "=")
		keyValues[temp[0]] = temp[1]
	}

	// 将数据与结构体进行绑定
	callback := &dao.Callback{
		Status:    "接收成功",
		SmsMsgId:  keyValues["smsMsgId"],
		ErrorCode: "",
	}
	if keyValues["status"] != "DELIVRD" {
		callback.Status = "发送失败"
		callback.ErrorCode = keyValues["status"]
	}

	// 将回调数据写入数据库
	return dao.Audit.SMSCallback(callback)
}
