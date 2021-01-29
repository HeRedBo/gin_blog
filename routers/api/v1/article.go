package v1

import (
	"gin-blog/models"
	"gin-blog/pkg/app"
	"gin-blog/pkg/e"
	"gin-blog/pkg/logging"
	"gin-blog/pkg/qrcode"
	"gin-blog/pkg/setting"
	"gin-blog/pkg/util"
	"gin-blog/service/article_service"
	"github.com/astaxie/beego/validation"
	"github.com/boombuler/barcode/qr"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"log"
	"net/http"
)

const (
	QRCODE_URL = "https://github.com/EDDYCJY/blog#gin%E7%B3%BB%E5%88%97%E7%9B%AE%E5%BD%95"
)

// 获取单个文章
func GetArticle(c *gin.Context) {
	appG := app.Gin{c}
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		app.MakeError(valid.Errors)
		appG.Response(http.StatusOK,e.INVALID_PARAMS, nil)
	}

	articleService := article_service.Article{ID: id}
	exists , err := articleService.ExistByID()
	if err != nil {
		appG.Response(http.StatusOK,e.ERROR_CHECK_EXIST_ARTICLE_FAIL,nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK,e.ERROR_CHECK_EXIST_ARTICLE_FAIL,nil)
		return
	}

	article, err := articleService.Get()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_ARTICLE_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, article)

}

func GetArticles(c *gin.Context) {
	data := make(map[string]interface{})
	maps := make(map[string]interface{})
	valid := validation.Validation{}

	var state int = -1

	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	var tagId int = -1
	if arg := c.Query("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		maps["tag_id"] = tagId

		valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
	}

	code := e.INVALID_PARAMS
	if ! valid.HasErrors() {
		code = e.SUCCESS

		data["lists"] = models.GetArticles(util.GetPage(c), setting.PageSize, maps)
		data["total"] = models.GetArticleTotal(maps)

	} else {
		for _, err := range valid.Errors {
			log.Printf("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code" : code,
		"msg" : e.GetMsg(code),
		"data" : data,
	})

}

/**
新增文章
 */
func AddArticle(c *gin.Context) {

	tagId := com.StrTo(c.PostForm("tag_id")).MustInt()
	title := c.PostForm("title")
	desc  := c.PostForm("desc")
	content := c.PostForm("content")
	createdBy := c.PostForm("created_by")
	state := com.StrTo(c.DefaultPostForm("state", "0")).MustInt()

	valid := validation.Validation{}
	valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
	valid.Required(title, "title").Message("标题不能为空")
	valid.Required(desc, "desc").Message("简述不能为空")
	valid.Required(content, "content").Message("内容不能为空")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	code := e.INVALID_PARAMS
	if ! valid.HasErrors() {
		if models.ExistTagByID(tagId) {
			data := make(map[string]interface {})
			data["tag_id"] = tagId
			data["title"] = title
			data["desc"] = desc
			data["content"] = content
			data["created_by"] = createdBy
			data["state"] = state

			models.AddArticle(data)
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	} else {
		for _, err := range valid.Errors {
			log.Printf("err.key: %s, err.message: %s", err.Key, err.Message)
			logging.Error(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code" : code,
		"msg" : e.GetMsg(code),
		"data" : make(map[string]interface{}),
	})
}


//修改文章
//func EditArticle(c *gin.Context) {
//	valid := validation.Validation{}
//
//	id := com.StrTo(c.Param("id")).MustInt()
//	tagId := com.StrTo(c.PostForm("tag_id")).MustInt()
//	title := c.PostForm("title")
//	desc := c.PostForm("desc")
//	content := c.PostForm("content")
//	modifiedBy := c.PostForm("modified_by")
//
//	var state int = -1
//	if arg := c.PostForm("state"); arg != "" {
//		state = com.StrTo(arg).MustInt()
//		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
//	}
//
//	valid.Min(id, 1, "id").Message("ID必须大于0")
//	valid.MaxSize(title, 100, "title").Message("标题最长为100字符")
//	valid.MaxSize(desc, 255, "desc").Message("简述最长为255字符")
//	valid.MaxSize(content, 65535, "content").Message("内容最长为65535字符")
//	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
//	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")
//
//	code := e.INVALID_PARAMS
//	if ! valid.HasErrors() {
//		if models.ExistArticleByID(id) {
//			if models.ExistTagByID(tagId) {
//				data := make(map[string]interface {})
//				if tagId > 0 {
//					data["tag_id"] = tagId
//				}
//				if title != "" {
//					data["title"] = title
//				}
//				if desc != "" {
//					data["desc"] = desc
//				}
//				if content != "" {
//					data["content"] = content
//				}
//
//				data["modified_by"] = modifiedBy
//
//				models.EditArticle(id, data)
//				code = e.SUCCESS
//			} else {
//				code = e.ERROR_NOT_EXIST_TAG
//			}
//		} else {
//			code = e.ERROR_NOT_EXIST_ARTICLE
//		}
//	} else {
//		for _, err := range valid.Errors {
//			log.Printf("err.key: %s, err.message: %s", err.Key, err.Message)
//			logging.Error(err.Key, err.Message)
//		}
//	}
//
//	c.JSON(http.StatusOK, gin.H{
//		"code" : code,
//		"msg" : e.GetMsg(code),
//		"data" : make(map[string]string),
//	})
//}

//删除文章
func DeleteArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := e.INVALID_PARAMS
	//if ! valid.HasErrors() {
	//	if models.ExistArticleByID(id) {
	//		models.DeleteArticle(id)
	//		code = e.SUCCESS
	//	} else {
	//		code = e.ERROR_NOT_EXIST_ARTICLE
	//	}
	//} else {
	//	for _, err := range valid.Errors {
	//		log.Printf("err.key: %s, err.message: %s", err.Key, err.Message)
	//		logging.Error(err.Key, err.Message)
	//	}
	//}

	c.JSON(http.StatusOK, gin.H{
		"code" : code,
		"msg" : e.GetMsg(code),
		"data" : make(map[string]string),
	})
}

func GenerateArticlePoster(c *gin.Context)  {
	appG := app.Gin{c}
	article := &article_service.Article{}
	qr := qrcode.NewQrcode(QRCODE_URL,300, 300, qr.M, qr.Auto) // 目前写死 gin 系列路径，可自行增加业务逻辑
	posterName := article_service.GetPosterFlag() + "-" + qrcode.GetQrCodeFileName(qr.URL) +  qr.GetQrCodeExt()
	articlePoster := article_service.NewArticlePoster(posterName, article, qr)
	articlePosterBgService := article_service.NewArticlePosterBg(
		"gb.jpg",
		articlePoster,
		&article_service.Rect{
			X0:   0,
			Y0:   0,
			X1:   550,
			Y1:   700,
		},
		&article_service.Pt{
			X: 125,
			Y: 298,
		},
		)

	_, filePath , err := articlePosterBgService.Generate()
	if err != nil {
		logging.Error(err)
		appG.Response(http.StatusOK, e.ERROR_GEN_ARTICLE_POSTER_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK,e.SUCCESS, map[string]string{
		"poster_url" : qrcode.GetQrCodeFullUrl(posterName),
		"post_save_url": filePath + posterName,
	})
}

