package jaeger

import (
	"context"

	"github.com/google/wire"
	"github.com/goriller/ginny-util/graceful"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-client-go/config"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

// ProviderSet
var ProviderSet = wire.NewSet(NewConfiguration, NewJaegerTracer)

// NewConfiguration
func NewConfiguration(v *viper.Viper) (*config.Configuration, error) {
	var (
		err error
		c   = new(config.Configuration)
	)

	if err = v.UnmarshalKey("jaeger", c); err != nil {
		return nil, errors.Wrap(err, "unmarshal jaeger configuration error")
	}

	return c, nil
}

// NewJaegerTracer NewJaegerTracer for current service
func NewJaegerTracer(cfg *config.Configuration) (opentracing.Tracer, error) {
	// Example logger and metrics factory. Use github.com/uber/jaeger-client-go/log
	// and github.com/uber/jaeger-lib/metrics respectively to bind to real logging and metrics
	// frameworks.
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	// Initialize tracer with a logger and a metrics factory
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory))

	opentracing.SetGlobalTracer(tracer)
	if err != nil {
		return nil, errors.Wrap(err, "create jaeger tracer error")
	}
	graceful.AddCloser(func(ctx context.Context) error {
		return closer.Close()
	})
	return tracer, nil
}
