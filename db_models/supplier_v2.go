package db_models

import "gorm.io/gorm"

type SupplierV2 struct {
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

type SupplierMarketplaceV2 struct {
	SupplierID  uint64         `gorm:"not null"`
	ID          uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
	MpType      int32          `gorm:"not null"`
	ShopName    string         `gorm:"not null;size:200;default:''"`
	ProductName string         `gorm:"not null;size:250;default:''"`
	URI         string         `gorm:"not null;size:500;default:''"`
	Description string         `gorm:"not null;size:500;default:''"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
