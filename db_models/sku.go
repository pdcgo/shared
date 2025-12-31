package db_models

import "time"

type Sku struct {
	ID          SkuID `json:"id" gorm:"primarykey"`
	VariantID   uint  `json:"variant_id"`
	TeamID      uint  `json:"team_id"`
	ProductID   uint  `json:"product_id"`
	WarehouseID uint  `json:"gudang_id"`

	StockReady   int `json:"stock_ready" gorm:"index"`
	StockPending int `json:"stock_pending"`
	StockTotal   int `json:"stock_total"` // deprecated stock total

	NextPrice     float64 `json:"next_price"`
	IsBlacklisted bool    `json:"is_blacklisted"`

	// LastRestock  time.Time `json:"last_restock" gorm:"index"`
	// statistik
	LastInbound  time.Time `json:"last_inbound" gorm:"index"`
	LastOutbound time.Time `json:"last_outbound" gorm:"index"`

	Warehouse *Warehouse      `json:"warehouse,omitempty"`
	Team      *Team           `json:"team,omitempty"`
	Product   *Product        `json:"product,omitempty"`
	Variant   *VariationValue `json:"variant,omitempty"`
}

// GetEntityID implements authorization.Entity.
func (sku *Sku) GetEntityID() string {
	return "sku"
}

func (sku *Sku) CalculateID() (SkuID, error) {

	idnya, err := NewSkuID(&SkuData{
		WarehouseID: sku.WarehouseID,
		TeamID:      sku.TeamID,
		ProductID:   sku.ProductID,
		VariantID:   sku.VariantID,
	})

	sku.ID = idnya

	return idnya, err
}
