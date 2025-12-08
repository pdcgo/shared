package custom_connect

import (
	"context"

	"connectrpc.com/connect"
)

type HaveDispatch interface {
	GetDispatch() bool
}

// var _ connect.Interceptor = (*dispatchImpl)(nil)

type dispatchImpl struct{}

// WrapStreamingClient implements connect.Interceptor.
func (d *dispatchImpl) WrapStreamingClient(handler connect.StreamingClientFunc) connect.StreamingClientFunc {
	return handler
}

// WrapStreamingHandler implements connect.Interceptor.
func (d *dispatchImpl) WrapStreamingHandler(handler connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return handler
}

// WrapUnary implements connect.Interceptor.
func (d *dispatchImpl) WrapUnary(handler connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, ar connect.AnyRequest) (connect.AnyResponse, error) {

		if ar.Spec().IsClient {
			pay, ok := ar.Any().(HaveDispatch)
			if !ok {
				return handler(ctx, ar)
			}

			if pay.GetDispatch() {

				// return nil, connect.NewError(connect.code)
			}

			return handler(ctx, ar)
		}

		return handler(ctx, ar)
	}
}

func NewDispatchInterceptor() *dispatchImpl {
	return &dispatchImpl{}
}
