package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"gorm.io/gorm"
	"io"
	"net/http"
	"ops-api/model"
	"ops-api/utils"
)

// ExcludedPaths 不记录日志的接口
var ExcludedPaths = map[string]bool{
	"/api/auth/login":                   true,
	"/api/v1/sso/saml/post":             true,
	"/api/auth/logout":                  true,
	"/api/auth/ww_login":                true,
	"/api/auth/dingtalk_login":          true,
	"/api/auth/feishu_login":            true,
	"/api/v1/user/avatarUpload":         true,
	"/api/v1/user/sync/ad":              true,
	"/api/v1/reset_password":            true,
	"/api/v1/site/logoUpload":           true,
	"/api/v1/sms/huawei/callback":       true,
	"/api/v1/sms/reset_password":        true,
	"/api/v1/user/mfa_qrcode":           true,
	"/api/v1/user/mfa_auth":             true,
	"/api/v1/sso/oauth/authorize":       true,
	"/api/v1/sso/oauth/token":           true,
	"/api/v1/sso/oauth/userinfo":        true,
	"/api/v1/sso/cas/authorize":         true,
	"/api/v1/sso/saml/authorize":        true,
	"/api/v1/sso/saml/metadata":         true,
	"/api/v1/account/code_verification": true,
}

// Oplog 操作日志
func Oplog(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// 只记录POST、PUT和DELETE请求
		if c.Request.Method != http.MethodPost && c.Request.Method != http.MethodPut && c.Request.Method != http.MethodDelete {
			c.Next()
			return
		}

		// 跳过不需要记录的接口
		if ExcludedPaths[path] {
			c.Next()
			return
		}

		// 获取请求参数
		var reqBodyBytes []byte
		if c.Request.Body != nil {
			reqBodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBodyBytes))
		}

		// 请求参数敏感信息过滤
		var requestData map[string]interface{}
		if err := json.Unmarshal(reqBodyBytes, &requestData); err == nil {
			utils.FilterFields(requestData)
		}

		// 将请求参数转换为JSON字符串
		requestDataStr, err := json.Marshal(requestData)
		if err != nil {
			requestDataStr = []byte("{}")
		}

		// 获取客户端信息
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()
		// 获取当前登录用户的用户名
		username, _ := c.Get("username")

		// 执行请求
		c.Next()

		// 记录响应数据
		var responseDataStr string
		responseData, exists := c.Get("response")
		if !exists {
			responseData = ""
		}
		if data, ok := responseData.(map[string]interface{}); ok {
			// 将map转换为JSON字符串
			jsonData, err := json.Marshal(data)
			if err != nil {
				// 转换失败
				responseDataStr = "{}"
			} else {
				responseDataStr = string(jsonData)
			}
		} else if data, ok := responseData.(string); ok {
			// 如果已经是字符串类型
			responseDataStr = data
		} else {
			// 如果没有数据
			responseDataStr = ""
		}

		// 记录操作日志
		log := model.LogOplog{
			Username:      username.(string),
			Endpoint:      path,
			Method:        c.Request.Method,
			RequestParams: string(requestDataStr),
			ResponseData:  responseDataStr,
			ClientIP:      clientIP,
			UserAgent:     userAgent,
		}

		// 将日志保存到数据库
		if err := db.Create(&log).Error; err != nil {
			logger.Warn("保存操作日志失败: %v", err)
		}
	}
}
