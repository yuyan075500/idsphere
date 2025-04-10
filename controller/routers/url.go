package routers

import (
	"github.com/gin-gonic/gin"
	"ops-api/controller"
)

func initUrlRouters(router *gin.Engine) {
	// 获取列表
	router.GET("/api/v1/urls", controller.UrlAddress.GetUrlList)

	url := router.Group("/api/v1/url")
	{
		// 新增
		url.POST("", controller.UrlAddress.AddUrl)

		// 删除
		url.DELETE("/:id", controller.UrlAddress.DeleteUrl)

		// 修改
		url.PUT("", controller.UrlAddress.UpdateUrl)

		// 检查
		url.POST("/check", controller.UrlAddress.CertificateCheck)
	}
}
