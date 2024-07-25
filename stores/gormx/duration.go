//@File     duration.go
//@Time     2024/7/18
//@Author   #Suyghur,

package gormx

import (
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

func durationHook() Hook {
	return func(name string, next Handler) Handler {
		return func(db *gorm.DB) {
			logx.Infof("duration hook start...")
			next(db)
			logx.Infof("duration hook end...")
		}
	}
}
