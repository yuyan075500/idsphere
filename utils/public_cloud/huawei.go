package public_cloud

import (
	"encoding/json"
	"fmt"
	hwBasic "github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	hwGlobal "github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	dns "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dns/v2"
	dnsModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dns/v2/model"
	dnsRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dns/v2/region"
	iam "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3"
	iamModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"
	iamRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/region"
	"io"
	"net/http"
	"ops-api/dao"
	"ops-api/global"
	"ops-api/utils"
	"strings"
	"time"
)

// HuaweiClient 华为云相关的 SDK 请求客户端
type HuaweiClient struct {
	GetTokenClient *iam.IamClient
	DNSClient      *dns.DnsClient
}

// CreateHuaweiClient 创建华为云客户端
func CreateHuaweiClient(accessKey, secretKey string) (*HuaweiClient, error) {

	// 创建 IAM 客户端
	GlobalAuth, err := hwGlobal.NewCredentialsBuilder().WithAk(accessKey).WithSk(secretKey).SafeBuild()
	if err != nil {
		return nil, err
	}
	iamReg, err := iamRegion.SafeValueOf("cn-east-3")
	if err != nil {
		return nil, err
	}
	iamBuilder := iam.IamClientBuilder().WithRegion(iamReg).WithCredential(GlobalAuth)
	iamHcClient, err := iamBuilder.SafeBuild()
	if err != nil {
		return nil, err
	}
	iamClient := iam.NewIamClient(iamHcClient)

	// 创建 DNS 客户端
	basicAuth, err := hwBasic.NewCredentialsBuilder().WithAk(accessKey).WithSk(secretKey).SafeBuild()
	if err != nil {
		return nil, err
	}
	dnsReg, err := dnsRegion.SafeValueOf("cn-east-3")
	if err != nil {
		return nil, err
	}
	dnsBuilder := dns.DnsClientBuilder().WithRegion(dnsReg).WithCredential(basicAuth)
	dnsHcClient, err := dnsBuilder.SafeBuild()
	if err != nil {
		return nil, err
	}
	dnsClient := dns.NewDnsClient(dnsHcClient)

	return &HuaweiClient{
		GetTokenClient: iamClient,
		DNSClient:      dnsClient,
	}, nil
}

// GetToken 获取用户Token
func (client *HuaweiClient) GetToken(accountName, iamUsername, iamPassword string) (*string, error) {
	request := &iamModel.KeystoneCreateUserTokenByPasswordRequest{}
	domainUser := &iamModel.PwdPasswordUserDomain{
		Name: accountName,
	}
	userPassword := &iamModel.PwdPasswordUser{
		Domain:   domainUser,
		Name:     iamUsername,
		Password: iamPassword,
	}
	passwordIdentity := &iamModel.PwdPassword{
		User: userPassword,
	}
	var listMethodsIdentity = []iamModel.PwdIdentityMethods{
		iamModel.GetPwdIdentityMethodsEnum().PASSWORD,
	}
	identityAuth := &iamModel.PwdIdentity{
		Methods:  listMethodsIdentity,
		Password: passwordIdentity,
	}
	authBody := &iamModel.PwdAuth{
		Identity: identityAuth,
	}
	request.Body = &iamModel.KeystoneCreateUserTokenByPasswordRequestBody{
		Auth: authBody,
	}
	response, err := client.GetTokenClient.KeystoneCreateUserTokenByPassword(request)
	if err != nil {
		return nil, err
	} else {
		return response.XSubjectToken, nil
	}
}

// GetZoneID 获取域名Zone ID
func (client *HuaweiClient) GetZoneID(domainName string) (zoneID string, err error) {

	request := &dnsModel.ListPublicZonesRequest{}

	// 请求的域名
	nameRequest := domainName
	request.Name = &nameRequest

	// 查找模式，equal 表示精确查找
	searchModeRequest := "equal"
	request.SearchMode = &searchModeRequest

	// 发送请求
	response, err := client.DNSClient.ListPublicZones(request)
	if err != nil {
		return "", nil
	}

	// 返回域名（Zone）ID
	for _, zone := range *response.Zones {
		if *zone.Name == fmt.Sprintf("%s.", domainName) {
			return *zone.Id, nil
		}
	}

	return "", nil
}

