package moretest_mock

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pdcgo/shared/pkg/moretest"
	"github.com/pdcgo/v2_gots_sdk"
)

func SetupSdk(group *v2_gots_sdk.SdkGroup, relativeuri string) moretest.SetupFunc {
	return func(t *testing.T) func() error {
		r := gin.Default()
		sdk := v2_gots_sdk.NewApiSdk(r)
		g := sdk.Group(relativeuri)
		*group = *g
		return nil
	}
}
