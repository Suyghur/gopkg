//@File     metrics.go
//@Time     2024/8/6
//@Author   #Suyghur,

package gormx

import gozerometric "github.com/zeromicro/go-zero/core/metric"

const namespace = "gormx_client"

var (
	metricReqDur = gozerometric.NewHistogramVec(&gozerometric.HistogramVecOpts{
		Namespace: namespace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "gormx client requests duration(ms).",
		Labels:    []string{"command"},
		Buckets:   []float64{0.25, 0.5, 1, 1.5, 2, 3, 5, 10, 25, 50, 100, 250, 500, 1000, 2000, 5000, 10000, 15000},
	})
	metricReqErr = gozerometric.NewCounterVec(&gozerometric.CounterVecOpts{
		Namespace: namespace,
		Subsystem: "requests",
		Name:      "error_total",
		Help:      "gormx client requests error count.",
		Labels:    []string{"command"},
	})
	metricSlowCount = gozerometric.NewCounterVec(&gozerometric.CounterVecOpts{
		Namespace: namespace,
		Subsystem: "requests",
		Name:      "slow_total",
		Help:      "gormx client requests slow count.",
		Labels:    []string{"command"},
	})
)
