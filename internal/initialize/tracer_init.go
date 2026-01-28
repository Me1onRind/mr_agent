package initialize

import (
	"context"
	"os"

	"github.com/opentracing/opentracing-go"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmot"
)

func InitOpentracing(serviceName, version string) func(context.Context) error {
	return func(context.Context) error {
		os.Setenv("ELASTIC_APM_SERVER_URL", "http://localhost:8200")
		os.Setenv("ELASTIC_APM_SECRET_TOKEN", "")
		os.Setenv("ELASTIC_APM_STACK_TRACE_LIMIT", "0")
		os.Setenv("ELASTIC_APM_USE_ELASTIC_TRACEPARENT_HEADER", "false")
		os.Setenv("ELASTIC_APM_TRANSACTION_SAMPLE_RATE", "1.0")
		tracer, err := apm.NewTracer(serviceName, version)
		if err != nil {
			return err
		}
		opentracing.SetGlobalTracer(apmot.New(apmot.WithTracer(tracer)))
		return nil
	}
}
