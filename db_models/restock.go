package db_models

import "time"

type RestockStatus string

func (RestockStatus) EnumList() []string {
	return []string{
		"owner_review",
		"admin_review",
		"reject",
		"amount_err",
		"paid",
	}
}

const (
	RestockOwnerReview RestockStatus = "owner_review"
	RestockAdminReview RestockStatus = "admin_review"
	RestockReject      RestockStatus = "reject"
	RestockAmountErr   RestockStatus = "amount_err"
	RestockPaid        RestockStatus = "paid"
)

type Restock struct {
	ID               uint   `json:"id" gorm:"primarykey"`
	TeamID           uint   `json:"team_id"`
	UserID           uint   `json:"user_id"`
	WarehouseID      uint   `json:"warehouse_id"`
	ShippingID       uint   `json:"shipping_id"`
	InvTransactionID *uint  `json:"tx_id"`
	ExternOrdID      string `json:"extern_ord_id"`

	Receipt   string        `json:"receipt"`
	ReStatus  RestockStatus `json:"re_status"`
	CreatedAt time.Time     `json:"created_at"`
	Deleted   bool          `json:"deleted" gorm:"index"`

	Items          []*RestockItem        `json:"items"`
	Invoices       []*RestockInvoiceItem `json:"invoices"`
	InvTransaction *InvTransaction       `json:"tx"`
	Team           *Team                 `json:"team"`
	User           *User                 `json:"user"`
}

// GetEntityID implements authorization.Entity.
func (r *Restock) GetEntityID() string {
	return "restock_submission"
}

type RestockInvoType string

func (RestockInvoType) EnumList() []string {
	return []string{
		"product",
		"shipping",
		"app_fee",
		"discount",
		"custom",
	}
}

const (
	RestockInvoProduct     RestockInvoType = "product"
	RestockInvoShippingFee RestockInvoType = "shipping"
	RestockInvoAppFee      RestockInvoType = "app_fee"
	RestockInvoDiscount    RestockInvoType = "discount"
	RestockInvoCustom      RestockInvoType = "custom"
)

type RestockInvoiceItem struct {
	ID               uint  `json:"id" gorm:"primarykey"`
	RestockID        uint  `json:"restock_id"`
	InvTransactionID *uint `json:"tx_id"`

	InvType RestockInvoType `json:"inv_type"`
	Label   string          `json:"label"`
	Amount  float64         `json:"amount"`
}

type RestockItem struct {
	ID        uint `json:"id" gorm:"primarykey"`
	RestockID uint `json:"restock_id"`
	VariantID uint `json:"variant_id"`

	Count   int             `json:"count"`
	Price   float64         `json:"price"`
	Variant *VariationValue `json:"variant"`
}

type RestockSuppplierTemp struct {
	ID        uint   `json:"id" gorm:"primarykey"`
	RestockID uint   `json:"restock_id"`
	VariantID uint   `json:"variant_id" binding:"required"`
	Link      string `json:"link" binding:"required,lte=500"`

	Restock *Restock
}

type RestockPaymentType string

const (
	RestockPaymentBankAccount RestockPaymentType = "bank_account"
	RestockPaymentShopeePay   RestockPaymentType = "shopee_pay"
	RestockPaymentNoPayment   RestockPaymentType = "no_payment"
)

func (RestockPaymentType) EnumList() []string {
	return []string{
		"bank_account",
		"shopee_pay",
	}
}

func (t RestockPaymentType) IsValid() bool {
	values := t.EnumList()
	for _, value := range values {
		if value == string(t) {
			return true
		}
	}

	return false
}

type RestockCost struct {
	ID               uint               `json:"id" gorm:"primarykey"`
	InvTransactionID uint               `json:"inv_transaction_id"`
	PaymentType      RestockPaymentType `json:"payment_type"`
	ShippingFee      float64            `json:"shipping_fee"`
	CodFee           float64            `json:"cod_fee"`
	OtherFee         float64            `json:"other_fee"`
	PerPieceFee      float64            `json:"per_piece_fee"`
}

func (r *RestockCost) TotalFee() float64 {
	fee := r.ShippingFee + r.OtherFee + r.CodFee
	return fee
}

func (r *RestockCost) CalculatePerPiece(count int) {
	fee := r.TotalFee()
	pieceFee := fee / float64(count)

	r.PerPieceFee = pieceFee
}
