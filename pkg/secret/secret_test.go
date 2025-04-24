package secret_test

import (
	"os"
	"testing"

	"github.com/pdcgo/shared/pkg/secret"
	"github.com/stretchr/testify/assert"
)

func TestSecret(t *testing.T) {
	if os.Getenv("GOOGLE_CLOUD_PROJECT") == "" {
		t.Skip("hanya yang punya kredensial")
	}

	hasil, err := secret.GetSecret("app_config_prod", "latest")
	assert.Nil(t, err)

	assert.NotEmpty(t, hasil)
}
