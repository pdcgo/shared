package custom_connect

import (
	"connectrpc.com/connect"
	"connectrpc.com/validate"
	"github.com/pdcgo/shared/custom_logging"
)

type DefaultInterceptor connect.HandlerOption

func NewDefaultInterceptor() (DefaultInterceptor, error) {
	interceptor := validate.NewInterceptor()

	defaultInterceptor := connect.WithInterceptors(
		&custom_logging.LoggingInterceptor{},
		interceptor,
		&custom_logging.DBLoggingInterceptor{},
	)

	return defaultInterceptor, nil
}
