package custom_connect

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"connectrpc.com/connect"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type errInterceptor struct{} // WrapStreamingClient implements connect.Interceptor.
func (*errInterceptor) WrapStreamingClient(handler connect.StreamingClientFunc) connect.StreamingClientFunc {
	return handler
}

// WrapStreamingHandler implements connect.Interceptor.
func (*errInterceptor) WrapStreamingHandler(handler connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return handler
}

// WrapUnary implements connect.Interceptor.
func (*errInterceptor) WrapUnary(handler connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, ar connect.AnyRequest) (connect.AnyResponse, error) {

		res, err := handler(ctx, ar)
		if err == nil {
			return res, err
		}

		ar.Any()

		span := trace.SpanFromContext(ctx)

		span.SetAttributes(attribute.String("rpc.error.message", err.Error()))

		for key, value := range ar.Header() {
			span.SetAttributes(attribute.String(fmt.Sprintf("rpc.headers.%s", key), strings.Join(value, ",")))
		}

		// handle protobuf request
		if msg, ok := ar.Any().(proto.Message); ok {
			b, _ := protojson.MarshalOptions{
				UseProtoNames:   true,
				EmitUnpopulated: false,
			}.Marshal(msg)

			span.SetAttributes(attribute.String("rpc.request.body", string(b)))
		} else {
			// handle native JSON
			b, _ := json.Marshal(ar.Any())
			span.SetAttributes(attribute.String("rpc.request.body", string(b)))
		}

		return res, err
	}
}
