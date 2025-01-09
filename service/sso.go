package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/LoginRadius/go-saml"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwk"
	"net/url"
	"ops-api/config"
	"ops-api/dao"
	"ops-api/middleware"
	"ops-api/model"
	"ops-api/utils"
	"strings"
	"time"
)

var SSO sso

type sso struct{}

var samlPostFormTemplate = utils.GenerateSAMLResponsePostForm()

// OAuthAuthorize OAuth2.0客户端获取授权请求参数
type OAuthAuthorize struct {
	ResponseType string `json:"response_type" binding:"required"`
	ClientId     string `json:"client_id" binding:"required"`
	RedirectURI  string `json:"redirect_uri"`
	State        string `json:"state"`
	Scope        string `json:"scope"`
	Nonce        string `json:"nonce"`
}

// CASAuthorize CAS3.0客户端获取授权请求参数
type CASAuthorize struct {
	Service string `form:"service" binding:"required"`
}

// Token OAuth2.0客户端获取token请求参数
type Token struct {
	GrantType    string `form:"grant_type"`
	Code         string `form:"code"`
	ClientId     string `form:"client_id"`
	RedirectURI  string `form:"redirect_uri"`
	ClientSecret string `form:"client_secret"`
}

// CASServiceValidate CAS3.0客户端票据校验请求参数
type CASServiceValidate struct {
	Service string `form:"service" binding:"required"`
	Ticket  string `form:"ticket" binding:"required"`
}

// SAMLRequest SAML2客户端授权请求参数
type SAMLRequest struct {
	SAMLRequest string `form:"SAMLRequest" binding:"required"` // SAMLRequest数据，通常该数据是DEFLATE压缩 + base64编码，获取此数据需要进行DEFLATE解压缩 + base64解码
	RelayState  string `form:"RelayState"`                     // SP的状态信息，防止跨站请求伪造攻击，功能与OAuth2.0客户端的state功能相同
	SigAlg      string `form:"SigAlg"`                         // 签名使用的算法
	Signature   string `form:"Signature"`                      // 签名，用于验证SP的身份，但需要配置SP的公钥
}

// ParseSPMetadata 获取SP Metadata信息请求参数
type ParseSPMetadata struct {
	SPMetadataURL string `json:"sp_metadata_url" binding:"required"`
}

// SPMetadata 返回给前端的SP Metadata数据
type SPMetadata struct {
	EntityID    string `json:"entity_id"`
	Certificate string `json:"certificate"`
}

// SAMLResponse IDP返回给浏览器的SAMLResponse数据
type SAMLResponse struct {
	URL          string
	SAMLResponse string
	RelayState   string
}

