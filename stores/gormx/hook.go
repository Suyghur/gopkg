//@File     hook.go
//@Time     2024/7/18
//@Author   #Suyghur,

package gormx

import (
	"gorm.io/gorm"
)

type (

	//ProducerHookFunc func(ctx context.Context, message *sarama.ProducerMessage) error

	//ProducerHook func(ctx context.Context, message *sarama.ProducerMessage, next ProducerHookFunc) error

	Handler func(*gorm.DB)

	Hook func(name string, next Handler) Handler

	processor interface {
		Get(name string) func(*gorm.DB)
		Replace(name string, handler func(*gorm.DB)) error
	}
)

func registerHook(db *gorm.DB, hooks ...Hook) {
	var processors = []struct {
		Name      string
		Processor processor
	}{
		{"gorm:create", db.Callback().Create()},
		{"gorm:query", db.Callback().Query()},
		{"gorm:delete", db.Callback().Delete()},
		{"gorm:update", db.Callback().Update()},
		{"gorm:row", db.Callback().Row()},
		{"gorm:raw", db.Callback().Raw()},
	}

	for _, hook := range hooks {
		for _, p := range processors {
			handler := p.Processor.Get(p.Name)
			handler = hook(p.Name, handler)

			if err := p.Processor.Replace(p.Name, handler); err != nil {
				panic(err)
			}
		}
	}
}
