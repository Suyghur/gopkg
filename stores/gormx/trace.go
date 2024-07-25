//@File     trace.go
//@Time     2024/7/18
//@Author   #Suyghur,

package gormx

import (
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

func traceHook() Hook {
	return func(name string, next Handler) Handler {
		return func(db *gorm.DB) {
			logx.Infof("trace hook start...")
			next(db)
			logx.Infof("trace hook end...")
		}
	}
}
