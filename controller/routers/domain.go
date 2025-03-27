package routers

import (
	"github.com/gin-gonic/gin"
	"ops-api/controller"
)

func initDomainRouters(router *gin.Engine) {
	// 获取域名服务商列表
	router.GET("/api/v1/domain/providers", controller.Domain.GetDomainServiceProviderList)

	// 获取域名列表
	router.GET("/api/v1/domains", controller.Domain.GetDomainList)

	provider := router.Group("/api/v1/domain/provider")
	{
		// 新增域名服务商
		provider.POST("", controller.Domain.AddDomainServiceProvider)

		// 删除域名服务商
		provider.DELETE("/:id", controller.Domain.DeleteDomainServiceProvider)

		// 修改域名服务商
		provider.PUT("", controller.Domain.UpdateDomainServiceProvider)
	}

	domain := router.Group("/api/v1/domain")
	{
		// 新增域名
		domain.POST("", controller.Domain.AddDomain)

		// 删除域名
		domain.DELETE("/:id", controller.Domain.DeleteDomain)

		// 修改域名
		domain.PUT("", controller.Domain.UpdateDomain)

		// 同步域名
		domain.POST("/sync", controller.Domain.SyncDomain)
	}

	dns := router.Group("/api/v1/dns")
	{
		// 获取域名DNS解析列表
		dns.GET("", controller.Domain.GetDomainDnsList)

		// 新增域名DNS解析
		dns.POST("", controller.Domain.AddDomainDns)

		// 修改域名DNS解析
		dns.PUT("", controller.Domain.UpdateDns)

		// 删除域名DNS解析
		dns.DELETE("", controller.Domain.DeleteDns)

		// 修改域名DNS解析状态
		dns.PUT("/status", controller.Domain.SetDomainStatus)
	}
}
