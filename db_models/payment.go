package db_models

type PaymentType string

func (PaymentType) EnumList() []string {
	return []string{
		"bank_account",
		"shopee_pay",
		"no_payment",
	}
}

const (
	PaymentBankAccount PaymentType = "bank_account"
	PaymentShopeePay   PaymentType = "shopee_pay"
	PaymentNoPayment   PaymentType = "no_payment"
)
