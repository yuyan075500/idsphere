package public_cloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	sdkerror "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	domain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/domain/v20180808"
	"strconv"
	"time"
)

// TencentClient 腾讯云相关的 SDK 请求客户端
type TencentClient struct {
	DomainClient *domain.Client
	DNSClient    *dnspod.Client
}

// CreateTencentClient 创建腾讯云客户端
func CreateTencentClient(accessKey, secretKey string) (*TencentClient, error) {
	credential := common.NewCredential(
		accessKey,
		secretKey,
	)

	// 实例化一个域名 client 选项，可选的，没有特殊需求可以跳过
	domainCpf := profile.NewClientProfile()
	domainCpf.HttpProfile.Endpoint = "domain.tencentcloudapi.com"
	domainClient, _ := domain.NewClient(credential, "", domainCpf)

	// 实例化一个 DNS client 选项，可选的，没有特殊需求可以跳过
	DNSCpf := profile.NewClientProfile()
	DNSCpf.HttpProfile.Endpoint = "dnspod.tencentcloudapi.com"
	dnsClient, _ := dnspod.NewClient(credential, "", DNSCpf)

	return &TencentClient{
		DomainClient: domainClient,
		DNSClient:    dnsClient,
	}, nil
}

// SyncDomains 域名同步
func (client *TencentClient) SyncDomains(serviceProviderID uint) ([]DomainList, error) {

	var (
		domains []DomainList
		offset  = uint64(0)
		limit   = uint64(100)
	)

	for {
		// 构造请求
		request := domain.NewDescribeDomainNameListRequest()
		request.Offset = common.Uint64Ptr(offset)
		request.Limit = common.Uint64Ptr(limit)

		// 发送请求
		response, err := client.DomainClient.DescribeDomainNameList(request)
		if err != nil {
			if err, ok := err.(*sdkerror.TencentCloudSDKError); ok {
				return nil, errors.New(err.GetMessage())
			}
			return nil, err
		}

		// 解析 JSON
		var respData struct {
			Response struct {
				TotalCount uint64 `json:"TotalCount"`
				DomainSet  []struct {
					DomainName     string `json:"DomainName"`
					CreationDate   string `json:"CreationDate"`
					ExpirationDate string `json:"ExpirationDate"`
				} `json:"DomainSet"`
			} `json:"Response"`
		}

		if err := json.Unmarshal([]byte(response.ToJsonString()), &respData); err != nil {
			return nil, err
		}

		// 解析数据
		for _, d := range respData.Response.DomainSet {
			createdAt, _ := time.Parse("2006-01-02", d.CreationDate)
			expiredAt, _ := time.Parse("2006-01-02", d.ExpirationDate)

			domains = append(domains, DomainList{
				Name:                    d.DomainName,
				RegistrationAt:          &createdAt,
				ExpirationAt:            &expiredAt,
				DomainServiceProviderID: serviceProviderID,
			})
		}

		// 如果获取的数量少于 limit，说明数据获取完了
		if uint64(len(respData.Response.DomainSet)) < limit {
			break
		}

		// 更新 offset 进行下一页查询
		offset += limit
	}

	return domains, nil
}

// GetDns 获取域名 DNS 记录
func (client *TencentClient) GetDns(pageNum, pageSize int64, domainName, keyWord string) (*DnsList, error) {

	var dnsItems []*DNS

	// 构造请求
	request := dnspod.NewDescribeRecordListRequest()
	request.Domain = common.StringPtr(domainName)
	request.Keyword = common.StringPtr(keyWord)
	request.Offset = common.Uint64Ptr(uint64((pageNum - 1) * pageSize))
	request.Limit = common.Uint64Ptr(uint64(pageSize))

	// 发送请求
	response, err := client.DNSClient.DescribeRecordList(request)
	if err != nil {
		if err, ok := err.(*sdkerror.TencentCloudSDKError); ok {
			return nil, errors.New(err.GetMessage())
		}
		return nil, err
	}

	// 解析返回数据
	recordList := response.Response.RecordList

	for _, record := range recordList {

		// 跳过默认 NS 记录
		if *record.DefaultNS {
			continue
		}

		// 转换格式
		dnsItem := &DNS{
			RR:       *record.Name,
			Type:     *record.Type,
			Value:    *record.Value,
			TTL:      int(*record.TTL),
			Status:   *record.Status,
			CreateAt: formatTime(*record.UpdatedOn),
			Remark:   *record.Remark,
			RecordId: fmt.Sprintf("%d", *record.RecordId),
			Priority: int(*record.MX),
		}

		if record.Weight != nil {
			dnsItem.Weight = new(int)
			*dnsItem.Weight = int(*record.Weight)
		} else {
			dnsItem.Weight = nil
		}

		dnsItems = append(dnsItems, dnsItem)
	}

	// 返回结果
	dnsList := DnsList{
		Items: dnsItems,
		Total: int64(*response.Response.RecordCountInfo.TotalCount),
	}

	return &dnsList, nil
}

