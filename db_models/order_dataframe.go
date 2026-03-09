package db_models

import (
	"time"

	"github.com/wargasipil/data_processing/types"
)

type InvoItemD struct {
	MpFrom          types.Series[OrderMpType]
	ExternalOrderID types.Series[string]
	Type            types.Series[AdjustmentType]
	TransactionDate types.Series[time.Time]
	Description     types.Series[string]
	Amount          types.Series[float64]
	BalanceAfter    types.Series[float64]
	Region          types.Series[string]
	IsOtherRegion   types.Series[bool]
	Failed          types.Series[bool]
}

type InvoItemDataFrame struct {
	count  int
	offset []int
	D      *InvoItemD
}

func NewInvoItemDataFrame(datas []*InvoItem) *InvoItemDataFrame {
	d := &InvoItemDataFrame{
		count:  0,
		offset: []int{},
		D: &InvoItemD{
			MpFrom:          types.Series[OrderMpType]{},
			ExternalOrderID: types.Series[string]{},
			Type:            types.Series[AdjustmentType]{},
			TransactionDate: types.Series[time.Time]{},
			Description:     types.Series[string]{},
			Amount:          types.Series[float64]{},
			BalanceAfter:    types.Series[float64]{},
			Region:          types.Series[string]{},
			IsOtherRegion:   types.Series[bool]{},
			Failed:          types.Series[bool]{},
		},
	}

	for _, item := range datas {
		d.D.MpFrom = append(d.D.MpFrom, item.MpFrom)
		d.D.ExternalOrderID = append(d.D.ExternalOrderID, item.ExternalOrderID)
		d.D.Type = append(d.D.Type, item.Type)
		d.D.TransactionDate = append(d.D.TransactionDate, item.TransactionDate)
		d.D.Description = append(d.D.Description, item.Description)
		d.D.Amount = append(d.D.Amount, item.Amount)
		d.D.BalanceAfter = append(d.D.BalanceAfter, item.BalanceAfter)
		d.D.Region = append(d.D.Region, item.Region)
		d.D.IsOtherRegion = append(d.D.IsOtherRegion, item.IsOtherRegion)
		d.D.Failed = append(d.D.Failed, item.Failed)
		d.offset = append(d.offset, d.count)
		d.count++
	}

	return d
}

func (d *InvoItemDataFrame) Query(filters ...types.OffsetFilter) *InvoItemDataFrame {
	offset := d.offset
	for _, filter := range filters {
		offset = filter(offset)
	}

	newdf := InvoItemDataFrame{
		count:  len(offset),
		D:      d.D,
		offset: offset,
	}

	return &newdf
}

func (d *InvoItemDataFrame) get(i int) *InvoItem {
	item := InvoItem{
		MpFrom:          d.D.MpFrom[i],
		ExternalOrderID: d.D.ExternalOrderID[i],
		Type:            d.D.Type[i],
		TransactionDate: d.D.TransactionDate[i],
		Description:     d.D.Description[i],
		Amount:          d.D.Amount[i],
		BalanceAfter:    d.D.BalanceAfter[i],
		Region:          d.D.Region[i],
		IsOtherRegion:   d.D.IsOtherRegion[i],
		Failed:          d.D.Failed[i],
	}

	return &item
}

func (d *InvoItemDataFrame) First() *InvoItem {
	if len(d.offset) == 0 {
		return nil
	}

	i := d.offset[0]
	return d.get(i)
}

func (d *InvoItemDataFrame) Last() *InvoItem {
	if len(d.offset) == 0 {
		return nil
	}

	i := d.offset[len(d.offset)-1]
	return d.get(i)
}

func (d *InvoItemDataFrame) Data() []*InvoItem {
	result := []*InvoItem{}
	for _, i := range d.offset {
		result = append(result, d.get(i))
	}

	return result
}
