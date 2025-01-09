package service

import (
	"errors"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/work"
	"ops-api/config"
)

// WechatClient 企业微信API客户端
type WechatClient struct {
	WeComApp *work.Work
}

// NewWeChatClient 企业微信实例化，参考文档：https://powerwechat.artisan-cloud.com/zh/wecom/
func NewWeChatClient() (*WechatClient, error) {

	externalUrl := config.Conf.Settings["externalUrl"].(string)

	// 读取配置信息
	corpId := config.Conf.Wechat.CorpId   // 企业微信的企业ID。
	agentId := config.Conf.Wechat.AgentId // 内部应用的AgentId
	secret := config.Conf.Wechat.Secret   // 内部应用的Secret

	// 客户端初始化
	WeComApp, err := work.NewWork(&work.UserConfig{
		CorpID:  corpId,
		AgentID: agentId,
		Secret:  secret,
		OAuth: work.OAuth{
			Callback: externalUrl,
			Scopes:   nil,
		},
		HttpDebug: false,
	})

	if err != nil {
		return nil, err
	}

	return &WechatClient{
		WeComApp: WeComApp,
	}, nil
}

// GetUserId 获取用户ID
func (client *WechatClient) GetUserId(code string) (string, error) {

	user, err := client.WeComApp.OAuth.Provider.GetUserInfo(code)
	if err != nil {
		return "", err
	}
	if user.ErrCode != 0 {
		return "", errors.New(user.ErrMSG)
	}

	// 不允许外部用户登录
	if user.ExternalUserID != "" {
		return "", errors.New("禁止本企业用户访问")
	}

	return user.UserID, nil
}
