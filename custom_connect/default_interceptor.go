package custom_connect

import (
	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"
	"connectrpc.com/validate"
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
		validator,
		telemetryInterceptor,
	), nil
}
