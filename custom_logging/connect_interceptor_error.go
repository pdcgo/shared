package custom_logging

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"
)

// LoggingInterceptor logs errors from RPC calls
type LoggingInterceptor struct{}

// WrapUnary satisfies connect.Interceptor
func (l *LoggingInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		res, err := next(ctx, req)
		if err != nil {
			if cErr, ok := err.(*connect.Error); ok {
				slog.Error("rpc error",
					"procedure", req.Spec().Procedure,
					"code", cErr.Code().String(),
					"msg", cErr.Message(),
				)
			} else {
				slog.Error("request_error",
					"procedure", req.Spec().Procedure,
					"error", err,
					"message", err.Error(),
					"token", req.Header().Get("Authorization"),
					slog.Any("payload", req.Any()),
				)
			}
		}
		return res, err
	}
}

// WrapStreamingClient (optional: if you want streaming client logs)
func (l *LoggingInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

// WrapStreamingHandler (optional: if you want streaming server logs)
func (l *LoggingInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return next
}
