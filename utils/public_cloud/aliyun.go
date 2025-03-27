package public_cloud

import (
	"encoding/json"
	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	domain20180129 "github.com/alibabacloud-go/domain-20180129/v5/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/wonderivan/logger"
	"ops-api/utils"
	"strings"
	"time"
)

// AliyunClient 结构体封装所有阿里云相关的 SDK 请求客户端
type AliyunClient struct {
	DomainClient *domain20180129.Client
	DnsClient    *alidns20150109.Client
}

// CreateAliyunClient 创建请求客户端
func CreateAliyunClient(accessKey, secretKey string) (*AliyunClient, error) {

	// 定义客户端配置
	config := &openapi.Config{
		AccessKeyId:     tea.String(accessKey),
		AccessKeySecret: tea.String(secretKey),
	}

	config.Endpoint = tea.String("domain.aliyuncs.com")
	domainClient, err := domain20180129.NewClient(config)
	if err != nil {
		return nil, err
	}

	config.Endpoint = tea.String("alidns.cn-hangzhou.aliyuncs.com")
	dnsClient, err := alidns20150109.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &AliyunClient{
		DomainClient: domainClient,
		DnsClient:    dnsClient,
	}, nil
}

// SyncDomains 域名同步
func (client *AliyunClient) SyncDomains(serviceProviderID uint) ([]DomainList, error) {

	// 初始化分页参数
	var (
		pageNum  = int32(1)
		pageSize = int32(100)
		domains  []DomainList
	)

	// 循环获取所有域名
	for {
		// 指定查询的分页和分页大小
		req := &domain20180129.QueryDomainListRequest{
			PageNum:  tea.Int32(pageNum),
			PageSize: tea.Int32(pageSize),
		}

		// 创建查询请求
		resp, err := client.DomainClient.QueryDomainListWithOptions(req, &util.RuntimeOptions{})

		// 错误处理
		if err != nil {
			return nil, handleError(err)
		}

		// 账号下没有域名
		if resp.Body.TotalPageNum == nil || *resp.Body.TotalPageNum == 0 {
			logger.Info("No domains found.")
			break
		}

		// 解析域名
		for _, domain := range resp.Body.Data.Domain {
			domains = append(domains, DomainList{
				Name:                    tea.StringValue(domain.DomainName),
				RegistrationAt:          utils.ParseTime(tea.StringValue(domain.RegistrationDate)),
				ExpirationAt:            utils.ParseTime(tea.StringValue(domain.ExpirationDate)),
				DomainServiceProviderID: serviceProviderID,
			})
		}

		// 如果当前页返回的数量小于 pageSize，则说明没有下一页
		if int32(len(resp.Body.Data.Domain)) < pageSize {
			break
		}
		pageNum++
	}
	return domains, nil
}

// GetDns 获取域名 DNS 记录
func (client *AliyunClient) GetDns(pageNum, pageSize int64, domainName, keyWord string) (*DnsList, error) {

	// 初始化请求参数
	describeDomainRecordsRequest := &alidns20150109.DescribeDomainRecordsRequest{
		DomainName: tea.String(domainName), // 查询的域名
		PageNumber: tea.Int64(pageNum),     // 分页
		PageSize:   tea.Int64(pageSize),    // 分页大小
		SearchMode: tea.String("LIKE"),     // 模糊查询
		KeyWord:    tea.String(keyWord),    // 查询关键字
	}

	// 创建查询请求
	resp, err := client.DnsClient.DescribeDomainRecordsWithOptions(describeDomainRecordsRequest, &util.RuntimeOptions{})

	// 错误处理
	if err != nil {
		return nil, handleError(err)
	}

	// 构造返回的 DnsList
	var dnsList DnsList
	dnsList.Total = tea.Int64Value(resp.Body.TotalCount)

	// 遍历 DomainRecords.Record 数组，填充 DNS 记录
	if resp.Body.DomainRecords != nil && resp.Body.DomainRecords.Record != nil {
		for _, record := range resp.Body.DomainRecords.Record {
			dns := &DNS{
				RR:       tea.StringValue(record.RR),
				Type:     tea.StringValue(record.Type),
				Value:    tea.StringValue(record.Value),
				TTL:      int(tea.Int64Value(record.TTL)),
				Status:   tea.StringValue(record.Status),
				CreateAt: time.Unix(tea.Int64Value(record.CreateTimestamp)/1000, 0).Format(time.RFC3339),
				Remark:   tea.StringValue(record.Remark),
				RecordId: tea.StringValue(record.RecordId),
			}
			if record.Weight != nil {
				dns.Weight = new(int)
				*dns.Weight = int(*record.Weight)
			} else {
				dns.Weight = nil
			}
			if record.Priority != nil {
				dns.Priority = int(tea.Int64Value(record.Priority))
			}
			dnsList.Items = append(dnsList.Items, dns)
		}
	}

	return &dnsList, nil
}

