package db_models

type ExpenseTypeAccount string

const (
	BankTypeAccount   ExpenseTypeAccount = "bank"
	ShopeeTypeAccount ExpenseTypeAccount = "shopeepay"
)

func (ExpenseTypeAccount) EnumList() []string {
	return []string{
		"bank",
		"shopeepay",
	}
}

type AccountType struct {
	ID   uint               `json:"id" gorm:"primarykey"`
	Key  string             `json:"key" gorm:"index:account_type_unique,unique"`
	Name string             `json:"name"`
	Type ExpenseTypeAccount `json:"type"`
}

func (*AccountType) GetEntityID() string {
	return "account_type"
}
