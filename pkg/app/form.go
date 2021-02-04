package app

import (
	"gin-blog/pkg/e"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)


// 参数绑定用户数据校验
func BindAndValid(c *gin.Context, form interface{}) (int, int) {
	err := c.Bind(form)
	if err != nil {
		return http.StatusBadRequest, e.INVALID_PARAMS
	}
	valid := validation.Validation{}
	check, err := valid.Valid(form)

	if err != nil {
		return http.StatusInternalServerError, e.ERROR
	}

	if !check {
		MakeError(valid.Errors)
		return http.StatusBadRequest,e.INVALID_PARAMS
	}
	return http.StatusOK, e.SUCCESS
}