// AddDns 添加域名 DNS 记录
func (client *AliyunClient) AddDns(domainName, rrType, rr, value, remark string, ttl int32, weight *int32, priority int32) (recordId string, err error) {

	// 初始化请求参数
	addDomainRecordRequest := &alidns20150109.AddDomainRecordRequest{
		DomainName: tea.String(domainName),
		RR:         tea.String(rr),
		Type:       tea.String(rrType),
		Value:      tea.String(value),
		TTL:        tea.Int64(int64(ttl)),
	}

	// 设置 MX 类型的优先级
	if rrType == "MX" {
		addDomainRecordRequest.Priority = tea.Int64(int64(priority))
	}

	// 创建新增请求
	res, err := client.DnsClient.AddDomainRecordWithOptions(addDomainRecordRequest, &util.RuntimeOptions{})

	if err != nil {
		return "", handleError(err)
	}
	return *res.Body.RecordId, nil
}

// UpdateDns 修改域名 DNS 记录
func (client *AliyunClient) UpdateDns(domainName, recordId, rrType, rr, value, remark string, ttl int32, weight *int32, priority int32) error {

	// 初始化请求参数
	updateDomainRecordRequest := &alidns20150109.UpdateDomainRecordRequest{
		RecordId: tea.String(recordId),
		RR:       tea.String(rr),
		Type:     tea.String(rrType),
		Value:    tea.String(value),
		TTL:      tea.Int64(int64(ttl)),
	}

	// 设置 MX 类型的优先级
	if rrType == "MX" {
		updateDomainRecordRequest.Priority = tea.Int64(int64(priority))
	}

	// 创建查询请求
	_, err := client.DnsClient.UpdateDomainRecordWithOptions(updateDomainRecordRequest, &util.RuntimeOptions{})

	if err != nil {
		return handleError(err)
	}
	return nil
}

// DeleteDns 删除域名 DNS 记录
func (client *AliyunClient) DeleteDns(domainName, recordId string) error {

	// 初始化请求参数
	deleteDomainRecordRequest := &alidns20150109.DeleteDomainRecordRequest{
		RecordId: tea.String(recordId),
	}

	// 创建查询请求
	_, err := client.DnsClient.DeleteDomainRecordWithOptions(deleteDomainRecordRequest, &util.RuntimeOptions{})

	if err != nil {
		return handleError(err)
	}
	return nil
}

// SetDnsStatus 设置域名 DNS 状态
func (client *AliyunClient) SetDnsStatus(domainName, recordId, status string) error {

	// 初始化请求参数
	setDomainRecordStatusRequest := &alidns20150109.SetDomainRecordStatusRequest{
		RecordId: tea.String(recordId),
		Status:   tea.String(status),
	}

	// 创建查询请求
	_, err := client.DnsClient.SetDomainRecordStatusWithOptions(setDomainRecordStatusRequest, &util.RuntimeOptions{})

	if err != nil {
		return handleError(err)
	}
	return nil
}

// handleError 统一错误处理
func handleError(err error) error {
	var sdkError *tea.SDKError
	if e, ok := err.(*tea.SDKError); ok {
		sdkError = e
	} else {
		sdkError = &tea.SDKError{Message: tea.String(err.Error())}
	}
	logger.Error("Error:", tea.StringValue(sdkError.Message))
	if sdkError.Data != nil {
		var data interface{}
		_ = json.NewDecoder(strings.NewReader(tea.StringValue(sdkError.Data))).Decode(&data)
		logger.Info("Recommend:", data)
	}
	return err
}
