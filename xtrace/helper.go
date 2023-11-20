//@File     : helper.go
//@Time     : 2023/2/21
//@Auther   : Kaishin

package xtrace

import (
	"bytes"
	"context"
	"fmt"
	gozerotrace "github.com/zeromicro/go-zero/core/trace"
	"github.com/zeromicro/go-zero/rest"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"io"
	"net/http"
)

func GetTraceId(ctx context.Context) trace.TraceID {
	/*
		get trace id
	*/
	span := trace.SpanFromContext(ctx)
	return span.SpanContext().TraceID()
}

func GetCarrier(ctx context.Context) *propagation.HeaderCarrier {
	carrier := &propagation.HeaderCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	return carrier
}

func AddTags(ctx context.Context, kv ...attribute.KeyValue) {
	/*
		add tags by ctx
	*/
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(kv...)
}

func AddTagsByMap(ctx context.Context, kvs map[string]interface{}) {
	for k, v := range kvs {
		vStr := fmt.Sprintf("%v", v)
		AddTags(ctx, attribute.String(k, vStr))
	}
}

func AddEvent(ctx context.Context, name string, kv ...attribute.KeyValue) {
	/*
		add evnet by ctx
	*/
	span := trace.SpanFromContext(ctx)
	span.AddEvent(name, trace.WithAttributes(kv...))
}

func TracingOnApiSvr(server *rest.Server) {
	/*
		running go-zero api server
	*/
	server.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			bytesData, err := io.ReadAll(r.Body)
			if err == nil {
				r.Body = io.NopCloser(bytes.NewBuffer(bytesData))
				jsonString := string(bytesData)

				// todo body parsed into tags
				//var data map[string]interface{}
				//err = json.Unmarshal([]byte(jsonString), &data)
				//if err != nil {
				//	next(w, r)
				//}
				//AddTagsByMap(r.Context(), data)

				// write body info into events
				AddEvent(r.Context(), r.RequestURI, attribute.String("params", jsonString))
			}
			next(w, r)
		}
	})
}

func MakeHeaderContext(name string) context.Context {
	ctx := context.Background()
	tracer := otel.GetTracerProvider().Tracer(gozerotrace.TraceName)
	_, span := tracer.Start(ctx, name, trace.WithSpanKind(trace.SpanKindInternal))
	return trace.ContextWithSpan(ctx, span)
}
