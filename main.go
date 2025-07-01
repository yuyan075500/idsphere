package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"ops-api/config"
	"ops-api/controller/routers"
	"ops-api/db"
	"ops-api/global"
	"ops-api/middleware"
	"ops-api/service"
)

func main() {

	// 配置文件初始化
	config.InitConfig()

	// 初始化Minio
	if err := db.MinioInit(); err != nil {
		logger.Error("Minio初始化失败：", err.Error())
		return
	}

	// 初始化MySQL
	if err := db.MySQLInit(); err != nil {
		logger.Error("MySQL初始化失败：", err.Error())
		return
	}

	// 初始Redis
	if err := db.RedisInit(); err != nil {
		logger.Error("Redis初始化失败：", err.Error())
		return
	}

	// 初始化CasBin权限
	if err := middleware.CasBinInit(); err != nil {
		logger.Error("权限初始化失败：", err.Error())
		return
	}

	// 初始化定时任务
	if err := service.TaskInit(); err != nil {
		logger.Error("定时任务初始化失败：", err.Error())
		return
	}

	r := gin.Default()

	// 加载跨域中间件
	r.Use(middleware.Cors())
	// 加载登录中间件，其中IgnorePaths()方法可以忽略不需要登录认证的路由，支持前缀匹配
	r.Use(middleware.LoginBuilder().
		IgnorePaths("/api/auth/login").
		IgnorePaths("/health").
		IgnorePaths("/swagger/").
		IgnorePaths("/debug/pprof/").
		IgnorePaths("/api/v1/sms/huawei/callback").
		IgnorePaths("/api/v1/sms/reset_password").
		IgnorePaths("/api/v1/reset_password").
		IgnorePaths("/api/v1/user/mfa_qrcode").
		IgnorePaths("/api/v1/user/mfa_auth").
		IgnorePaths("/api/v1/sso/oauth/token").
		IgnorePaths("/api/v1/sso/oauth/userinfo").
		IgnorePaths("/p3/serviceValidate").
		IgnorePaths("/api/v1/sso/saml/metadata").
		IgnorePaths("/api/v1/sso/saml/post").
		IgnorePaths("/api/v1/sso/saml/authorize").
		IgnorePaths("/.well-known/openid-configuration").
		IgnorePaths("/api/v1/sso/oidc/jwks").
		IgnorePaths("/api/v1/sso/cookie/auth").
		IgnorePaths("/api/auth/dingtalk_login").
		IgnorePaths("/api/auth/ww_login").
		IgnorePaths("/api/auth/feishu_login").
		IgnorePaths("/api/v1/site/guide").
		Build())
	// 加载权限中间件
	r.Use(middleware.PermissionCheck())
	// 加载记录日志中间件
	r.Use(middleware.Oplog(global.MySQLClient))

	// 加载模板文档
	r.LoadHTMLGlob("templates/*")

	// 注册路由
	routers.Router.InitRouter(r)

	// 启动服务
	if err := r.Run(fmt.Sprintf("%v", config.Conf.Server)); err != nil {
		logger.Error("IDSphere服务启动失败：", err.Error())
	}
}
