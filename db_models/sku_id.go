package db_models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type SkuID string

func (id SkuID) String() string {
	return string(id)
}

type SkuData struct {
	WarehouseID uint `json:"warehouse_id"`
	TeamID      uint `json:"team_id"`
	ProductID   uint `json:"product_id"`
	VariantID   uint `json:"variant_id"`
}

func (skuID SkuID) Extract() (*SkuData, error) {
	hasil := SkuData{}

	header := skuID[:4]
	data := skuID[4:]

	stlen := 0
	for idx, elem := range header {
		lendata, err := hexToUint(string(elem))
		if err != nil {
			return nil, err
		}
		raw := data[stlen : stlen+int(lendata)]
		stlen += int(lendata)

		val, err := hexToUint(string(raw))
		if err != nil {
			return nil, err
		}

		switch idx {
		case 0:
			hasil.WarehouseID = val
		case 1:
			hasil.TeamID = val
		case 2:
			hasil.ProductID = val
		case 3:
			hasil.VariantID = val
		}

	}

	return &hasil, nil
}

func NewSkuID(data *SkuData) (SkuID, error) {
	ceks := []uint{data.WarehouseID, data.TeamID, data.ProductID, data.VariantID}
	var hasil []string = make([]string, 5)
	var header string
	for c, cek := range ceks {
		if cek == 0 {
			return "", errors.New("sku incomplete information for id generation")
		}
		index := c + 1
		segment := fmt.Sprintf("%x", cek)
		header += fmt.Sprintf("%x", len(segment))
		hasil[index] = segment

	}

	hasil[0] = header
	idnya := strings.Join(hasil, "")

	return SkuID(idnya), nil
}

func hexToUint(data string) (uint, error) {
	val, err := strconv.ParseUint(data, 16, 64)

	return uint(val), err
}
