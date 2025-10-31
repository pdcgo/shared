package custom_connect_test

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/pdcgo/schema/services/access_iface/v1"
	"github.com/stretchr/testify/assert"
)

func TestCustomValidate(t *testing.T) {
	validator := protovalidate.GlobalValidator
	msg := &access_iface.RequestSource{
		TeamId:      1,
		RequestFrom: access_iface.RequestFrom_REQUEST_FROM_ADMIN,
	}

	err := validator.Validate(msg)
	assert.Nil(t, err)

	msg = &access_iface.RequestSource{
		TeamId: 1,
	}

	err = validator.Validate(msg)
	assert.NotNil(t, err)

}
