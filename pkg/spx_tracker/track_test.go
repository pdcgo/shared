package spx_tracker_test

import (
	"net/http"
	"testing"

	"github.com/pdcgo/shared/pkg/spx_tracker"
	"github.com/stretchr/testify/assert"
)

func TestSpxTracker(t *testing.T) {

	c := spx_tracker.TrackClient{
		C: http.DefaultClient,
	}
	// c.Track("SPXID054339093423")
	_, err := c.Track("JX3660411373")
	assert.NotNil(t, err)

	// debugtool.LogJson(res)

}
