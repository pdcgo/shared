package custom_connect

import (
	"context"
	"encoding/base64"
	"errors"

	"buf.build/go/protovalidate"
	"connectrpc.com/connect"
	"github.com/pdcgo/schema/services/access_iface/v1"
	"google.golang.org/protobuf/proto"
)

// var _ connect.Interceptor = (*RequestSource)(nil)

type RequestSourceIntercept struct{}

type KeyContext string

const (
	SourceKey KeyContext = "source"
)

// WrapStreamingClient implements connect.Interceptor.
func (r *RequestSourceIntercept) WrapStreamingClient(handler connect.StreamingClientFunc) connect.StreamingClientFunc {
	return handler
}

// WrapStreamingHandler implements connect.Interceptor.
func (r *RequestSourceIntercept) WrapStreamingHandler(handler connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return handler
}

// WrapUnary implements connect.Interceptor.
func (r *RequestSourceIntercept) WrapUnary(handler connect.UnaryFunc) connect.UnaryFunc {
	validator := protovalidate.GlobalValidator

	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {

		raw := req.Header().Get("X-Pdc-Source")
		if raw == "" {
			return connect.NewResponse(&access_iface.RequestSourceError{
				Message: "invalid base64 source",
			}), errors.New("empty source")
		}

		data, err := base64.StdEncoding.DecodeString(raw)

		if err != nil {
			return connect.NewResponse(&access_iface.RequestSourceError{
				Message: "invalid base64 source",
			}), err
		}

		source := &access_iface.RequestSource{}
		err = proto.Unmarshal(data, source)
		if err != nil {
			return connect.NewResponse(&access_iface.RequestSourceError{
				Message: "proto cannot decode source",
			}), err
		}

		err = validator.Validate(source)
		if err != nil {
			return connect.NewResponse(&access_iface.RequestSourceError{
				Message: "incomplete source",
			}), err
		}

		cctx := context.WithValue(ctx, SourceKey, source)
		return handler(cctx, req)
	}
}

func GetRequestSource(ctx context.Context) *access_iface.RequestSource {
	return ctx.Value(SourceKey).(*access_iface.RequestSource)
}

type RequestSourceClientIntercept struct {
	// source *access_iface.RequestSource
}

func RequestSourceSerialize(msg *access_iface.RequestSource) (string, error) {

	raw, err := proto.Marshal(msg)
	if err != nil {
		return "", err
	}

	var encode []byte = make([]byte, base64.StdEncoding.EncodedLen(len(raw)))
	base64.StdEncoding.Encode(encode, raw)

	return string(encode), err
}
