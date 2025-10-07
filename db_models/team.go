package db_models

type TeamType string

const (
	RootTeamType      TeamType = "root"
	AdminTeamType     TeamType = "admin"
	WarehouseTeamType TeamType = "warehouse"
	SellingTeamType   TeamType = "selling"
)

func (o TeamType) DefaultField() string {
	return "selling"
}

func (o TeamType) EnumList() []string {
	return []string{
		"root",
		"admin",
		"warehouse",
		"selling",
	}
}

func (o TeamType) IsEmpty() bool {
	return o.String() == ""
}

func (o TeamType) String() string {
	return string(o)
}

type TeamStat struct {
	InvoiceUnpaid   float64 `json:"invoice_unpaid"`
	InvoiceNotFinal float64 `json:"invoice_not_final"`
}

type Team struct {
	ID                uint     `json:"id" gorm:"primarykey"`
	Type              TeamType `json:"type"`
	Name              string   `json:"name"`
	TeamCode          TeamCode `json:"team_code" gorm:"index:team_code_unique,unique" binding:"required"`
	Description       string   `json:"desc"`
	ProductLimitCount int      `json:"product_limit_count"`
	ProductCount      int      `json:"product_count"`
	*TeamStat

	Deleted  bool         `json:"deleted"`
	Feature  *TeamFeature `json:"feature"`
	TeamInfo *TeamInfo    `json:"team_info"`
}

func (t *Team) GetEntityID() string {
	return "team"
}

type TeamFeature struct {
	ID     uint `json:"id" gorm:"primarykey"`
	TeamID uint `json:"team_id"`

	RestockSubmission bool `json:"restock_submission"`
	ProductPriority   bool `json:"product_priority"`
	PreventOrder      bool `json:"prevent_order"`
}

// GetEntityID implements authorization.Entity.
func (t *TeamFeature) GetEntityID() string {
	return "team_feature"
}

type TeamInfo struct {
	ID     uint `json:"id" gorm:"primaryKey"`
	TeamID uint `json:"team_id"`

	ReturnWarehouseID *uint `json:"return_warehouse_id"` // untuk destinasi gudang return team
	ReturnUserID      *uint `json:"return_user_id"`      // untuk destinasi gudang return team

	// untuk info agar team lain bisa transfer
	ContactNumber     string `json:"contact_number"`
	BankType          string `json:"bank_type"`
	BankOwnerName     string `json:"bank_owner_name"`
	BankAccountNumber string `json:"bank_account_number"`

	ReturnWarehouse *Warehouse `json:"return_warehouse,omitempty"`
	ReturnUser      *User      `json:"return_user,omitempty"`
}
