package logging

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
)

type key int

const loggerIDKey key = 119

// LoggerInContext will embed zap.Logger into request context for handler to use
func LoggerInContext(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			c.SetRequest(req.WithContext(addLoggerToContext(req.Context(), logger)))
			return next(c)
		}
	}
}

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
		LogHeaders:   []string{"X-Wikia-Internal-Request","X-Original-Host"},
		BeforeNextFunc: func(c echo.Context) {
			var traceID string
			req := c.Request()

			if span := opentracing.SpanFromContext(req.Context()); span != nil {
				if sc, ok := span.Context().(jaeger.SpanContext); ok {
					traceID = sc.TraceID().String()
				}
			}

			wrapped := logger.With(zap.String("trace_id", traceID))
			c.SetRequest(req.WithContext(addLoggerToContext(req.Context(), wrapped)))
		},
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log := FromEchoContext(c)
			wrapped := log.With(
				zap.String("request_id", v.RequestID),
				zap.Duration("latency", v.Latency),
				zap.String("host", v.Host),
				zap.String("remote_ip", v.RemoteIP),
				zap.String("referer", v.Referer),
				zap.String("uri_path", v.URIPath),
				zap.String("user_agent", v.UserAgent),
				zap.String("method", v.Method),
				zap.Int("status", v.Status),
				zap.Any("headers", v.Headers),
			)

			wrapped.Info("request")
			return nil
		},
	})
}

func FromEchoContext(ctx echo.Context) *zap.Logger {
	return FromContext(ctx.Request().Context())
}

// FromRequest will return current logger embedded in the given request object
func FromRequest(r *http.Request) *zap.Logger {
	return FromContext(r.Context())
}

// FromContext will return current logger from the given context.Context object
func FromContext(ctx context.Context) *zap.Logger {
	logger := ctx.Value(loggerIDKey)
	if logger == nil {
		return nil
	}

	return logger.(*zap.Logger)
}

func AddToContext(c echo.Context, logger *zap.Logger) {
	c.SetRequest(c.Request().WithContext(addLoggerToContext(c.Request().Context(), logger)))
}

// addLoggerToContext adds given logger to the context.Context and returns new context
func addLoggerToContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerIDKey, logger)
}
