package custom_connect_test

import (
	"testing"

	"github.com/pdcgo/schema/services/access_iface/v1"
	"github.com/pdcgo/shared/custom_connect"
	"github.com/stretchr/testify/assert"
)

func TestSerializeSource(t *testing.T) {
	token, err := custom_connect.RequestSourceSerialize(&access_iface.RequestSource{
		TeamId:      12,
		RequestFrom: access_iface.RequestFrom_REQUEST_FROM_ADMIN,
	})

	assert.Nil(t, err)
	t.Log(token)
}
