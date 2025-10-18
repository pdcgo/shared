package custom_logging

import (
	"context"

	"connectrpc.com/connect"
)

type DBLoggingInterceptor struct{}

// WrapStreamingClient implements connect.Interceptor.
func (d *DBLoggingInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

// WrapStreamingHandler implements connect.Interceptor.
func (d *DBLoggingInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return next
}

// WrapUnary implements connect.Interceptor.
func (d *DBLoggingInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(pctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		ctx := context.WithValue(pctx, "route", req.Spec().Procedure)
		return next(ctx, req)
	}
}