// ResponseToken 返回给OAuth2.0客户端客户端的Token信息
type ResponseToken struct {
	IdToken      string `json:"id_token"`
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

// CASServiceResponse CAS3.0客户端返回给客户端的用户信息
type CASServiceResponse struct {
	XMLName               xml.Name               `xml:"cas:serviceResponse"`
	Xmlns                 string                 `xml:"xmlns:cas,attr"`
	AuthenticationSuccess *AuthenticationSuccess `xml:"cas:authenticationSuccess"`
}
type AuthenticationSuccess struct {
	User       string     `xml:"cas:user"`
	Attributes Attributes `xml:"cas:attributes"`
}
type Attributes struct {
	Id          uint   `xml:"id"`
	Name        string `xml:"name"`
	Username    string `xml:"username"`
	Email       string `xml:"email"`
	PhoneNumber string `xml:"phone_number"`
}

// ResponseUserinfo 返回给客户端的用户信息
type ResponseUserinfo struct {
	Id                uint   `json:"id"`
	Name              string `json:"name"`               // 用户姓名
	Username          string `json:"username"`           // 用户名
	PreferredUsername string `json:"preferred_username"` // 首选用户名
	Email             string `json:"email"`              // 邮箱地址
	PhoneNumber       string `json:"phone_number"`       // 电话号码
}

type Claims struct {
	PreferredUsername string `json:"preferred_username"`
	NickName          string `json:"nickName"`
}

// OIDCConfig 返回给前端的OIDC配置信息
type OIDCConfig struct {
	Issuer                           string   `json:"issuer"`
	AuthorizationEndpoint            string   `json:"authorization_endpoint"`
	TokenEndpoint                    string   `json:"token_endpoint"`
	UserInfoEndpoint                 string   `json:"userinfo_endpoint"`
	JwksURI                          string   `json:"jwks_uri"`
	ScopesSupported                  []string `json:"scopes_supported"`
	ResponseTypesSupported           []string `json:"response_types_supported"`
	GrantTypesSupported              []string `json:"grant_types_supported"`
	SubjectTypesSupported            []string `json:"subject_types_supported"`
	IDTokenSigningAlgValuesSupported []string `json:"id_token_signing_alg_values_supported"`
}

// GetOIDCConfig 获取OIDC配置信息
func (s *sso) GetOIDCConfig() (configuration *OIDCConfig, err error) {
	externalUrl := config.Conf.Settings["externalUrl"].(string)
	var config = &OIDCConfig{
		Issuer:                           externalUrl,
		AuthorizationEndpoint:            externalUrl + "/login",
		TokenEndpoint:                    externalUrl + "/api/v1/sso/oauth/token",
		UserInfoEndpoint:                 externalUrl + "/api/v1/sso/oauth/userinfo",
		JwksURI:                          externalUrl + "/api/v1/sso/oidc/jwks",
		ScopesSupported:                  []string{"openid"},
		ResponseTypesSupported:           []string{"code"},
		GrantTypesSupported:              []string{"authorization_code"},
		SubjectTypesSupported:            []string{"public"},
		IDTokenSigningAlgValuesSupported: []string{"RS256"},
	}

	return config, nil
}

// GetCASAuthorize CAS3.0客户端授权
func (s *sso) GetCASAuthorize(data *CASAuthorize, userId uint, username string) (callbackUrl, siteName string, err error) {

	// 获取客户端应用
	site, err := dao.Site.GetCASSite(data.Service)
	if err != nil {
		return "", "", errors.New("应用未注册或配置错误")
	}

	// 判断用户是否有权限访问
	if !site.AllOpen {
		if !dao.Site.IsUserInSite(userId, site) {
			return "", site.Name, errors.New("您无权访问该应用")
		}
	}

	// 生成票据（固定格式）
	st := fmt.Sprintf("ST-%d-%s", time.Now().Unix(), username)

	// 使用HMAC SHA-256对票据进行签名
	secret := config.Conf.Settings["secret"].(string)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(st))
	signature := hex.EncodeToString(mac.Sum(nil))

	// 将授权票据写入数据库
	st = fmt.Sprintf("%s-%s", st, signature)
	ticket := &model.SsoCASTicket{
		Ticket:    st,                               // 票据信息
		Service:   site.CallbackUrl,                 // 回调地址
		UserID:    userId,                           // 用户ID
		ExpiresAt: time.Now().Add(10 * time.Second), // 票据的有效期为10秒
	}
	if err = dao.SSO.CreateAuthorizeTicket(ticket); err != nil {
		return "", site.Name, err
	}

	// 返回票据
	redirectURI := fmt.Sprintf("%s?ticket=%s", site.CallbackUrl, st)
	return redirectURI, site.Name, nil
}

