package db_models

type OrderPayment struct {
	ID                   uint `json:"id" gorm:"primarykey"`
	OrderID              uint `json:"order_id" gorm:"index:order_id_unique,unique"`
	IsReceivableAdjusted bool `json:"is_receivable_adjusted"`

	Order *Order `gorm:"-"`
}