// SyncDomains 域名同步
func (client *HuaweiClient) SyncDomains(serviceProviderID uint) ([]DomainList, error) {

	var (
		domains []DomainList
		offset  = uint64(0)
		limit   = uint64(100)
	)

	// 获取服务商信息
	provider, err := dao.Domain.GetDomainServiceProviderForID(int(serviceProviderID))
	if err != nil {
		return nil, err
	}
	if provider.AccountName == nil || provider.IamUsername == nil || provider.IamPassword == nil {
		return nil, fmt.Errorf("域名服务商配置信息不完整")
	}
	accountName := *provider.AccountName
	iamUsername := *provider.IamUsername
	iamPassword, _ := utils.Decrypt(*provider.IamPassword)

	// 从缓存用户Token
	token, err := global.RedisClient.Get("hw_iam_user_token").Result()
	if err != nil {
		// 从华为云获取用户Token
		t, err := client.GetToken(accountName, iamUsername, iamPassword)
		if err != nil {
			return nil, err
		}

		token = *t

		// 将Token存入Redis缓存
		if err := global.RedisClient.Set("hw_iam_user_token", *t, 12*time.Hour).Err(); err != nil {
			return nil, err
		}
	}

	for {
		// 创建HTTP请求
		url := fmt.Sprintf("https://domains-external.myhuaweicloud.cn/v2/domains?offset=%d&limit=%d", offset, limit)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		// 设置请求头
		req.Header.Set("X-Auth-Token", token)
		req.Header.Set("Content-Type", "application/json")

		// 发送请求
		c := &http.Client{}
		resp, err := c.Do(req)
		if err != nil {
			return nil, err
		}

		// 读取响应体
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		// 解析 JSON
		var respData struct {
			Domains []struct {
				DomainName   string `json:"domain_name"`
				RegisterDate string `json:"register_date"`
				ExpireDate   string `json:"expire_date"`
			} `json:"domains"`
			Total int `json:"total"`
		}
		if err := json.Unmarshal(body, &respData); err != nil {
			return nil, err
		}

		// 解析数据
		for _, d := range respData.Domains {
			createdAt, _ := time.Parse("2006-01-02", d.RegisterDate)
			expiredAt, _ := time.Parse("2006-01-02", d.ExpireDate)

			domains = append(domains, DomainList{
				Name:                    d.DomainName,
				RegistrationAt:          &createdAt,
				ExpirationAt:            &expiredAt,
				DomainServiceProviderID: serviceProviderID,
			})
		}

		// 如果获取的数量少于 limit，说明数据获取完了
		if uint64(len(respData.Domains)) < limit {
			break
		}

		// 更新 offset 进行下一页查询
		offset += limit
	}

	return domains, nil
}

