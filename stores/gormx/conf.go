//@File     config.go
//@Time     2024/7/18
//@Author   #Suyghur,

package gormx

import "time"

const (
	defaultMaxIdleConns    = 10
	defaultMaxOpenConns    = 100
	defaultConnMaxLifetime = 300 * time.Second
	defaultSlowThreshold   = 500 * time.Millisecond
	defaultDialTimeout     = 3 * time.Second
)

type config struct {
	// 最大空闲连接数
	MaxIdleConns int
	// 最大活动连接数
	MaxOpenConns int
	// 连接最大存活时间
	ConnMaxLifetime time.Duration
	// 慢日志阈值
	SlowThreshold time.Duration
	// 拨号超时时间
	DialTimeout time.Duration
}
