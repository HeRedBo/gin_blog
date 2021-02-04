package app

import (
	"gin-blog/pkg/logging"
	"github.com/astaxie/beego/validation"
)

func MakeError(errors []*validation.Error) {
	for _, err := range errors {
		logging.Error(err.Key, err.Message)
	}
	return
}
