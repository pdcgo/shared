package db_models

import "gorm.io/gorm"

type V2Supplier struct {
	ID          uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
	TeamID      uint64         `json:"team_id" gorm:"not null;index"`
	Code        string         `json:"code" gorm:"uniqueIndex:uidx_code_active,where:deleted_at IS NULL"`
	Name        string         `json:"name"`
	Contact     string         `json:"contact,omitempty"`
	Province    string         `json:"province,omitempty"`
	City        string         `json:"city,omitempty"`
	Description string         `json:"description,omitempty"`
	Address     string         `json:"address,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type V2SupplierMarketplace struct {
	SupplierID  uint64         `gorm:"not null"`
	ID          uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
	MpType      int32          `gorm:"not null"`
	ShopName    string         `gorm:"not null;size:200;default:''"`
	ProductName string         `gorm:"not null;size:250;default:''"`
	URI         string         `gorm:"not null;size:500;default:''"`
	Description string         `gorm:"not null;size:500;default:''"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type V2SupplierInvTxItem struct {
	ID          uint `json:"id" gorm:"primarykey"`
	InvTxItemID uint `json:"inv_tx_item_id"`
	SupplierID  uint `json:"supplier_id"`

	Supplier  *V2SupplierMarketplace `json:"supplier"`
	InvTxItem *InvTxItem             `json:"inv_tx_item"`
}
