package stat_iface

type ExecRes struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *ExecRes) SetError(err error) *ExecRes {
	if err != nil {
		e.Err = err
		e.Message = err.Error()
	}

	return e
}

type MarketplaceHoldFundPipeline interface {
	Start() *ExecRes
	Wait() *ExecRes
	Stop() *ExecRes
	Info() *ExecRes
}

type StatService interface {
	MarketplaceHoldFundPipeline() MarketplaceHoldFundPipeline
}
