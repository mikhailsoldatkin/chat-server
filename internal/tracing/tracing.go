package tracing

import (
	"github.com/mikhailsoldatkin/platform_common/pkg/closer"
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
)

// Init initializes the Jaeger tracer for distributed tracing.
func Init(logger *zap.Logger, serviceName string, jaegerAddress string) {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LocalAgentHostPort: jaegerAddress,
		},
	}

	c, err := cfg.InitGlobalTracer(serviceName)

	closerFunc := func() error {
		errClose := c.Close()
		if errClose != nil {
			return errClose
		}
		return nil
	}
	closer.Add(closerFunc)

	if err != nil {
		logger.Fatal("failed to init tracing", zap.Error(err))
	}
}
