//@File     hook.go
//@Time     2024/7/18
//@Author   #Suyghur,

package gormx

import (
	"gorm.io/gorm"
)

type (
	Handler func(*gorm.DB)

	Hook func(name, command string, next Handler) Handler

	processor interface {
		Get(name string) func(*gorm.DB)
		Replace(name string, handler func(*gorm.DB)) error
	}
)

func registerHook(db *gorm.DB, hooks ...Hook) {
	var processors = []struct {
		Name      string
		Command   string
		Processor processor
	}{
		{"gorm:create", "create", db.Callback().Create()},
		{"gorm:query", "query", db.Callback().Query()},
		{"gorm:delete", "delete", db.Callback().Delete()},
		{"gorm:update", "update", db.Callback().Update()},
		{"gorm:row", "row", db.Callback().Row()},
		{"gorm:raw", "raw", db.Callback().Raw()},
	}

	for _, hook := range hooks {
		for _, p := range processors {
			handler := p.Processor.Get(p.Name)
			handler = hook(p.Name, p.Command, handler)

			if err := p.Processor.Replace(p.Name, handler); err != nil {
				panic(err)
			}
		}
	}
}
