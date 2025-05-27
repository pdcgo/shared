package stat_service

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/pdcgo/shared/interfaces/stat_iface"
)

type statServiceSdkImpl struct {
	endpoint string
}

type pipelineImpl struct {
	endpoint string
}

// Cancel implements stat_iface.Pipeline.
func (p *pipelineImpl) Cancel() *stat_iface.ExecRes {
	var res stat_iface.ExecRes
	httpres, err := http.DefaultClient.Get(p.endpoint + "pipeline/cancel")
	if err != nil {
		return res.SetError(err)
	}

	err = json.NewDecoder(httpres.Body).Decode(&res)
	if err != nil {
		return res.SetError(err)
	}

	return res.DetectError()
}

// Start implements stat_iface.Pipeline.
func (p *pipelineImpl) Start(ctx context.Context, startTime time.Time) *stat_iface.ExecRes {
	var res stat_iface.ExecRes
	httpres, err := http.DefaultClient.Get(p.endpoint + "pipeline/start")
	if err != nil {
		return res.SetError(err)
	}

	err = json.NewDecoder(httpres.Body).Decode(&res)
	if err != nil {
		return res.SetError(err)
	}

	return res.DetectError()
}

// Wait implements stat_iface.Pipeline.
func (p *pipelineImpl) Wait(ctx context.Context) *stat_iface.ExecRes {
	panic("unimplemented")
}

// Pipeline implements stat_iface.StatService.
func (s *statServiceSdkImpl) Pipeline() stat_iface.Pipeline {
	return &pipelineImpl{
		endpoint: s.endpoint,
	}
}

// MarketplaceHoldFundPipeline implements stat_iface.StatService.
func (s *statServiceSdkImpl) MarketplaceHoldFundPipeline() stat_iface.MarketplaceHoldFundPipeline {
	panic("unimplemented")
}

func NewStatService(endpoint string) stat_iface.StatService {
	return &statServiceSdkImpl{
		endpoint: endpoint,
	}
}
