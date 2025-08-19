//@File     controller_context.go
//@Time     2024/8/10
//@Author   #Suyghur,

package v2

import (
	"github.com/robfig/cron/v3"
	"github.com/yyxxgame/gopkg/infrastructure/cron/v2/internal"
	"github.com/zeromicro/go-zero/core/logx"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type (
	ICronTaskController interface {
		Start()
		Stop()
		RegisterJobs(jobs ...ICronJob)
	}

	controller struct {
		conf     CronTaskConf
		cron     *cron.Cron
		tracer   oteltrace.Tracer
		handlers map[string]*WrapperJob
		hooks    []Hook
	}
)

func NewCronTaskController(conf CronTaskConf, opts ...Option) ICronTaskController {
	c := &controller{
		conf: conf,
		cron: cron.New(
			cron.WithSeconds(),
			cron.WithChain(
				cron.SkipIfStillRunning(internal.DefaultLogger),
				cron.Recover(internal.DefaultLogger),
			),
		),
		handlers: make(map[string]*WrapperJob),
		hooks:    []Hook{},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *controller) Start() {
	c.cron.Start()
}

func (c *controller) Stop() {
	c.cron.Stop()
}

func (c *controller) RegisterJobs(jobs ...ICronJob) {
	for _, job := range jobs {
		for _, item := range c.conf.Jobs {
			if !item.Enable {
				continue
			}

			if item.Name != job.Named() {
				continue
			}

			if _, exists := c.handlers[item.Name]; exists {
				logx.Errorf("[CRON-TASK-CONTROLLER-ERROR]: register job: %s on duplicated error, skip it ...", item.Name)
				continue
			}

			wrapperHooks := []Hook{
				NewDurationHook(c.tracer).ExecHook(),
			}
			wrapperHooks = append(wrapperHooks, c.hooks...)

			wrapper := &WrapperJob{
				cronJob:   job,
				params:    item.Params,
				finalHook: chainHooks(wrapperHooks...),
			}

			_, err := c.cron.AddJob(item.Expression, wrapper)
			if err != nil {
				logx.Errorf("[CRON-TASK-CONTROLLER-ERROR]: register job: %s on error: %v", item.Name, err)
			} else {
				logx.Infof("[CRON-TASK-CONTROLLER]: register job: %s on success", item.Name)
			}

			c.handlers[item.Name] = wrapper
		}
	}
}
