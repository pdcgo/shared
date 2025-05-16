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
	Data Wallet     `json:"data"`
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

type CreateTransactionPayload struct {
	TeamID uint `json:"team_id"`

	Type   TransactionType `json:"type"`
	TxType TxType          `json:"tx_type"`
	Amount float64         `json:"amount"`
	Note   string          `json:"note"`
}
type CreateTransactionRes struct {
	Err  *ResultErr        `json:"err"`
	Data WalletTransaction `json:"data"`
}
type UpdateTransactionRes struct {
}

type SendPayload struct {
	ToTeamID    uint    `json:"to_team_id"`
	RefID       string  `json:"ref_id"`
	Type        TxType  `json:"type"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

type TeamWallet interface {
	Info() *WalletInfoRes
	SendAmount(payload *SendPayload) *CreateTransactionRes
	Err() error
}

type WalletService interface {
	// kaitan wallet
	CreateWallet(agent identity_iface.Agent, teamID uint) *WalletInfoRes
	TeamWallet(agent identity_iface.Agent, teamID uint) TeamWallet
	// kaitan payment
	GetPayment(agent identity_iface.Agent) PaymentDetailRes
	ListPayment(agent identity_iface.Agent) ListPaymentRes
	CreatePayment(agent identity_iface.Agent) CreatePaymentRes
	CancelPayment(agent identity_iface.Agent) CancelPaymentRes

	// kaitan transaksi
	GetTransactions(agent identity_iface.Agent) GetTransactionsRes
	CreateTransaction(agent identity_iface.Agent, payload *CreateTransactionPayload) CreateTransactionRes
	UpdateTransaction(agent identity_iface.Agent) UpdateTransactionRes
}
