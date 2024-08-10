//@File     duration.go
//@Time     2024/7/18
//@Author   #Suyghur,

package gormx

import (
	"time"

	"github.com/zeromicro/go-zero/core/timex"
	"gorm.io/gorm"
)

func durationHook(slowThreshold time.Duration) Hook {
	return func(name, command string, next Handler) Handler {
		return func(db *gorm.DB) {
			start := timex.Now()
			next(db)
			time.Sleep(time.Second)
			duration := timex.Since(start)

			metricReqDur.Observe(duration.Milliseconds(), command)

			if db.Error != nil {
				metricReqErr.Inc(command)
			}

			if slowThreshold > 0 && duration > slowThreshold {
				metricSlowCount.Inc(command)
			}
		}
	}
}
