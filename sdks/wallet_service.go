package sdks

import (
	"github.com/pdcgo/shared/interfaces/identity_iface"
	"github.com/pdcgo/shared/interfaces/wallet_iface"
)

type httpWalletImpl struct {
	endpoint string
}

// CancelPayment implements wallet_iface.WalletService.
func (h *httpWalletImpl) CancelPayment(agent identity_iface.Agent) wallet_iface.CancelPaymentRes {
	panic("unimplemented")
}

// CreatePayment implements wallet_iface.WalletService.
func (h *httpWalletImpl) CreatePayment(agent identity_iface.Agent) wallet_iface.CreatePaymentRes {
	panic("unimplemented")
}

// CreateWallet implements wallet_iface.WalletService.
func (h *httpWalletImpl) CreateWallet(agent identity_iface.Agent, teamID uint) *wallet_iface.WalletInfoRes {
	panic("unimplemented")
}

// GetPayment implements wallet_iface.WalletService.
func (h *httpWalletImpl) GetPayment(agent identity_iface.Agent) wallet_iface.PaymentDetailRes {
	panic("unimplemented")
}

// GetTransactions implements wallet_iface.WalletService.
func (h *httpWalletImpl) GetTransactions(agent identity_iface.Agent) wallet_iface.GetTransactionsRes {
	panic("unimplemented")
}

// ListPayment implements wallet_iface.WalletService.
func (h *httpWalletImpl) ListPayment(agent identity_iface.Agent) wallet_iface.ListPaymentRes {
	panic("unimplemented")
}

// TeamWallet implements wallet_iface.WalletService.
func (h *httpWalletImpl) TeamWallet(agent identity_iface.Agent, teamID uint, tipe wallet_iface.BalanceType) wallet_iface.TeamWallet {
	panic("unimplemented")
}

// TransactionCreate implements wallet_iface.WalletService.
func (h *httpWalletImpl) TransactionCreate(agent identity_iface.Agent, payload *wallet_iface.TransactionCreatePayload) wallet_iface.CreateTransactionRes {
	panic("unimplemented")
}

// UpdateTransaction implements wallet_iface.WalletService.
func (h *httpWalletImpl) UpdateTransaction(agent identity_iface.Agent) wallet_iface.UpdateTransactionRes {
	panic("unimplemented")
}

func NewWalletService(endpoint string) wallet_iface.WalletService {
	return &httpWalletImpl{
		endpoint: endpoint,
	}
}
