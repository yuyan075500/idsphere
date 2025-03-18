package routers

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"ops-api/config"
	"ops-api/docs"
)

var Router router

type router struct{}

func (r *router) InitRouter(router *gin.Engine) {

	swagger := config.Conf.Settings["swagger"].(bool)

	// Swagger接口文档
	if swagger {
		docs.SwaggerInfo.BasePath = ""
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	// 注册 pprof 路由
	pprof.Register(router)

	// 初始化不同类型路由
	initUserRouters(router)
	initGroupRouters(router)
	initSiteRouters(router)
	initAuditRouters(router)
	initSmsRouters(router)
	initMenuRouters(router)
	initAuthRouters(router)
	initSSORouters(router)
	initTagRouters(router)
	initTaskRouters(router)
	initAccountRouters(router)
	initSettingsRouters(router)
	initDomainRouters(router)
	initDomainCertificateRouters(router)

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.String(200, "ok")
	})
}
