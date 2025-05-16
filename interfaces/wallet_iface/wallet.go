package wallet_iface

import (
	"encoding/json"

	"github.com/pdcgo/shared/interfaces/identity_iface"
)

type ResultErr struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Extra   any    `json:"extra"`
}

// Error implements error.
func (r *ResultErr) Error() string {
	raw, _ := json.Marshal(r)
	return string(raw)
}

type WalletInfoRes struct {
	Err  *ResultErr `json:"err"`
	Data *Wallet    `json:"data"`
}

type PaymentDetailRes struct {
}

type ListPaymentRes struct {
}
type CreatePaymentRes struct {
}
type CancelPaymentRes struct {
}
type GetTransactionsRes struct {
}
type CreateTransactionRes struct {
}
type UpdateTransactionRes struct {
}

type WalletService interface {
	// kaitan wallet
	GetWallet(agent identity_iface.Agent, teamID uint) WalletInfoRes
	CreateWallet(agent identity_iface.Agent, teamID uint) *WalletInfoRes
	// kaitan payment
	GetPayment(agent identity_iface.Agent) PaymentDetailRes
	ListPayment(agent identity_iface.Agent) ListPaymentRes
	CreatePayment(agent identity_iface.Agent) CreatePaymentRes
	CancelPayment(agent identity_iface.Agent) CancelPaymentRes

	// kaitan transaksi
	GetTransactions(agent identity_iface.Agent) GetTransactionsRes
	CreateTransaction(agent identity_iface.Agent) CreateTransactionRes
	UpdateTransaction(agent identity_iface.Agent) UpdateTransactionRes
}
