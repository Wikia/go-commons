package logging

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
)

// EchoLogger is a middleware and zap to provide an "access log" like logging for each request.
func EchoLogger(logger *zap.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogRequestID: true,
		LogLatency:   true,
		LogHost:      true,
		LogRemoteIP:  true,
		LogReferer:   true,
		LogURIPath:   true,
		LogUserAgent: true,
		LogMethod:    true,
		LogStatus:    true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			req := c.Request()

			var traceID string

			if span := opentracing.SpanFromContext(req.Context()); span != nil {
				if sc, ok := span.Context().(jaeger.SpanContext); ok {
					traceID = sc.TraceID().String()
				}
			}

			logger.Info("request",
				zap.String("request_id", v.RequestID),
				zap.Duration("latency", v.Latency),
				zap.String("host", v.Host),
				zap.String("remote_ip", v.RemoteIP),
				zap.String("referer", v.Referer),
				zap.String("uri_path", v.URIPath),
				zap.String("user_agent", v.UserAgent),
				zap.String("method", v.Method),
				zap.Int("status", v.Status),
				zap.String("trace_id", traceID),
			)
			return nil
		},
	})
}
