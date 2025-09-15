package db_models

type MarketplaceType string

func (m MarketplaceType) EnumList() []string {
	return []string{
		"tokopedia",
		"shopee",
		"tiktok",
		"lazada",
		"mengantar",
		"custom",
	}
}

const (
	MpShopee    MarketplaceType = "shopee"
	MpTokopedia MarketplaceType = "tokopedia"
	MpTiktok    MarketplaceType = "tiktok"
	MpMengantar MarketplaceType = "mengantar"
	MpCustom    MarketplaceType = "custom"
	MpLazada    MarketplaceType = "lazada"
)

type Marketplace struct {
	ID            uint  `json:"id" gorm:"primarykey"`
	TeamID        uint  `json:"team_id" gorm:"index:mp_team_id_unique,unique"`
	HoldAssetID   *uint `json:"asset_id"`
	BankAccountID *uint `json:"bank_account_id"`

	MpUsername  string          `json:"mp_username" gorm:"index:mp_team_id_unique,unique"`
	MpName      string          `json:"mp_name"`
	MpType      MarketplaceType `json:"mp_type" gorm:"index:mp_team_id_unique,unique"`
	Uri         string          `json:"uri"`
	IsDuplicate bool            `json:"is_duplicate"`
	Deleted     bool            `json:"deleted" gorm:"index"`

	HoldAsset   *Asset       `json:"hold_asset"`
	BankAccount *BankAccount `json:"bank_account"`
	Team        *Team        `json:"team"`
}

func (m *Marketplace) GetEntityID() string {
	return "marketplace"
}

func (m *Marketplace) IsHaveHoldAsset() bool {
	if m.HoldAssetID == nil {
		return false
	}
	if *m.HoldAssetID == 0 {
		return false
	}

	return true
}

type UserMarketplace struct {
	ID            uint         `json:"id" gorm:"primarykey"`
	UserID        uint         `json:"user_id" gorm:"index:user_mp_unique,unique"`
	MarketplaceID uint         `json:"marketplace_id" gorm:"index:user_mp_unique,unique"`
	User          *User        `json:"user"`
	Marketplace   *Marketplace `json:"marketplace"`
}

type BankAccount struct {
	ID      uint   `json:"id" gorm:"primarykey"`
	AssetID uint   `json:"asset_id"`
	Asset   *Asset `json:"asset"`
}