// ServiceValidate CAS3.0客户端票据校验
func (s *sso) ServiceValidate(param *CASServiceValidate) (data *CASServiceResponse, err error) {
	// 客户端验证
	_, err = dao.Site.GetCASSite(param.Service)
	if err != nil {
		return nil, errors.New("service string is invalid")
	}

	// 获取票据（如果有数据则表明：1、Code存在，2、在有效期内，3、未使用）
	ticketInfo, err := dao.SSO.GetAuthorizeTicket(param.Ticket)
	if err != nil {
		return nil, errors.New("ticket string is invalid")
	}

	// 分离票据
	parts := strings.Split(param.Ticket, "-")

	// 票据验证：结构验证
	if len(parts) != 4 {
		return nil, errors.New("ticket string is invalid")
	}

	// 获取票据本体
	ticket := fmt.Sprintf("%s-%s-%s", parts[0], parts[1], parts[2])
	// 获取票据签名
	signature := parts[3]

	// 生成新的签名
	secret := config.Conf.Settings["secret"].(string)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(ticket))
	newSignature := hex.EncodeToString(mac.Sum(nil))

	// 票据验证：比较签名
	if !hmac.Equal([]byte(newSignature), []byte(signature)) {
		return nil, errors.New("ticket string is invalid")
	}

	// 获取用户信息
	user, err := dao.User.GetUserInfo(ticketInfo.UserID)
	if err != nil {
		return nil, err
	}

	return &CASServiceResponse{
		Xmlns: "http://www.yale.edu/tp/cas",
		AuthenticationSuccess: &AuthenticationSuccess{
			User: user.Username,
			Attributes: Attributes{
				Id:          uint(user.ID),
				Email:       user.Email,
				Name:        user.Name,
				PhoneNumber: user.PhoneNumber,
				Username:    user.Username,
			},
		},
	}, nil
}

// GetOAuthAuthorize OAuth2.0客户端授权
func (s *sso) GetOAuthAuthorize(data *OAuthAuthorize, userId uint) (callbackUrl, siteName string, err error) {

	// 获取客户端应用
	site, err := dao.Site.GetOAuthSite(data.ClientId)
	if err != nil {
		return "", "", errors.New("应用未注册或配置错误")
	}

	// 判断用户是否有权限访问
	if !site.AllOpen {
		if !dao.Site.IsUserInSite(userId, site) {
			return "", site.Name, errors.New("您无权访问该应用")
		}
	}

	// 创建随机字符串（长度建议>16）
	str := utils.GenerateRandomString(32)
	// 字符串加密，用于返回给客户端授权码
	code, err := utils.Encrypt(str)
	if err != nil {
		return "", site.Name, err
	}

	// 将授权票据写入数据库
	ticket := &model.SsoOAuthTicket{
		Code:        str,                              // 数据库中存放未加密的code，客户端来认证的时候使用的是加密后的code，这样在验证code的时候将前端加密的进行解密判断是否与数据库中的相等即可
		RedirectURI: site.CallbackUrl,                 // 回调地址
		UserID:      userId,                           // 用户ID
		ExpiresAt:   time.Now().Add(10 * time.Second), // 票据的有效期为10秒
		Nonce:       &data.Nonce,
	}
	if err = dao.SSO.CreateAuthorizeCode(ticket); err != nil {
		return "", site.Name, err
	}

	// 返回授权码
	redirectURI := fmt.Sprintf("%s?code=%s&state=%s", site.CallbackUrl, code, data.State)
	return redirectURI, site.Name, nil
}

// GetToken OAuth2.0客户端Token获取
func (s *sso) GetToken(param *Token) (token *ResponseToken, err error) {

	var user *dao.UserInfoWithMenu

	// 客户端验证
	site, err := dao.Site.GetOAuthSite(param.ClientId)
	if err != nil {
		return nil, errors.New("client_id string is invalid")
	}
	if site.ClientSecret != param.ClientSecret {
		return nil, errors.New("client_secret string is invalid")
	}

	// 获取Code（如果有数据则表明：1、Code存在，2、在有效期内，3、未使用）
	code, _ := utils.Decrypt(param.Code)
	ticket, err := dao.SSO.GetAuthorizeCode(code)
	if err != nil {
		return nil, errors.New("code string is invalid")
	}

	// 生成token供access_token和id_token使用（OIDC认证使用的id_token，OAuth认证使用的access_token）
	user, err = dao.User.GetUserInfo(ticket.UserID)
	idToken, err := middleware.GenerateOAuthToken(uint(user.ID), user.Name, user.Username, site.ClientId, "readwrite", *ticket.Nonce)
	if err != nil {
		return nil, err
	}

	token = &ResponseToken{
		IdToken:     idToken,
		AccessToken: idToken,
		TokenType:   "bearer", // 固定值
		ExpiresIn:   3600,     // Token过期时间，这里和配置文件中的JWT过期时间保持一致，也可以独立配置
		Scope:       "openid", // 固定值
	}

	return token, err
}

