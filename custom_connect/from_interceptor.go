package custom_connect

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	"buf.build/go/protovalidate"
	"connectrpc.com/connect"
	"github.com/pdcgo/schema/services/access_iface/v1"
	"google.golang.org/protobuf/proto"
)

// var _ connect.Interceptor = (*RequestSource)(nil)

type RequestSourceInterceptor struct{}

type KeyContext string

const (
	SourceKey    KeyContext = "source"
	AuthTokenKey KeyContext = "auth_token"
)

type failingStream struct {
	err error
}

// CloseRequest implements connect.StreamingClientConn.
func (f *failingStream) CloseRequest() error {
	return f.err
}

// CloseResponse implements connect.StreamingClientConn.
func (f *failingStream) CloseResponse() error {
	return f.err
}

// Peer implements connect.StreamingClientConn.
func (f *failingStream) Peer() connect.Peer {
	return connect.Peer{}
}

// Receive implements connect.StreamingClientConn.
func (f *failingStream) Receive(any) error {
	return f.err
}

// RequestHeader implements connect.StreamingClientConn.
func (f *failingStream) RequestHeader() http.Header {
	return http.Header{}
}

// ResponseHeader implements connect.StreamingClientConn.
func (f *failingStream) ResponseHeader() http.Header {
	return http.Header{}
}

// ResponseTrailer implements connect.StreamingClientConn.
func (f *failingStream) ResponseTrailer() http.Header {
	return http.Header{}
}

// Send implements connect.StreamingClientConn.
func (f *failingStream) Send(any) error {
	return f.err
}

// Spec implements connect.StreamingClientConn.
func (f *failingStream) Spec() connect.Spec {
	return connect.Spec{}
}

// WrapStreamingClient implements connect.Interceptor.
func (r *RequestSourceInterceptor) WrapStreamingClient(handler connect.StreamingClientFunc) connect.StreamingClientFunc {
	return func(ctx context.Context, s connect.Spec) connect.StreamingClientConn {

		if s.IsClient {
			conn := handler(ctx, s)
			source, _ := GetRequestSource(ctx)

			if source != nil {
				sourceString, err := RequestSourceSerialize(source)
				if err != nil {
					return &failingStream{
						err: err,
					}
				}
				conn.RequestHeader().Set("X-Pdc-Source", sourceString)
				// conn.RequestHeader().Set("Authorization", )
			}

			token, _ := GetAuthToken(ctx)
			if token != "" {
				conn.RequestHeader().Set("Authorization", token)
			}

			return conn
		}

		conn := handler(ctx, s)
		return conn
	}
}

// WrapStreamingHandler implements connect.Interceptor.
func (r *RequestSourceInterceptor) WrapStreamingHandler(handler connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	validator := protovalidate.GlobalValidator
	return func(ctx context.Context, shc connect.StreamingHandlerConn) error {
		var err error

		var raw string
		raw = shc.Peer().Query.Get("x-pdc-source")
		raw = strings.TrimSpace(raw)
		if raw == "" {
			raw = shc.RequestHeader().Get("X-Pdc-Source")
		}

		if raw == "" {
			return handler(ctx, shc)
		}

		data, err := base64.StdEncoding.DecodeString(raw)

		if err != nil {
			connect.NewError(connect.CodeInvalidArgument, err)
		}

		source := &access_iface.RequestSource{}
		err = proto.Unmarshal(data, source)
		if err != nil {
			connect.NewError(connect.CodeInvalidArgument, err)
		}

		err = validator.Validate(source)
		if err != nil {
			connect.NewError(connect.CodeInvalidArgument, err)
		}

		// cctx := context.WithValue(ctx, SourceKey, source)
		cctx := SetRequestSource(ctx, source)

		token := shc.RequestHeader().Get("Authorization")
		tctx := SetAuthToken(cctx, token)
		return handler(tctx, shc)
	}
}

// WrapUnary implements connect.Interceptor.
func (r *RequestSourceInterceptor) WrapUnary(handler connect.UnaryFunc) connect.UnaryFunc {
	validator := protovalidate.GlobalValidator

	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		var err error
		if req.Spec().IsClient {
			source, _ := GetRequestSource(ctx)
			if source == nil {
				return handler(ctx, req)
			}

			sourceString, err := RequestSourceSerialize(source)
			if err != nil {
				return nil, err
			}

			// req.Peer().Query.Set("x-pdc-source", sourceString)
			req.Header().Set("X-Pdc-Source", sourceString)

			if req.Header().Get("Authorization") == "" {
				token, _ := GetAuthToken(ctx)
				if token != "" {
					req.Header().Set("Authorization", token)
				}
			}

			return handler(ctx, req)

		}
		// set token
		ctx = SetAuthToken(ctx, req.Header().Get("Authorization"))

		var raw string
		raw = req.Peer().Query.Get("x-pdc-source")
		raw = strings.TrimSpace(raw)
		if raw == "" {
			raw = req.Header().Get("X-Pdc-Source")
		}

		if raw == "" {
			return handler(ctx, req)
		}

		data, err := base64.StdEncoding.DecodeString(raw)

		if err != nil {
			return connect.NewResponse(&access_iface.RequestSourceError{
				Message: "invalid base64 source",
			}), connect.NewError(connect.CodeInvalidArgument, err)
		}

		source := &access_iface.RequestSource{}
		err = proto.Unmarshal(data, source)
		if err != nil {
			return connect.NewResponse(&access_iface.RequestSourceError{
				Message: "proto cannot decode source",
			}), connect.NewError(connect.CodeInvalidArgument, err)
		}

		err = validator.Validate(source)
		if err != nil {
			return connect.NewResponse(&access_iface.RequestSourceError{
				Message: "incomplete source",
			}), connect.NewError(connect.CodeInvalidArgument, err)
		}

		// cctx := context.WithValue(ctx, SourceKey, source)
		cctx := SetRequestSource(ctx, source)

		return handler(cctx, req)
	}
}

func GetRequestSource(ctx context.Context) (*access_iface.RequestSource, error) {
	source, ok := ctx.Value(SourceKey).(*access_iface.RequestSource)
	if !ok {
		return source, errors.New("x-pdc-source not set")
	}

	return source, nil
}

func SetRequestSource(ctx context.Context, source *access_iface.RequestSource) context.Context {
	return context.WithValue(ctx, SourceKey, source)
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

func GetAuthToken(ctx context.Context) (string, error) {
	source, ok := ctx.Value(AuthTokenKey).(string)
	if !ok {
		return source, errors.New("token not set")
	}

	return source, nil
}

func SetAuthToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, AuthTokenKey, token)
}
