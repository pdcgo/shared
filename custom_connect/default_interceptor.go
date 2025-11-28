package custom_connect

import (
	"context"

	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"
	"connectrpc.com/validate"
	"github.com/gin-gonic/gin"
	"github.com/pdcgo/shared/custom_logging"
)

type DefaultInterceptor connect.HandlerOption

func NewDefaultInterceptor() (DefaultInterceptor, error) {
	interceptor := validate.NewInterceptor()

	telemetryInterceptor, err := otelconnect.NewInterceptor(
		otelconnect.WithTrustRemote(),
		otelconnect.WithPropagateResponseHeader(),
	)

	if err != nil {
		return nil, err
	}

	defaultInterceptor := connect.WithInterceptors(
		telemetryInterceptor,
		&RequestSourceInterceptor{},
		&errInterceptor{},
		&custom_logging.LoggingInterceptor{},
		interceptor,
		&custom_logging.DBLoggingInterceptor{},
	)

	return defaultInterceptor, nil
}

type DefaultClientInterceptor connect.ClientOption

func NewDefaultClientInterceptor() (DefaultClientInterceptor, error) {
	validator := validate.NewInterceptor()

	telemetryInterceptor, err := otelconnect.NewInterceptor(
		otelconnect.WithTrustRemote(),
		otelconnect.WithPropagateResponseHeader(),
	)

	if err != nil {
		return nil, err
	}

	return connect.WithInterceptors(
		&ginInterceptor{},
		&RequestSourceInterceptor{},
		validator,
		telemetryInterceptor,
		&errInterceptor{},
	), nil
}

type ginInterceptor struct{}

// WrapStreamingClient implements connect.Interceptor.
func (g *ginInterceptor) WrapStreamingClient(handler connect.StreamingClientFunc) connect.StreamingClientFunc {
	// return func(ctx context.Context, s connect.Spec) connect.StreamingClientConn {
	// 	if !s.IsClient {
	// 		return handler(ctx, s)
	// 	}

	// 	return handler(ctx, s)
	// }

	return handler
}

// WrapStreamingHandler implements connect.Interceptor.
func (g *ginInterceptor) WrapStreamingHandler(handler connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {

	// return func(ctx context.Context, shc connect.StreamingHandlerConn) error {
	// 	if !shc.Spec().IsClient {
	// 		return handler(ctx, shc)
	// 	}

	// 	return handler(ctx, shc)
	// }
	return handler
}

// WrapUnary implements connect.Interceptor.
func (g *ginInterceptor) WrapUnary(handler connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, ar connect.AnyRequest) (connect.AnyResponse, error) {
		if !ar.Spec().IsClient {
			return handler(ctx, ar)
		}

		ginCtx, ok := ctx.(*gin.Context)
		if !ok {
			return handler(ctx, ar)
		}

		token := ginCtx.Request.Header.Get("Authorization")
		if token == "" {
			return handler(ctx, ar)
		}

		ar.Header().Set("Authorization", token)
		return handler(ctx, ar)
	}
}
