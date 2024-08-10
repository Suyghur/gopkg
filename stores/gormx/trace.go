//@File     trace.go
//@Time     2024/7/18
//@Author   #Suyghur,

package gormx

import (
	"go.opentelemetry.io/otel/attribute"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

var sqlAttributeKey = attribute.Key("sql.method")

const spanName = "sql"

func traceHook() Hook {
	return func(name, command string, next Handler) Handler {
		return func(db *gorm.DB) {
			logx.Infof("trace hook start...")
			next(db)
			if db.Error != nil {
			}
			logx.Infof("trace hook end...")
		}
	}
}

//func startSpan(ctx context.Context, method string) (context.Context, oteltrace.Span) {
//	tracer := trace.TracerFromContext(ctx)
//	start, span := tracer.Start(ctx, spanName, oteltrace.WithSpanKind(oteltrace.SpanKindClient))
//	span.SetAttributes(sqlAttributeKey.String(method))
//
//	return start, span
//}
//
//func endSpan(span oteltrace.Span, err error) {
//	defer span.End()
//
//	e := logger.ErrRecordNotFound
//
//	if err == nil || errors.Is(err, sql.ErrNoRows) {
//		span.SetStatus(codes.Ok, "")
//		return
//	}
//
//	span.SetStatus(codes.Error, err.Error())
//	span.RecordError(err)
//}
