package db_models

import (
	"errors"
	"math"
	"time"

	"github.com/pdcgo/shared/interfaces/identity_iface"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type MarkupValue float64

func (mark MarkupValue) Markup(price float64) float64 {
	// pp := (float64(100) * price) / (100 - float64(mark))
	pp := (float64(price) * 100.0) / (100.0 - float64(mark))
	return math.Ceil(pp)
}

type Product struct {
	ID      uint  `json:"id" gorm:"primarykey"`
	TeamID  uint  `json:"team_id"`
	Deleted bool  `json:"deleted" gorm:"index"`
	UserID  uint  `json:"user_id"`
	RefID   RefID `json:"ref_id"`

	Name            string                      `json:"name"`
	Image           datatypes.JSONSlice[string] `json:"image"`
	Desc            string                      `json:"desc"`
	VariationNames  []*VariationName            `json:"variation_names"`
	VariationValues []*VariationValue           `json:"variation_values"`

	MarkupPercent MarkupValue `json:"markup_percent"`

	StockReady    int       `json:"stock_ready"`
	StockPending  int       `json:"stock_pending"`
	BundleCount   int       `json:"bundle_count"`
	Priority      bool      `json:"priority"`
	CrossLocked   bool      `json:"cross_locked"`
	StockReserved int       `json:"stock_reserved"`
	Created       time.Time `json:"created"`

	Team     *Team       `json:"team,omitempty"`
	User     *User       `json:"user,omitempty"`
	Category []*Category `json:"category,omitempty" gorm:"many2many:product_category;"`
	Tags     []*Tag      `json:"tags,omitempty" gorm:"many2many:product_tags;"`
}

func (prod *Product) GetUserCode() (string, error) {

	data, err := prod.RefID.ExtractData()
	return data.UserCode, err
}

func (prod *Product) GetImageIDs(db *gorm.DB) ([]string, error) {
	ids := []string{}
	ids = append(ids, prod.Image...)

	idmap := map[string]bool{}

	for _, varival := range prod.VariationValues {
		image := varival.Image
		if image == "" {
			continue
		}

		if idmap[image] {
			continue
		}

		idmap[image] = true
		ids = append(ids, image)

	}

	return ids, nil
}

// GetEntityID implements authorize.Entity.
func (prod *Product) GetEntityID() string {
	return "product"
}

func (prod *Product) VarNameExist(id uint) bool {
	for _, vname := range prod.VariationNames {
		if vname.ID == id {
			return true
		}
	}

	return false
}

func (prod *Product) VarValueExist(id uint) *VariationValue {
	for _, vname := range prod.VariationValues {
		if vname.ID == id {
			return vname
		}
	}

	return nil
}

type ProductLogType string

const (
	ProductLogDelete   ProductLogType = "delete"
	ProductLogUndelete ProductLogType = "undelete"
)

type ProductLog struct {
	ID        uint                     `json:"id" gorm:"primarykey"`
	UserID    uint                     `json:"user_id"`
	ProductID uint                     `json:"product_id"`
	From      identity_iface.AgentType `json:"from"`
	Type      ProductLogType           `json:"type"`
	At        time.Time                `json:"at"`
}

type VariationName struct {
	ID        uint                        `json:"id" gorm:"primarykey"`
	ProductID uint                        `json:"product_id"`
	Name      string                      `json:"name" gorm:"index"`
	Options   datatypes.JSONSlice[string] `json:"options"`

	Product *Product `json:"-"`
}

func (vary *VariationName) CheckOptionExist(value string) bool {
	for _, nn := range vary.Options {
		if nn == value {
			return true
		}
	}
	return false
}

type VariationValue struct {
	ID             uint                        `json:"id" gorm:"primarykey"`
	RefID          RefID                       `json:"ref_id"`
	ProductID      uint                        `json:"product_id"`
	Image          string                      `json:"image"`
	VariationName  datatypes.JSONSlice[string] `json:"variation_name"`
	VariationValue datatypes.JSONSlice[string] `json:"variation_value"`

	PriceMin float64 `json:"price_min"`
	PriceMax float64 `json:"price_max"`

	HasSku      bool       `json:"has_sku"`
	LastOrdered *time.Time `json:"last_ordered"`
	IsDefault   bool       `json:"is_default"`

	StockReady   int `json:"stock_ready"`
	StockPending int `json:"stock_pending"`
	BundleCount  int `json:"bundle_count"`

	Product *Product `json:"product,omitempty"`
}

func (vv *VariationValue) ReturnID(db *gorm.DB, teamID uint) (uint, error) {
	var toReturnID uint

	if vv.ID == 0 {
		return toReturnID, errors.New("variation not loaded maybe")
	}

	err := db.
		Select("to_var_id").
		Model(&VariationReturnMap{}).
		Where("team_id = ?", teamID).
		Where("var_id = ?", vv.ID).
		Find(&toReturnID).
		Error

	if toReturnID == 0 {
		return toReturnID, errors.New("variation not mapped")
	}

	return toReturnID, err
}

func (vv *VariationValue) AddMapReturnID(db *gorm.DB, retVarID, teamID uint) error {
	if vv.ID == 0 {
		return errors.New("variation not loaded maybe")
	}
	maper := VariationReturnMap{
		ToVarID:  retVarID,
		ToTeamID: teamID,
		VarID:    vv.ID,
	}

	err := db.
		Model(&VariationReturnMap{}).
		Where("var_id = ?", vv.ID).
		Where("to_team_id = ?", teamID).
		Find(&maper).
		Error

	if err != nil {
		return err
	}

	err = db.Save(&maper).Error
	return err
}

type VariationValues []*VariationValue

type VariationReturnMap struct {
	ToVarID  uint // milik team return
	ToTeamID uint

	VarID uint // milik sendiri

	ToVar *VariationValue `gorm:"foreignKey:ParentVarID;"`
	Var   *VariationValue `gorm:"foreignKey:VarID;"`
	// Team  *Team
}
