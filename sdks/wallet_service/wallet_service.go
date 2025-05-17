package wallet_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/pdcgo/shared/interfaces/identity_iface"
	"github.com/pdcgo/shared/interfaces/wallet_iface"
)

var ErrReq = wallet_iface.ResultErr{
	Code:    "sdk_request_error",
	Message: "service wallet sedang error",
}

var ErrSdk = wallet_iface.ResultErr{
	Code:    "sdk_error",
	Message: "service wallet sedang error",
}

var ErrDeprecatedOrNotFound = wallet_iface.ResultErr{
	Code:    "sdk_deprecated_or_not_found",
	Message: "",
}

type httpWalletImpl struct {
	config     *wallet_iface.WalletServiceConfig
	httpClient *httpClient
}

func NewWalletService(config *wallet_iface.WalletServiceConfig) wallet_iface.WalletService {
	return &httpWalletImpl{
		config: config,
		httpClient: &httpClient{
			secret:   config.AppSecret,
			endpoint: config.Endpoint,
		},
	}
}

// CancelPayment implements wallet_iface.WalletService.
func (h *httpWalletImpl) CancelPayment(ctx context.Context, agent identity_iface.Agent) wallet_iface.CancelPaymentRes {
	panic("unimplemented")
}

// CreatePayment implements wallet_iface.WalletService.
func (h *httpWalletImpl) CreatePayment(ctx context.Context, agent identity_iface.Agent) wallet_iface.CreatePaymentRes {
	panic("unimplemented")
}

// CreateWallet implements wallet_iface.WalletService.
func (h *httpWalletImpl) CreateWallet(ctx context.Context, agent identity_iface.Agent, teamID uint) *wallet_iface.WalletInfoRes {
	var res wallet_iface.WalletInfoRes

	req, err := h.httpClient.createReq(ctx, agent, http.MethodPost, fmt.Sprintf("wallets/%d", teamID), nil, nil)
	if err != nil {
		res.Err = &ErrSdk
	}

	httpres, err := h.httpClient.Do(req)
	if err != nil {
		var reserr *wallet_iface.ResultErr
		if errors.As(err, &reserr) {
			res.Err = reserr
		} else {
			res.Err = &ErrReq
		}

		return &res
	}

	err = json.NewDecoder(httpres.Body).Decode(&res)
	if err != nil {
		res.Err = &ErrSdk
		return &res
	}
	return &res
}

// GetPayment implements wallet_iface.WalletService.
func (h *httpWalletImpl) GetPayment(ctx context.Context, agent identity_iface.Agent) wallet_iface.PaymentDetailRes {
	panic("unimplemented")
}

// GetTransactions implements wallet_iface.WalletService.
func (h *httpWalletImpl) GetTransactions(ctx context.Context, agent identity_iface.Agent) wallet_iface.GetTransactionsRes {
	panic("unimplemented")
}

// ListPayment implements wallet_iface.WalletService.
func (h *httpWalletImpl) ListPayment(ctx context.Context, agent identity_iface.Agent) wallet_iface.ListPaymentRes {
	panic("unimplemented")
}

// TeamWallet implements wallet_iface.WalletService.
func (h *httpWalletImpl) TeamWallet(ctx context.Context, agent identity_iface.Agent, teamID uint, tipe wallet_iface.BalanceType) wallet_iface.TeamWallet {
	res := &teamWalletImpl{
		teamID:     teamID,
		tipe:       tipe,
		agent:      agent,
		httpClient: h.httpClient,
		ctx:        ctx,
	}
	return res
}

// TransactionCreate implements wallet_iface.WalletService.
func (h *httpWalletImpl) TransactionCreate(ctx context.Context, agent identity_iface.Agent, payload *wallet_iface.TransactionCreatePayload) wallet_iface.CreateTransactionRes {
	panic("unimplemented")
}

// UpdateTransaction implements wallet_iface.WalletService.
func (h *httpWalletImpl) UpdateTransaction(ctx context.Context, agent identity_iface.Agent) wallet_iface.UpdateTransactionRes {
	panic("unimplemented")
}
