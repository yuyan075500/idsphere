package service

import (
	"bytes"
	"errors"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"
	"image/png"
	"ops-api/config"
	"ops-api/global"
	"ops-api/middleware"
	"ops-api/model"
)

var MFA mfa

type mfa struct{}

// MFAValidate MFA认证接口请求参数（支持CAS3.0和OAuth2.0）
type MFAValidate struct {
	Username         string `json:"username" binding:"required"`
	Code             string `json:"code" binding:"required"`
	Token            string `json:"token" binding:"required"`
	ResponseType     string `json:"response_type"`      // OAuth2.0客户端：授权类型，固定值：code
	ClientId         string `json:"client_id"`          // OAuth2.0客户端：客户端ID
	RedirectURI      string `json:"redirect_uri"`       // OAuth2.0客户端：重定向URL
	State            string `json:"state"`              // OAuth2.0客户端：客户端状态码
	Scope            string `json:"scope"`              // OAuth2.0客户端：申请权限范围
	Service          string `json:"service"`            // CAS3.0客户端：回调地址
	SAMLRequest      string `json:"SAMLRequest"`        // SAML2客户端：SAMLRequest
	RelayState       string `json:"RelayState"`         // SAML2客户端：客户端状态码
	SigAlg           string `json:"SigAlg"`             // SAML2客户端：签名算法
	Signature        string `json:"Signature"`          // SAML2客户端：签名
	NginxRedirectURI string `json:"nginx_redirect_uri"` // Nginx代理客户端：回调地址
}

// GetGoogleQrcode 生成Google MFA认证二维码
func (m *mfa) GetGoogleQrcode(token string) (image []byte, err error) {

	// 获取登录用户名
	username, err := global.RedisClient.Get(token).Result()
	if err != nil {
		return nil, err
	}

	// 创建TOTP
	issuer := config.Conf.Settings["issuer"].(string)
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: username,
	})
	if err != nil {
		return nil, err
	}

	// 使用TOTP Key获取MFA密钥
	mfaSecret := key.Secret()

	// 使用TOTP Key生成二维码图片
	var buf bytes.Buffer
	img, err := key.Image(256, 256)
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}

	// 将mfaSecret更新至缓存，如果用户MFA检验成功则将mfaSecret与用户进行绑定
	if err := global.RedisClient.Set(token, mfaSecret, 0).Err(); err != nil {
		return nil, err
	}

	return buf.Bytes(), err
}

// GoogleQrcodeValidate Google MFA认证校验
func (m *mfa) GoogleQrcodeValidate(params *MFAValidate) (jwtToken, redirectUri, application string, err error) {

	var (
		user   model.AuthUser
		secret string
	)

	// 获取登录用户信息
	tx := global.MySQLClient.First(&user, "username = ?", params.Username)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return "", "", "", errors.New("用户不存在")
	}

	// 获取Secret，如果用户还没有绑定MFA，则从Redis中获取Secret
	if user.MFACode == nil {
		srt, err := global.RedisClient.Get(params.Token).Result()
		if err != nil {
			return "", "", "", err
		}
		secret = srt
	} else {
		secret = *user.MFACode
	}

	// 校验MFA
	valid := totp.Validate(params.Code, secret)
	if !valid {
		return "", "", "", errors.New("验证码错误")
	}

	// 生成用户Token
	jwtToken, err = middleware.GenerateJWT(user.ID, user.Name, user.Username)
	if err != nil {
		return "", "", "", err
	}

	// 更新用户MFA绑定信息
	if user.MFACode == nil {
		user.MFACode = &secret
		if err := tx.Save(&user).Error; err != nil {
			return "", "", "", err
		}
	}

	// 处理单点登录请求
	if params.SAMLRequest != "" || params.Service != "" || params.ClientId != "" {
		loginParams := &UserLogin{
			ResponseType: params.ResponseType,
			ClientId:     params.ClientId,
			RedirectURI:  params.RedirectURI,
			State:        params.State,
			Scope:        params.Scope,
			Service:      params.Service,
			SAMLRequest:  params.SAMLRequest,
			RelayState:   params.RelayState,
			SigAlg:       params.SigAlg,
			Signature:    params.Signature,
		}

		callbackData, siteName, err := SSO.Login(loginParams, user)
		if err != nil {
			return "", "", siteName, err
		}
		return jwtToken, callbackData, siteName, nil
	}

	return jwtToken, "", "", nil

}