// GetUserinfo 客户端获取用户信息
func (s *sso) GetUserinfo(token string) (user *ResponseUserinfo, err error) {
	// 验证Token
	mc, err := middleware.ValidateJWT(token)
	if err != nil {
		return nil, err
	}

	// 获取用户信息
	userinfo, err := dao.User.GetUserInfo(mc.ID)

	user = &ResponseUserinfo{
		Id:                uint(userinfo.ID),
		Name:              userinfo.Name,
		Username:          userinfo.Username,
		PreferredUsername: userinfo.Username,
		Email:             userinfo.Email,
		PhoneNumber:       userinfo.PhoneNumber,
	}

	return user, err
}

// GetJwks OIDC客户端获取Jwks
func (s *sso) GetJwks() ([]byte, error) {

	// 读取公钥文件
	pubKey, err := utils.LoadPublicKey()
	if err != nil {
		return nil, err
	}

	// 转换公钥为JWK
	jwkKey, err := jwk.New(pubKey)
	if err != nil {
		return nil, err
	}

	// 将公钥转换为PKIX格式的字节
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return nil, err
	}

	// 基于公钥内容生成kid
	hash := sha256.Sum256(pubKeyBytes)
	kid := base64.URLEncoding.EncodeToString(hash[:])

	// 设置其它参数
	_ = jwkKey.Set(jwk.KeyIDKey, kid)
	_ = jwkKey.Set(jwk.AlgorithmKey, "RS256")
	_ = jwkKey.Set("use", "sig")

	// 创建JWK Set
	jwkSet := jwk.NewSet()
	jwkSet.Add(jwkKey)

	// 将JWK Set序列化为JSON
	jwksJSON, err := json.Marshal(jwkSet)
	if err != nil {
		return nil, err
	}

	return jwksJSON, nil
}

// GetIdPMetadata 获取SAML2 IDP Metadata
func (s *sso) GetIdPMetadata() (metadata string, err error) {

	externalUrl := config.Conf.Settings["externalUrl"].(string)

	// 获取证书
	cert, err := utils.LoadIdpCertificate()
	if err != nil {
		return "", err
	}

	// 创建IDP实例
	idp := saml.IdentityProvider{
		IsIdpInitiated:       false,       // 是否是IdP Initiated模式，true：表示认证请求是通过IdP发起的，false：表示认证请求是客户端（SP）发起的
		Issuer:               externalUrl, // IDP实体，默认为当前服务器地址
		IDPCert:              base64.StdEncoding.EncodeToString(cert.Raw),
		NameIdentifierFormat: saml.AttributeFormatUnspecified,
	}

	// 添加单点登录接口信息
	idp.AddSingleSignOnService(saml.MetadataBinding{
		Binding:  saml.HTTPRedirectBinding, // 由于IDP是前后端分离架构，所以这里使用HTTPRedirectBinding
		Location: externalUrl + "/login",   // 单点登录接口地址
	})

	// 添加IDP组件相关信息
	idp.AddOrganization(saml.Organization{
		OrganizationDisplayName: "IDSphere 统一认证平台", // 组织显示名称
		OrganizationName:        "IDSphere",        // 组织正式名称
		OrganizationURL:         externalUrl,
	})

	// 添加单点登录接口信息（实际不支持单点登出）
	idp.AddSingleSignOutService(saml.MetadataBinding{
		Binding:  saml.HTTPPostBinding,
		Location: externalUrl + "/api/auth/logout",
	})

	// 生成metadata元数据
	metadata, msg := idp.MetaDataResponse()
	if msg != nil {
		return "", msg.Error
	}

	return metadata, nil
}

// ParseSPMetadata SP Metadata解析
func (s *sso) ParseSPMetadata(metadataUrl string) (data *SPMetadata, err error) {

	metadata, err := utils.ParseSPMetadata(metadataUrl)
	if err != nil {
		return nil, err
	}

	// 提取IDP的签名证书
	var signingCertData string
	for _, keyDescriptor := range metadata.SPSSODescriptor.KeyDescriptors {
		if keyDescriptor.Use == "signing" {
			signingCertData = keyDescriptor.KeyInfo.X509Data.X509Certificate
			break
		}
	}
	if signingCertData == "" {
		return nil, errors.New("未找到签名证书")
	}

	return &SPMetadata{
		Certificate: signingCertData,
		EntityID:    metadata.EntityID,
	}, nil
}

