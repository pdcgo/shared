package db_models

type Tag struct {
	ID   uint   `json:"id" gorm:"primarykey"`
	Name string `json:"name" gorm:"index:tag_name_unique,unique" binding:"required,lte=100"`
}

type ProductTag struct {
	ProductID uint `gorm:"primaryKey"`
	TagID     uint `gorm:"primaryKey"`

	Product *Product `json:"product"`
	Tag     *Tag     `json:"tag"`
}
