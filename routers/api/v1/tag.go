package v1

import (
	"gin-blog/models"
	"gin-blog/pkg/app"
	"gin-blog/pkg/e"
	"gin-blog/pkg/export"
	"gin-blog/pkg/logging"
	"gin-blog/pkg/setting"
	"gin-blog/pkg/util"
	"gin-blog/service/tag_service"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

// 获取文章多个标签
func GetTags(c *gin.Context) {
	name := c.Query("name")

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	if name != "" {
		maps["mame"] = name
	}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
	}

	code := e.SUCCESS

	lists, err := models.GetTags(util.GetPage(c), setting.PageSize, maps)
	if err != nil {
		code = e.ERROR_EXIST_TAG
	}
	data["lists"] = lists
	data["total"] = models.GetTagTotal(maps)

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

// 新增文章标签
func AddTag(c *gin.Context) {
	name := c.PostForm("name")
	state := com.StrTo(c.DefaultPostForm("state", "0")).MustInt()
	createdBy := c.PostForm("created_by")

	valid := validation.Validation{}
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.MaxSize(createdBy, 100, "created_by").Message("创建人最长为100字符")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	code := e.INVALID_PARAMS

	if !valid.HasErrors() {
		check, err := models.ExistTagByName(name)
		if err != nil {
			code = e.ERROR_EXIST_TAG
		}

		if check {
			code = e.SUCCESS
			models.AddTag(name, state, createdBy)
		}

		//if !models.ExistTagByName(name) {
		//	code = e.SUCCESS
		//	models.AddTag(name, state, createdBy)
		//} else {
		//	code = e.ERROR_EXIST_TAG
		//}

	} else {
		for _, err := range valid.Errors {
			logging.Error(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}

// 修改文章标签
func EditTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	name := c.PostForm("name")
	modifiedBy := c.PostForm("modified_by")

	valid := validation.Validation{}
	var state int = 1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	valid.Required(id, "id").Message("ID不能为空")
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")
	valid.Required(modifiedBy, "created_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "created_by").Message("修改人最长为100字符")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		code = e.SUCCESS
		//exists , err := articleService.ExistByID()
		//if err != nil {
		//	appG.Response(http.StatusOK,e.ERROR_CHECK_EXIST_ARTICLE_FAIL,nil)
		//	return
		//}
		exists , err := models.ExistTagByID(id)
		if err != nil {
			code = e.ERROR_NOT_EXIST_TAG
		}
		if exists {
			data := make(map[string]interface{})
			data["modified_by"] = modifiedBy
			if name != "" {
				data["name"] = name
			}
			if state != -1 {
				data["state"] = state
			}

			models.EditTag(id, data)
		}
		//if models.ExistTagByID(id) {
		//	data := make(map[string]interface{})
		//	data["modified_by"] = modifiedBy
		//	if name != "" {
		//		data["name"] = name
		//	}
		//	if state != -1 {
		//		data["state"] = state
		//	}
		//
		//	models.EditTag(id, data)
		//} else {
		//	code = e.ERROR_NOT_EXIST_TAG
		//}
	} else {
		for _, err := range valid.Errors {
			logging.Error(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})

}

// 删除文章标签
func DeleteTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		code = e.SUCCESS

		exists , err := models.ExistTagByID(id)
		if err != nil {
			code = e.ERROR_NOT_EXIST_TAG
		}
		if exists {
			models.DeleteTag(id)
		}
		//if models.ExistTagByID(id) {
		//	models.DeleteTag(id)
		//} else {
		//	code = e.ERROR_NOT_EXIST_TAG
		//}
	} else {
		for _, err := range valid.Errors {
			logging.Error(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}

func ExportTag(c *gin.Context) {
	appG := app.Gin{C: c}

	name := c.PostForm("name")
	state := 1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	tagService := tag_service.Tag{
		Name:  name,
		State: state,
	}

	filename, err := tagService.Export()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_EXPORT_TAG_FAIL, nil)
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"export_url":      export.GetExcelFullUrl(filename),
		"export_save_url": export.GetExcelPath() + filename,
	})
}

func ImportTag(c *gin.Context) {
	appG := app.Gin{C: c}

	file, _,err := c.Request.FormFile("file")
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}

	tagService := tag_service.Tag{}
	err = tagService.Import(file)
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
