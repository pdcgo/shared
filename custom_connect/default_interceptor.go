package custom_connect

import (
	"connectrpc.com/connect"
	"connectrpc.com/validate"
	"github.com/pdcgo/shared/custom_logging"
)

type DefaultInterceptor connect.HandlerOption

func NewDefaultInterceptor() (DefaultInterceptor, error) {
	interceptor, err := validate.NewInterceptor()
	if err != nil {
		return nil, err
	}

	defaultInterceptor := connect.WithInterceptors(
		interceptor,
		&custom_logging.LoggingInterceptor{},
		&custom_logging.DBLoggingInterceptor{},
	)

	return defaultInterceptor, nil
}
