package service

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dingtalkcontact_1_0 "github.com/alibabacloud-go/dingtalk/contact_1_0"
	dingtalkoauth2_1_0 "github.com/alibabacloud-go/dingtalk/oauth2_1_0"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"ops-api/config"
	"ops-api/utils"
)

// DingTalkClient 钉钉API客户端
type DingTalkClient struct {
	OAuthClient   *dingtalkoauth2_1_0.Client
	ContactClient *dingtalkcontact_1_0.Client
}

// NewDingTalkClient 创建一个新的钉钉API客户端
func NewDingTalkClient() (*DingTalkClient, error) {
	conf := &openapi.Config{
		Protocol: tea.String("https"),
		RegionId: tea.String("central"),
	}

	oauthClient, err := dingtalkoauth2_1_0.NewClient(conf)
	if err != nil {
		return nil, err
	}

	contactClient, err := dingtalkcontact_1_0.NewClient(conf)
	if err != nil {
		return nil, err
	}

	return &DingTalkClient{
		OAuthClient:   oauthClient,
		ContactClient: contactClient,
	}, nil
}

// GetUserAccessToken 获取用户Token
func (client *DingTalkClient) GetUserAccessToken(code string) (userAccessToken string, err error) {

	var (
		dingdingAppKey    = config.Conf.Settings["dingdingAppKey"].(string)
		dingdingAppSecret = config.Conf.Settings["dingdingAppSecret"].(string)
	)

	// 解密
	appSecret, _ := utils.Decrypt(dingdingAppSecret)

	request := &dingtalkoauth2_1_0.GetUserTokenRequest{
		ClientSecret: tea.String(appSecret),
		ClientId:     tea.String(dingdingAppKey),
		Code:         tea.String(code),
		GrantType:    tea.String("authorization_code"),
	}

	response, err := client.OAuthClient.GetUserToken(request)
	if err != nil {
		return "", err
	}

	return tea.StringValue(response.Body.AccessToken), nil
}

// GetUserInfo 获取用户信息
func (client *DingTalkClient) GetUserInfo(userAccessToken string) (*dingtalkcontact_1_0.GetUserResponse, error) {
	headers := &dingtalkcontact_1_0.GetUserHeaders{
		XAcsDingtalkAccessToken: tea.String(userAccessToken),
	}

	userInfo, err := client.ContactClient.GetUserWithOptions(tea.String("me"), headers, &util.RuntimeOptions{})
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}
