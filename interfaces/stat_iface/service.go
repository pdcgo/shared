package stat_iface

import (
	"context"
	"errors"
	"time"
)

type ExecRes struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *ExecRes) DetectError() *ExecRes {
	if e.Message != "" {
		e.Err = errors.New(e.Message)
	}
	return e
}

func (e *ExecRes) SetError(err error) *ExecRes {
	if err != nil {
		e.Err = err
		e.Message = err.Error()
	}

	return e
}

type MarketplaceHoldFundPipeline interface {
	Start(startTime time.Time) *ExecRes
	Wait() *ExecRes
	Stop() *ExecRes
	Info() *ExecRes
}

type Pipeline interface {
	Cancel() *ExecRes
	Start(ctx context.Context, startTime time.Time) *ExecRes
	Wait(ctx context.Context) *ExecRes
}

type StatService interface {
	MarketplaceHoldFundPipeline() MarketplaceHoldFundPipeline
	Pipeline() Pipeline
}
