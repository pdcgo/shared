package tracking_iface

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pdcgo/shared/db_models"
	"github.com/pdcgo/v2_gots_sdk/pdc_api"
)

type UnsupportedTrackErr struct {
	ShippingID uint
}

type MetaApi interface {
	Handler() gin.HandlerFunc
	ApiMeta(uri string) *pdc_api.Api
}

// Error implements error.
func (u *UnsupportedTrackErr) Error() string {
	return fmt.Sprintf("shipping id %d not supported", u.ShippingID)
}

type ThirdParty interface {
	Track(shippingID uint, receipt string) (*db_models.TrackInfo, error)
}

type TrackingService interface {
	Track(shippingID uint, receipt string) (*db_models.TrackInfo, error)
}

type TrackingApi interface {
	Track(service TrackingService) MetaApi
	TrackOrder() MetaApi
}
