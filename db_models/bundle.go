package db_models

import "time"

type Bundle struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	UserID    uint      `json:"user_id"`
	TeamID    uint      `json:"team_id"`
	RefID     RefID     `json:"ref_id"`
	Name      string    `json:"name"`
	PairCount uint      `json:"pair_count"`
	PriceMin  float64   `json:"price_min"`
	PriceMax  float64   `json:"price_max"`
	CreateAt  time.Time `json:"create_at"`

	Items []*BundleItem `json:"items"`
	Team  *Team         `json:"team"`
	User  *User         `json:"user"`
}

func (b *Bundle) GetEntityID() string {
	return "bundle"
}

type BundleItem struct {
	ID        uint `json:"id" gorm:"primarykey"`
	BundleID  uint `json:"bundle_id"`
	VariantID uint `json:"variant_id"`
	ProductID uint `json:"product_id"`
	Count     uint `json:"count"`

	Product *Product        `json:"product,omitempty"`
	Variant *VariationValue `json:"variant,omitempty"`
}

type OrderBundle struct {
	ID       uint `json:"id" gorm:"primarykey"`
	BundleID uint `json:"bundle_id"`
	OrderID  uint `json:"order_id"`

	Bundle *Bundle `json:"bundle"`
}

type WarehouseBundle struct {
	ID          uint `json:"id" gorm:"primarykey"`
	BundleID    uint `json:"bundle_id"`
	WarehouseID uint `json:"warehouse_id"`
	PairCount   uint `json:"pair_count"`

	Bundle    *Bundle    `json:"bundle"`
	Warehouse *Warehouse `json:"warehouse"`
}
