package custom_connect

import (
	"context"
	"log"

	"connectrpc.com/connect"
)

var _ connect.Interceptor = (*ScopeIntercept)(nil)

type ScopeIntercept struct{}

// WrapStreamingClient implements connect.Interceptor.
func (s *ScopeIntercept) WrapStreamingClient(handler connect.StreamingClientFunc) connect.StreamingClientFunc {
	return handler
}

// WrapStreamingHandler implements connect.Interceptor.
func (s *ScopeIntercept) WrapStreamingHandler(handler connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return handler
}

// WrapUnary implements connect.Interceptor.
func (s *ScopeIntercept) WrapUnary(handler connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, ar connect.AnyRequest) (connect.AnyResponse, error) {
		scope := ar.Peer().Query.Get("resource_scope_id")

		log.Println(scope)
		return handler(ctx, ar)
	}
}
