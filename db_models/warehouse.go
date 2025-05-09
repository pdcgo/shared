package db_models

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"
)

type Warehouse struct {
	ID uint `json:"id" gorm:"primarykey;autoIncrement:false"`
	*WarehouseStat

	Name        string  `json:"name"`
	IsFull      bool    `json:"is_full"`
	UseFixedFee bool    `json:"use_fixed_fee"`
	FeeFix      float64 `json:"basic_fee_fix"`
	FeePercent  float32 `json:"fee_percent"`
	MaxFee      float64 `json:"max_fee"`

	Desc    string `json:"desc"`
	Address string `json:"address"`

	OpenTime   *time.Time `json:"open_time"`
	CloseTime  *time.Time `json:"close_time"`
	CloseOrder *time.Time `json:"close_order"`

	IsClosed bool `json:"is_closed"`

	Created time.Time `json:"created"`
	Deleted bool      `json:"deleted" gorm:"index"`

	Racks []*Rack `json:"racks"`
}

func (w *Warehouse) GetWarehouseFee(tprice float64) (float64, error) {
	if w.UseFixedFee {
		if w.FeeFix == 0 {
			return 0, errors.New("fixed fee empty")
		}
		return w.FeeFix, nil
	}

	if w.FeePercent == 0 {
		return 0, nil
	}

	var fee float64 = float64(tprice) * float64(w.FeePercent)
	fee = fee * 0.01
	fee = math.Ceil(fee)
	fee = fee * 100
	if fee > w.MaxFee {
		return w.MaxFee, nil
	}

	return fee, nil
}

type WarehouseStat struct {
	RackCount    uint `json:"rack_count"`
	OrderCount   uint `json:"order_count"`
	Capacity     uint `json:"capacity"`
	MaxCapacity  uint `json:"max_capacity"`
	ProductCount uint `json:"product_count"`
}

// GetEntityID implements authorize.Entity.
func (w *Warehouse) GetEntityID() string {
	return "warehouse"
}

func (w *Warehouse) UnmarshalJSON(data []byte) error {
	type Alias Warehouse
	aux := &struct {
		OpenTime   string `json:"open_time"`
		CloseTime  string `json:"close_time"`
		CloseOrder string `json:"close_order"`
		*Alias
	}{
		Alias: (*Alias)(w),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	strToTime := func(t string) *time.Time {
		if t == "" {
			return nil
		}

		tim, err := time.Parse("03:04", t)
		if err != nil {
			return &time.Time{}
		}

		return &tim
	}

	w.OpenTime = strToTime(aux.OpenTime)
	w.CloseTime = strToTime(aux.CloseTime)
	w.CloseOrder = strToTime(aux.CloseOrder)

	return nil
}

func (w *Warehouse) MarshalJSON() ([]byte, error) {
	type Alias Warehouse
	aux := &struct {
		OpenTime   string `json:"open_time"`
		CloseTime  string `json:"close_time"`
		CloseOrder string `json:"close_order"`
		FeePercent int    `json:"fee_percent"`
		*Alias
	}{
		Alias: (*Alias)(w),
	}

	getTime := func(t *time.Time) string {
		if t == nil {
			return ""
		}

		h := t.Hour()
		hh := strconv.Itoa(h)
		if h < 10 {
			hh = "0" + hh
		}

		m := t.Minute()
		mm := strconv.Itoa(m)
		if m < 10 {
			mm = "0" + mm
		}

		return fmt.Sprintf("%s:%s", hh, mm)
	}

	aux.FeePercent = int(w.FeePercent * 100)
	aux.OpenTime = getTime(w.OpenTime)
	aux.CloseTime = getTime(w.CloseTime)
	aux.CloseOrder = getTime(w.CloseOrder)

	return json.Marshal(aux)
}

type Rack struct {
	ID          uint   `json:"id" gorm:"primarykey"`
	WarehouseID uint   `json:"warehouse_id"`
	Name        string `json:"name"`
	// UseTransit  bool      `json:"use_transit"`
	IsSystem  bool      `json:"is_system"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime:milli"`
	Deleted   bool      `json:"deleted" gorm:"index"`
	// Warehouse *Warehouse `json:"warehouse"`
}

type Placement struct {
	ID     uint  `json:"id" gorm:"primarykey"`
	RackID uint  `json:"rack_id" gorm:"index:sku_rack_id,unique"`
	SkuID  SkuID `json:"sku_id" gorm:"index:sku_rack_id,unique"`
	Count  int   `json:"count"`

	Rack *Rack `json:"rack"`
	Sku  *Sku  `json:"sku"`
}

func (Placement) GetEntityID() string {
	return "placement"
}

// bakalan remove
type PlacementHistory struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	PlacementID uint      `json:"placement_id"`
	SkuID       SkuID     `json:"sku_id"`
	ByUserID    uint      `json:"by_user_id"`
	InvTxID     uint      `json:"inv_tx_id"`
	Type        InvTxType `json:"type"`

	Count int `json:"count"`

	InvTx     *InvTxItem `json:"inv_tx"`
	ByUser    *User      `json:"user"`
	Sku       *Sku       `json:"sku"`
	Placement *Placement `json:"placement"`
	CreateAt  time.Time  `json:"create_at"`
}

type WarehouseProduct struct {
	ID          uint `json:"id" gorm:"primarykey"`
	WarehouseID uint `json:"warehouse_id" gorm:"index:ware_product,unique"`
	ProductID   uint `json:"product_id" gorm:"index:ware_product,unique"`
	Stock       uint `json:"stock"`

	Product   *Product   `json:"product"`
	Warehouse *Warehouse `json:"warehouse"`
}
