//@File     option.go
//@Time     2024/7/18
//@Author   #Suyghur,

package gormx

import "time"

type (
	Option func(*config)
)

func WithMaxIdleConns(num int) Option {
	return func(c *config) {
		c.MaxIdleConns = num
	}
}

func WithMaxOpenConns(num int) Option {
	return func(c *config) {
		c.MaxOpenConns = num
	}
}

func WithConnMaxLifetime(duration time.Duration) Option {
	return func(c *config) {
		c.ConnMaxLifetime = duration
	}
}

func WithSlowThreshold(threshold time.Duration) Option {
	return func(c *config) {
		c.SlowThreshold = threshold
	}
}

func WithDialTimeout(timeout time.Duration) Option {
	return func(c *config) {
		c.DialTimeout = timeout
	}
}
