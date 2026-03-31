package db_models

import (
	"time"

	"github.com/pdcgo/schema/services/selling_iface/v1"
	"gorm.io/gorm"
)

type Supplier struct {
	ID        uint64                     `gorm:"primaryKey;autoIncrement"`
	TeamID    uint64                     `gorm:"not null;index"`
	Type      selling_iface.SupplierType `gorm:"not null"`
	DeletedAt gorm.DeletedAt             `gorm:"index"`
}

type SupplierCustom struct {
	SupplierID  uint64         `gorm:"primaryKey"`
	Name        string         `gorm:"not null;size:200"`
	Contact     string         `gorm:"not null;size:50"`
	Description string         `gorm:"not null;size:500"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type SupplierMarketplace struct {
	SupplierID  uint64         `gorm:"primaryKey"`
	MpType      int32          `gorm:"not null"`
	ShopName    string         `gorm:"not null;size:200;default:''"`
	ProductName string         `gorm:"not null;size:200;default:''"`
	URI         string         `gorm:"not null;size:500"`
	Description string         `gorm:"not null;size:500"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type VariantSupplierV2 struct {
	ID           uint          `json:"id" gorm:"primarykey"`
	TeamID       uint          `json:"team_id"`
	VariantID    uint          `json:"variant_id"`
	SupplierID   uint          `json:"supplier_id"`
	PreOrderTime time.Duration `json:"pre_order_time"`

	Team     *Team           `json:"team,omitempty"`
	Variant  *VariationValue `json:"variant,omitempty"`
	Supplier *Supplier       `json:"supplier,omitempty"`
}

type SupplierInvTxItemV2 struct {
	ID          uint `json:"id" gorm:"primarykey"`
	InvTxItemID uint `json:"inv_tx_item_id"`
	SupplierID  uint `json:"supplier_id"`

	Supplier  *Supplier  `json:"supplier"`
	InvTxItem *InvTxItem `json:"inv_tx_item"`
}