// GetDns 获取域名 DNS 记录
func (client *HuaweiClient) GetDns(pageNum, pageSize int64, domainName, keyWord string) (*DnsList, error) {

	var dnsItems []*DNS

	// 获取Zone ID
	zoneID, err := client.GetZoneID(domainName)
	if err != nil {
		return nil, err
	}

	request := &dnsModel.ShowRecordSetByZoneRequest{}
	request.ZoneId = zoneID

	// 分页大小
	limitRequest := int32(pageSize)
	request.Limit = &limitRequest

	// 页数
	offsetRequest := int32((pageNum - 1) * pageSize)
	request.Offset = &offsetRequest

	// 过滤
	nameRequest := keyWord
	request.Name = &nameRequest

	// 过滤模式
	searchModeRequest := "like"
	request.SearchMode = &searchModeRequest

	// 发送请求
	response, err := client.DNSClient.ShowRecordSetByZone(request)
	if err != nil {
		return nil, err
	}

	for _, record := range *response.Recordsets {

		// 先去掉末尾的点
		name := strings.TrimSuffix(*record.Name, ".")
		// 如果name等于domainName，说明是根域名，置为空
		if name != domainName {
			// 再去掉domainName
			name = strings.TrimSuffix(name, "."+domainName)
		} else {
			continue
		}

		// 转换格式
		dnsItem := &DNS{
			RR:       name,
			Type:     *record.Type,
			TTL:      int(*record.Ttl),
			Status:   getStatus(record.Status),
			CreateAt: *record.CreatedAt,
			Remark:   getString(record.Description),
			Value:    getValue(record.Records),
			RecordId: *record.Id,
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
		Total: int64(*response.Metadata.TotalCount),
	}

	return &dnsList, nil
}

// AddDns 添加域名 DNS 记录
func (client *HuaweiClient) AddDns(domainName, rrType, rr, value, remark string, ttl int32, weight *int32, priority int32) (recordId string, err error) {

	// 如果value两端没有双引号，则加上双引号
	if !strings.HasPrefix(value, "\"") {
		value = "\"" + value
	}
	if !strings.HasSuffix(value, "\"") {
		value = value + "\""
	}

	// 获取Zone ID
	zoneID, err := client.GetZoneID(domainName)
	if err != nil {
		return "", err
	}

	request := &dnsModel.CreateRecordSetWithLineRequest{}
	request.ZoneId = zoneID

	// 将value 转换为数组
	listRecordsbody := strings.Split(value, ",")

	ttlCreateRecordSetRequestBody := ttl

	descriptionCreateRecordSetRequestBody := remark

	request.Body = &dnsModel.CreateRecordSetWithLineRequestBody{
		Records:     &listRecordsbody,
		Ttl:         &ttlCreateRecordSetRequestBody,
		Type:        rrType,
		Description: &descriptionCreateRecordSetRequestBody,
		Name:        fmt.Sprintf("%s.%s", rr, domainName),
		Weight:      weight,
	}

	res, err := client.DNSClient.CreateRecordSetWithLine(request)
	if err != nil {
		return "", err
	}

	return *res.Id, nil
}

// UpdateDns 修改域名 DNS 记录
func (client *HuaweiClient) UpdateDns(domainName, recordId, rrType, rr, value, remark string, ttl int32, weight *int32, priority int32) error {

	// 获取Zone ID
	zoneID, err := client.GetZoneID(domainName)
	if err != nil {
		return err
	}

	request := &dnsModel.UpdateRecordSetsRequest{}
	request.ZoneId = zoneID
	request.RecordsetId = recordId

	// 将value 转换为数组
	listRecordsbody := strings.Split(value, ",")

	ttlCreateRecordSetRequestBody := ttl

	descriptionCreateRecordSetRequestBody := remark

	request.Body = &dnsModel.UpdateRecordSetsReq{
		Records:     &listRecordsbody,
		Ttl:         &ttlCreateRecordSetRequestBody,
		Type:        rrType,
		Description: &descriptionCreateRecordSetRequestBody,
		Name:        fmt.Sprintf("%s.%s", rr, domainName),
		Weight:      weight,
	}
	_, err = client.DNSClient.UpdateRecordSets(request)
	if err != nil {
		return err
	}

	return nil
}

// DeleteDns 删除域名 DNS 记录
func (client *HuaweiClient) DeleteDns(domainName, recordId string) error {

	// 获取Zone ID
	zoneID, err := client.GetZoneID(domainName)
	if err != nil {
		return err
	}

	request := &dnsModel.DeleteRecordSetsRequest{}
	request.ZoneId = zoneID
	request.RecordsetId = recordId

	_, err = client.DNSClient.DeleteRecordSets(request)
	if err != nil {
		return err
	}

	return nil
}

// SetDnsStatus 设置域名 DNS 状态
func (client *HuaweiClient) SetDnsStatus(domainName, recordId, status string) error {

	request := &dnsModel.SetRecordSetsStatusRequest{}
	request.RecordsetId = recordId

	request.Body = &dnsModel.SetRecordSetsStatusRequestBody{
		Status: status,
	}
	_, err := client.DNSClient.SetRecordSetsStatus(request)
	if err != nil {
		return err
	}

	return nil
}

// 处理可能为空的字符串
func getString(ptr *string) string {
	if ptr != nil {
		return *ptr
	}
	return ""
}

// DNS解析状态处理
func getStatus(status *string) string {
	if *status == "ACTIVE" {
		return "ENABLE"
	}

	return *status
}

// 记录值处理
func getValue(value *[]string) string {
	v := *value
	if len(v) > 1 {
		return strings.Join(v, ",")
	} else if len(v) == 1 {
		return v[0]
	}
	return ""
}
