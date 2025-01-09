package sms

import (
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	openapiutil "github.com/alibabacloud-go/openapi-util/service"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"ops-api/config"
	"ops-api/utils"
)

func CreateApiInfo(action string) (_result *openapi.Params) {
	params := &openapi.Params{
		// 接口名称
		Action: tea.String(action),

		// 接口版本
		Version: tea.String("2017-05-25"),

		// 接口协议
		Protocol: tea.String("HTTPS"),

		// 接口 HTTP 方法
		Method:   tea.String("POST"),
		AuthType: tea.String("AK"),
		Style:    tea.String("RPC"),

		// 接口 PATH
		Pathname: tea.String("/"),

		// 接口请求体内容格式
		ReqBodyType: tea.String("json"),

		// 接口响应体内容格式
		BodyType: tea.String("json"),
	}

	_result = params
	return _result
}

// CreateClient 创建客户端
func CreateClient() (_result *openapi.Client, _err error) {

	var (
		accessKeyId     = config.Conf.Settings["smsAppKey"].(string)
		accessKeySecret = config.Conf.Settings["smsAppSecret"].(string)
		endpoint        = config.Conf.Settings["smsEndpoint"].(string)
	)

	// secret解密
	str, _ := utils.Decrypt(accessKeySecret)

	// 指定客户端配置
	conf := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(str),
		Endpoint:        tea.String(endpoint),
	}

	// 客户端实例化
	_result = &openapi.Client{}
	_result, _err = openapi.NewClient(conf)

	return _result, _err
}

func AliyunSend(receiver, templateParas string) (resp *string, err error) {

	var (
		smsSignature  = config.Conf.Settings["smsSignature"].(string)
		smsTemplateId = config.Conf.Settings["smsTemplateId"].(string)
	)

	// 创建客户端
	client, _err := CreateClient()
	if _err != nil {
		return nil, _err
	}

	params := CreateApiInfo("SendSms")

	// 指定请求参数
	queries := map[string]interface{}{}
	queries["PhoneNumbers"] = tea.String(receiver)
	queries["SignName"] = tea.String(smsSignature)
	queries["TemplateCode"] = tea.String(smsTemplateId)
	queries["TemplateParam"] = tea.String(fmt.Sprintf("{\"code\":\"%s\"}", templateParas))

	// 指定运行时选项
	runtime := &util.RuntimeOptions{}

	// 创建API请求
	request := &openapi.OpenApiRequest{
		Query: openapiutil.Query(queries),
	}

	// 请求发送短信
	result, _err := client.CallApi(params, request, runtime)
	if _err != nil {
		return nil, _err
	}

	return util.ToJSONString(result), nil
}

// GetSMSReceipt 获取短信回执
func GetSMSReceipt(phoneNumber, bizId, sendData string) (resp *string, err error) {
	// 创建客户端
	client, _err := CreateClient()
	if _err != nil {
		return nil, _err
	}

	params := CreateApiInfo("QuerySendDetails")

	// 指定请求参数
	queries := map[string]interface{}{}
	queries["PhoneNumber"] = tea.String(phoneNumber)
	queries["BizId"] = tea.String(bizId)
	queries["SendDate"] = tea.String(sendData)
	queries["PageSize"] = tea.Int(1)
	queries["CurrentPage"] = tea.Int(1)

	// 指定运行时选项
	runtime := &util.RuntimeOptions{}

	// 创建API请求
	request := &openapi.OpenApiRequest{
		Query: openapiutil.Query(queries),
	}

	// 请求回执

	result, _err := client.CallApi(params, request, runtime)
	if _err != nil {
		return nil, _err
	}

	return util.ToJSONString(result), nil
}
