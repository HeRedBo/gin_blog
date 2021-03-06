package routers

import (
	"gin-blog/middleware/jwt"
	"gin-blog/pkg/export"
	"gin-blog/pkg/qrcode"
	"gin-blog/pkg/setting"
	"gin-blog/pkg/upload"
	"gin-blog/routers/api"
	"gin-blog/routers/api/v1"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter() *gin.Engine {
	r :=  gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	gin.SetMode(setting.RunMode)

	apiGroup := r.Group("/api")
	{
		apiGroup.POST("/auth", api.PostAuth)
	}

	
	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))
	r.POST("/upload", api.UploadImage)
	r.StaticFS("/qrcode", http.Dir(qrcode.GetQrCodeFullPath()))


	apiv1 := r.Group("/api/v1")
	apiv1.Use(jwt.JWT())
	{
		//获取标签列表
		apiv1.GET("/tags", v1.GetTags)
		//新建标签
		apiv1.POST("/tags", v1.AddTag)
		//更新指定标签
		apiv1.PUT("/tags/:id", v1.EditTag)
		//删除指定标签
		apiv1.DELETE("/tags/:id", v1.DeleteTag)
		//导出标签
		apiv1.StaticFS("/export", http.Dir(export.GetExcelFullPath()))
		apiv1.POST("/tags/export", v1.ExportTag)
		//导入标签
		apiv1.POST("/tags/import", v1.ImportTag)
		//获取文章列表
		apiv1.GET("/articles", v1.GetArticles)
		//获取指定文章
		apiv1.GET("/articles/:id", v1.GetArticle)
		//新建文章
		apiv1.POST("/articles", v1.AddArticle)
		// 生成二维码
		apiv1.POST("/articles/poster/generate", v1.GenerateArticlePoster)
		//更新指定文章
		//apiv1.PUT("/articles/:id", v1.EditArticle)
		//删除指定文章
		apiv1.DELETE("/articles/:id", v1.DeleteArticle)
	}

	return r
}
