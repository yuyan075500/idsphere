package routers

import (
	"github.com/gin-gonic/gin"
	"ops-api/controller"
)

func initDomainCertificateRouters(router *gin.Engine) {
	// 获取域名证书列表
	router.GET("/api/v1/certificates", controller.Certificate.GetDomainCertificateList)

	certificate := router.Group("/api/v1/certificate")
	{
		// 上传域名证书
		certificate.POST("/upload", controller.Certificate.UploadDomainCertificate)

		// 删除域名证书
		certificate.DELETE("/:id", controller.Certificate.DeleteDomainCertificate)

		// 证书下载
		certificate.GET("/:id", controller.Certificate.DownloadDomainCertificate)
	}
}
