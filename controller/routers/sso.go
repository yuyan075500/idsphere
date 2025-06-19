package routers

import (
	"github.com/gin-gonic/gin"
	"ops-api/controller"
)

// 单点登录相关路由
func initSSORouters(router *gin.Engine) {

	sso := router.Group("/api/v1/sso")
	{
		// 获取授权（Nginx）
		sso.GET("/nginx/authorize", controller.SSO.NginxAuthorize)
		// 获取授权（OAuth2.0）
		sso.POST("/oauth/authorize", controller.SSO.OAuthAuthorize)
		// 获取Token（OAuth2.0）
		sso.POST("/oauth/token", controller.SSO.GetToken)
		// 获取用户信息（OAuth2.0 GET请求）
		sso.GET("/oauth/userinfo", controller.SSO.GetUserInfo)
		// 获取用户信息（OAuth2.0 POST请求）
		sso.POST("/oauth/userinfo", controller.SSO.GetUserInfo)
		// 获取Jwks配置
		sso.GET("/oidc/jwks", controller.SSO.GetJwksConfig)
		// 获取授权（CAS3.0）
		sso.POST("/cas/authorize", controller.SSO.CASAuthorize)
		// 获取IDP元数据（SAML2）
		sso.GET("/saml/metadata", controller.SSO.GetIdPMetadata)
		// SP授权（SAML2）
		sso.POST("/saml/authorize", controller.SSO.SPAuthorize)
		// SP元数据解析
		sso.POST("/saml/metadata", controller.Site.ParseSPMetadata)
		// Cookie认证（用于Nginx转发过来的认证请求）
		sso.GET("/cookie/auth", controller.SSO.CookieAuth)
	}

	// CAS3.0客户端票据校验
	router.GET("/p3/serviceValidate", controller.SSO.CASServiceValidate)
	// 获取OIDC配置
	router.GET("/.well-known/openid-configuration", controller.SSO.GetOIDCConfig)
}
