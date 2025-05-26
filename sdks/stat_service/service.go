package stat_service

import "github.com/pdcgo/shared/interfaces/stat_iface"

type statServiceSdkImpl struct{}

// Pipeline implements stat_iface.StatService.
func (s *statServiceSdkImpl) Pipeline() stat_iface.Pipeline {
	panic("unimplemented")
}

// MarketplaceHoldFundPipeline implements stat_iface.StatService.
func (s *statServiceSdkImpl) MarketplaceHoldFundPipeline() stat_iface.MarketplaceHoldFundPipeline {
	panic("unimplemented")
}

func NewStatService() stat_iface.StatService {
	return &statServiceSdkImpl{}
}