// GetSampleAuthnRequest 获取简单SAMLRequest数据
func (s *sso) GetSampleAuthnRequest(samlRequest *SAMLRequest) url.Values {
	payload := url.Values{}
	payload.Add("RelayState", samlRequest.RelayState)
	payload.Add("SAMLRequest", samlRequest.SAMLRequest)
	payload.Add("sigAlg", samlRequest.SigAlg)
	payload.Add("signature", samlRequest.Signature)
	return payload
}

// GetSPAuthorize SP授权
func (s *sso) GetSPAuthorize(samlRequest *SAMLRequest, userId uint) (html, siteName string, err error) {

	var b bytes.Buffer
	externalUrl := config.Conf.Settings["externalUrl"].(string)

	// 获取SAMLRequest数据
	requestData, err := utils.ParseSAMLRequest(samlRequest.SAMLRequest)
	if err != nil {
		return "", "", err
	}

	// 获取SP应用
	site, err := dao.Site.GetSamlSite(requestData.Issuer.Value)
	if err != nil {
		return "", "", errors.New("应用未注册或配置错误")
	}

	// 判断用户是否有权限访问
	if !site.AllOpen {
		if !dao.Site.IsUserInSite(userId, site) {
			return "", site.Name, errors.New("您无权访问该应用")
		}
	}

	// 获取IDP私钥
	privateKeySrt := config.Conf.Settings["privateKey"].(string)

	// 获取IDP证书
	certificate := config.Conf.Settings["certificate"].(string)

	// 获取SP证书（给证书加上头尾）
	SPCert := site.Certificate
	if !strings.HasPrefix(SPCert, "-----BEGIN CERTIFICATE-----") && !strings.HasSuffix(SPCert, "-----END CERTIFICATE-----") {
		SPCert = fmt.Sprintf("-----BEGIN CERTIFICATE-----\n%s\n-----END CERTIFICATE-----\n", site.Certificate)
	}

	// 获取用户信息
	userinfo, err := dao.User.GetUserInfo(userId)
	if err != nil {
		return "", site.Name, err
	}

	// 初始化IDP实（注：也可以在结构体中使用IDPCertFilePath和IDPKeyFilePath从指定路径中读取IDP的证书和私钥，但经测试有Bug）
	idp := saml.IdentityProvider{
		IsIdpInitiated:       false,                                   // 是否为IDP发起认证
		Issuer:               externalUrl,                             // IDP实体
		Audiences:            []string{requestData.Issuer.Value},      // SP实体
		IDPKey:               privateKeySrt,                           // IDP私钥
		IDPCert:              certificate,                             // IDP证书
		SPCert:               SPCert,                                  // SP证书
		NameIdentifier:       userinfo.Username,                       // 用户的唯一标识符
		NameIdentifierFormat: saml.NameIdFormatUnspecified,            // 用户唯一标识符格式
		ACSLocation:          requestData.AssertionConsumerServiceURL, // SP回调地址
		ACSBinging:           saml.HTTPPostBinding,                    // 将SAMLResponse发送到SP的方法
		SessionIndex:         uuid.New().String(),                     // 会话唯一标识符,常用用于会议跟踪
	}

	// 阿里云相关配置（需要给NameID加上域名）
	if strings.HasPrefix(requestData.Issuer.Value, "https://signin.aliyun.com") {
		idp.NameIdentifier = fmt.Sprintf("%s@%s", userinfo.Username, site.DomainId)
	}

	// 添加其它用户属性
	idp.AddAttribute("name", userinfo.Name, saml.AttributeFormatUnspecified)                // 用户姓名
	idp.AddAttribute("username", userinfo.Username, saml.AttributeFormatUnspecified)        // 用户名
	idp.AddAttribute("email", userinfo.Email, saml.AttributeFormatUnspecified)              // 邮箱地址
	idp.AddAttribute("phone_number", userinfo.PhoneNumber, saml.AttributeFormatUnspecified) // 电话号码

	// 华为云专属配置
	idp.AddAttribute("IAM_SAML_Attributes_xUserId", userinfo.Username, saml.AttributeFormatUnspecified)
	idp.AddAttribute("IAM_SAML_Attributes_redirect_url", site.RedirectUrl, saml.AttributeFormatUnspecified) // 登录后跳转的地址
	idp.AddAttribute("IAM_SAML_Attributes_domain_id", site.DomainId, saml.AttributeFormatUnspecified)
	idp.AddAttribute("IAM_SAML_Attributes_idp_id", site.IDPName, saml.AttributeFormatUnspecified)

	// 天翼云专属配置
	idp.AddAttribute("nickName", userinfo.Name, saml.AttributeFormatUnspecified)  // 用户姓名
	idp.AddAttribute("accountId", site.DomainId, saml.AttributeFormatUnspecified) //  天翼云账号ID
	idp.AddAttribute("userId", userinfo.CtyunId, saml.AttributeFormatUnspecified) // 天翼云IAM用户ID
	idp.AddAttribute("idpId", site.DomainId, saml.AttributeFormatUnspecified)     // 天翼云IDP ID

	// 设置认证请求有效期
	idp.AuthnRequestTTL(time.Minute * 10)

	// SAMLRequest请求校验,由于IDP的Metadata元数据中指定的SP Binding方法为HTTPRedirectBinding, 所以这里传入的时候是GET方法
	_, validationError := idp.ValidateAuthnRequest("GET", s.GetSampleAuthnRequest(samlRequest), url.Values{})
	if validationError != nil {
		return "", site.Name, validationError.Error
	}

	// 生成签名后XML数据
	signedXML, signedXMLErr := idp.NewSignedLoginResponse()
	if signedXMLErr != nil {
		return "", site.Name, signedXMLErr.Error
	}

	// 生成HTML响应
	var htmlData = SAMLResponse{
		URL:          idp.ACSLocation,
		SAMLResponse: base64.StdEncoding.EncodeToString([]byte(signedXML)),
		RelayState:   idp.RelayState,
	}
	if err := samlPostFormTemplate.Execute(&b, htmlData); err != nil {
		return "", site.Name, err
	}

	return b.String(), site.Name, nil
}

