package wallet_iface

import (
	"context"
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

type TransactionCreatePayload struct {
	TeamID   uint    `json:"team_id"`
	ToTeamID uint    `json:"to_team_id"`
	TxType   TxType  `json:"tx_type"`
	Amount   float64 `json:"amount"`
	Note     string  `json:"note"`
}
type CreateTransactionRes struct {
	Err   *ResultErr           `json:"err"`
	Datas []*WalletTransaction `json:"data"`
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
	Detail() *WalletInfoRes
	SendAmount(payload *SendPayload) *CreateTransactionRes
	Err() error
}

type WalletServiceConfig struct {
	AppSecret string `yaml:"app_secret"`
	Endpoint  string `yaml:"endpoint"`
}

type WalletService interface {
	// kaitan wallet
	CreateWallet(ctx context.Context, agent identity_iface.Agent, teamID uint) *WalletInfoRes
	TeamWallet(ctx context.Context, agent identity_iface.Agent, teamID uint, tipe BalanceType) TeamWallet
	// kaitan payment
	GetPayment(ctx context.Context, agent identity_iface.Agent) PaymentDetailRes
	ListPayment(ctx context.Context, agent identity_iface.Agent) ListPaymentRes
	CreatePayment(ctx context.Context, agent identity_iface.Agent) CreatePaymentRes
	CancelPayment(ctx context.Context, agent identity_iface.Agent) CancelPaymentRes

	// kaitan transaksi
	TransactionCreate(ctx context.Context, agent identity_iface.Agent, payload *TransactionCreatePayload) CreateTransactionRes

	GetTransactions(ctx context.Context, agent identity_iface.Agent) GetTransactionsRes

	UpdateTransaction(ctx context.Context, agent identity_iface.Agent) UpdateTransactionRes
}
