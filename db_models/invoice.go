package db_models

import (
	"time"
)

type InvoiceStatus string

const (
	InvoicePaid     InvoiceStatus = "paid"
	InvoiceNotPaid  InvoiceStatus = "not_paid"
	InvoiceNotFinal InvoiceStatus = "not_final"
)

func (InvoiceStatus) EnumList() []string {
	return []string{
		"paid",
		"not_paid",
	}
}

type InvoiceType string

func (InvoiceType) EnumList() []string {
	return []string{
		"shipping_fee",
		"product",
		"warehouse_fee",
		"prod_adjust",
		"ware_adjust",
		"common_adjust",
		"problem",
		"problem_adjust",
	}
}

const (
	InvoShipFeeType         InvoiceType = "shipping_fee"
	InvoProductType         InvoiceType = "product"
	InvoWarehouseFeeType    InvoiceType = "warehouse_fee"
	InvoProductAdjustment   InvoiceType = "prod_adjust"
	InvoWarehouseAdjustment InvoiceType = "ware_adjust"
	InvoCommonAdjustment    InvoiceType = "common_adjust"
	InvoProblem             InvoiceType = "problem"
	InvoProblemAdjustment   InvoiceType = "problem_adjust"
)

type Invoice struct {
	ID         uint  `json:"id" gorm:"primarykey"`
	OrderID    *uint `json:"order_id"`
	TxID       *uint `json:"tx_id"`
	FromTeamID uint  `json:"from_team_id"`
	ToTeamID   uint  `json:"to_team_id"`
	HistID     *uint `json:"hist_id"`

	// Type       string        `json:"type"`
	Status        InvoiceStatus `json:"status" gorm:"index"`
	Amount        float64       `json:"amount"`
	PaidAt        time.Time     `json:"paid_at"` // manut tanggal submission
	AcceptedAt    time.Time     `json:"accepted_at"`
	Created       time.Time     `json:"created"`
	Type          InvoiceType   `json:"type"`
	NeedAdj       bool          `json:"need_adj"`
	HasSubmission bool          `json:"has_submission"`

	Order              *Order               `json:"order,omitempty"`
	Tx                 *InvTransaction      `json:"inv_transaction,omitempty"`
	FromTeam           *Team                `json:"-" gorm:"foreignkey:FromTeamID"`
	ToTeam             *Team                `json:"-" gorm:"foreignkey:ToTeamID"`
	Hist               *PaymentHistory      `json:"hist"`
	PaymentSubmissions []*PaymentSubmission `json:"payment_submission,omitempty" gorm:"many2many:invoice_payment_submission;"`
}

// GetEntityID implements authorization.Entity.
func (i *Invoice) GetEntityID() string {
	return "invoices"
}

func (d *Invoice) GetAdjustmentType() InvoiceType {
	var newtype InvoiceType = InvoCommonAdjustment

	switch d.Type {
	case InvoProductType:
		newtype = InvoProductAdjustment
	case InvoWarehouseFeeType:
		newtype = InvoWarehouseAdjustment
	case InvoProblem:
		newtype = InvoProblemAdjustment
	default:
		newtype = InvoCommonAdjustment
	}

	return newtype
}

type PaymentHistory struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	CreatedByID uint      `json:"created_by_id"`
	Amount      float64   `json:"amount"`
	CreatedAt   time.Time `json:"created_at"`

	CreatedBy *User `json:"createdBy"`
}

type PaymentSubmissionStatus string

const (
	PaymentSubmissionStatusSubmitted PaymentSubmissionStatus = "submitted"
	PaymentSubmissionStatusAccepted  PaymentSubmissionStatus = "accepted"
	PaymentSubmissionStatusRejected  PaymentSubmissionStatus = "rejected"
)

func (PaymentSubmissionStatus) EnumList() []string {
	return []string{
		"submitted",
		"accepted",
		"rejected",
	}
}

func (s PaymentSubmissionStatus) IsValid() bool {
	values := s.EnumList()
	for _, value := range values {
		if value == string(s) {
			return true
		}
	}

	return false
}

type PaymentSubmission struct {
	ID            uint                    `json:"id" gorm:"primarykey"`
	CreatedByID   uint                    `json:"created_by_id"`
	CompletedByID *uint                   `json:"completed_by_id"`
	Status        PaymentSubmissionStatus `json:"status"`
	Receipt       string                  `json:"receipt"`
	Amount        float64                 `json:"amount"`
	CreatedAt     time.Time               `json:"created_at"`

	CreatedBy   *User `json:"created_by"`
	CompletedBy *User `json:"completed_by"`
}

type PSubmissionInv struct {
	InvoiceID           uint
	PaymentSubmissionID uint

	Invoice           *Invoice           `json:"invoice"`
	PaymentSubmission *PaymentSubmission `json:"payment_submission"`
}

func (PSubmissionInv) TableName() string {
	return "invoice_payment_submission"
}