// AddDns 添加域名 DNS 记录
func (client *TencentClient) AddDns(domainName, rrType, rr, value, remark string, ttl int32, weight *int32, priority int32) error {

	// 构造请求
	request := dnspod.NewCreateRecordRequest()
	request.Domain = common.StringPtr(domainName)
	request.RecordType = common.StringPtr(rrType)
	request.RecordLine = common.StringPtr("默认")
	request.Value = common.StringPtr(value)
	request.SubDomain = common.StringPtr(rr)
	request.TTL = common.Uint64Ptr(uint64(ttl))
	request.Remark = common.StringPtr(remark)

	if weight != nil {
		request.Weight = common.Uint64Ptr(uint64(*weight))
	}

	// 设置 MX 类型的优先级
	if rrType == "MX" {
		request.MX = common.Uint64Ptr(uint64(priority))
	}

	// 发送请求
	_, err := client.DNSClient.CreateRecord(request)
	if err != nil {
		if err, ok := err.(*sdkerror.TencentCloudSDKError); ok {
			return errors.New(err.GetMessage())
		}
		return err
	}

	return nil
}

// UpdateDns 修改域名 DNS 记录
func (client *TencentClient) UpdateDns(domainName, recordId, rrType, rr, value, remark string, ttl int32, weight *int32, priority int32) error {

	// 将 recordId 从字符串转换为 uint64
	recordIdUint, err := strconv.ParseUint(recordId, 10, 64)
	if err != nil {
		return err
	}

	// 构造请求
	request := dnspod.NewModifyRecordRequest()
	request.Domain = common.StringPtr(domainName)
	request.RecordType = common.StringPtr(rrType)
	request.RecordLine = common.StringPtr("默认")
	request.Value = common.StringPtr(value)
	request.RecordId = common.Uint64Ptr(recordIdUint)
	request.SubDomain = common.StringPtr(rr)
	request.TTL = common.Uint64Ptr(uint64(ttl))
	request.Remark = common.StringPtr(remark)

	if weight != nil {
		request.Weight = common.Uint64Ptr(uint64(*weight))
	}

	// 设置 MX 类型的优先级
	if rrType == "MX" {
		request.MX = common.Uint64Ptr(uint64(priority))
	}

	// 发送请求
	_, err = client.DNSClient.ModifyRecord(request)
	if err != nil {
		if err, ok := err.(*sdkerror.TencentCloudSDKError); ok {
			return errors.New(err.GetMessage())
		}
		return err
	}

	return nil
}

// DeleteDns 删除域名 DNS 记录
func (client *TencentClient) DeleteDns(domainName, recordId string) error {

	// 将 recordId 从字符串转换为 uint64
	recordIdUint, err := strconv.ParseUint(recordId, 10, 64)
	if err != nil {
		return err
	}

	// 构造请求
	request := dnspod.NewDeleteRecordRequest()
	request.Domain = common.StringPtr(domainName)
	request.RecordId = common.Uint64Ptr(recordIdUint)

	// 发送请求
	_, err = client.DNSClient.DeleteRecord(request)
	if err != nil {
		if err, ok := err.(*sdkerror.TencentCloudSDKError); ok {
			return errors.New(err.GetMessage())
		}
		return err
	}

	return nil
}

// SetDnsStatus 设置域名 DNS 状态
func (client *TencentClient) SetDnsStatus(domainName, recordId, status string) error {

	// 将 recordId 从字符串转换为 uint64
	recordIdUint, err := strconv.ParseUint(recordId, 10, 64)
	if err != nil {
		return err
	}

	// 构造请求
	request := dnspod.NewModifyRecordStatusRequest()
	request.Domain = common.StringPtr(domainName)
	request.RecordId = common.Uint64Ptr(recordIdUint)
	request.Status = common.StringPtr(status)

	// 发送请求
	_, err = client.DNSClient.ModifyRecordStatus(request)
	if err != nil {
		if err, ok := err.(*sdkerror.TencentCloudSDKError); ok {
			return errors.New(err.GetMessage())
		}
		return err
	}

	return nil
}

// 时间处理
func formatTime(timeStr string) string {
	parsedTime, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		return timeStr
	}
	return parsedTime.Format(time.RFC3339)
}
