package stat_service_test

import (
	"testing"
	"time"

	"github.com/pdcgo/shared/pkg/debugtool"
	"github.com/pdcgo/shared/sdks/stat_service"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationStatService(t *testing.T) {
	t.Skip("cuma integration test")
	sdk := stat_service.NewStatService("http://localhost:8080/stat/")
	t.Run("testing pipeline", func(t *testing.T) {
		pipe := sdk.Pipeline()
		t.Run("testing start", func(t *testing.T) {
			res := pipe.Start(t.Context(), time.Time{})
			assert.Nil(t, res.Err)
			debugtool.LogJson(res)
		})

		// t.Run("testing cancel", func(t *testing.T) {
		// 	res := pipe.Cancel()
		// 	assert.Nil(t, res.Err)
		// })
	})
}
