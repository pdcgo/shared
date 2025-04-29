package db_models

type Category struct {
	ID       uint        `json:"id" gorm:"primarykey"`
	ParentID *uint       `json:"parent_id" gorm:"index"`
	Level    int         `json:"level"`
	Name     string      `json:"name"`
	IsLast   bool        `json:"is_last"`
	Children []*Category `json:"children" gorm:"foreignKey:ParentID"`
}

type CategoryStat struct {
	ID             uint      `json:"id" gorm:"primarykey"`
	CategoryID     uint      `json:"category_id"`
	ProductCount   int       `json:"product_count"`
	OrderCount     int       `json:"order_count"`
	OrderItemCount int       `json:"order_item_count"`
	Category       *Category `json:"category"`
}