// Login 单点登录
func (s *sso) Login(queryParams AuthorizeParam, user model.AuthUser) (callbackData, appName string, err error) {

	var (
		data        string
		application string
	)

	if queryParams.GetResponseType() != "" && queryParams.GetClientId() != "" && queryParams.GetRedirectURI() != "" {
		// OAuth认证返回
		params := &OAuthAuthorize{
			ClientId:     queryParams.GetClientId(),
			RedirectURI:  queryParams.GetRedirectURI(),
			ResponseType: queryParams.GetResponseType(),
			Scope:        queryParams.GetScope(),
			State:        queryParams.GetState(),
			Nonce:        queryParams.GetNonce(),
		}
		callbackUrl, siteName, err := s.GetOAuthAuthorize(params, user.ID)
		if err != nil {
			return "", siteName, err
		}
		data = callbackUrl
		application = siteName
	} else if queryParams.GetService() != "" {
		// CAS认证返回
		params := &CASAuthorize{
			Service: queryParams.GetService(),
		}
		callbackUrl, siteName, err := s.GetCASAuthorize(params, user.ID, user.Username)
		if err != nil {
			return "", siteName, err
		}
		data = callbackUrl
		application = siteName
	} else if queryParams.GetSAMLRequest() != "" {
		// SAML2认证返回
		params := &SAMLRequest{
			SAMLRequest: queryParams.GetSAMLRequest(),
			RelayState:  queryParams.GetRelayState(),
			SigAlg:      queryParams.GetSigAlg(),
			Signature:   queryParams.GetSignature(),
		}
		html, siteName, err := s.GetSPAuthorize(params, user.ID)
		if err != nil {
			return "", siteName, err
		}
		data = html
		application = siteName
	}

	return data, application, nil
}
