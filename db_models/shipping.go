package db_models

type Shipping struct {
	ID          uint   `gorm:"primarykey" json:"id"`
	Key         string `gorm:"index:shipping_key_unique,unique" json:"key"`
	DisplayName string `json:"display_name"`
}
