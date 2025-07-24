package db_models

import (
	"errors"
	"sort"
	"time"

	"github.com/pdcgo/shared/interfaces/identity_iface"
	"gorm.io/datatypes"
)

type TransactionType string

const (
	TransactionTypeInbound  TransactionType = "inbound"
	TransactionTypeOutbound TransactionType = "outbound"
)

func (TransactionType) EnumList() []string {
	return []string{
		"inbound",
		"outbound",
	}
}

var TransactionTypeMap = map[TransactionType][]InvTxType{
	TransactionTypeInbound:  {InvTxRestock, InvTxAdjRestock, InvTxReturn, InvTxTransferIn},
	TransactionTypeOutbound: {InvTxOrder, InvTxTransferOut},
}

type InvTxType string

var outMap = map[InvTxType]bool{
	InvTxOrder:       true,
	InvTxTransferOut: true,
}
var inMap = map[InvTxType]bool{
	InvTxRestock:    true,
	InvTxAdjRestock: true,
	InvTxReturn:     true,
	InvTxTransferIn: true,
}

func (d InvTxType) IsInbound() bool {
	return inMap[d]
}

func (d InvTxType) IsOutbound() bool {
	return outMap[d]
}

func (InvTxType) EnumList() []string {
	return []string{
		"restock",
		"adj_restock",
		"order",
		"return",
		"transfer_out",
		"transfer_in",
		"transit",
		"broken",
		"ch_sku_out",
		"ch_sku_in",
		"adj_in",
		"adj_out",
	}
}

const (
	InvTxRestock      InvTxType = "restock"
	InvTxAdjRestock   InvTxType = "adj_restock"
	InvTxOrder        InvTxType = "order"
	InvTxReturn       InvTxType = "return"
	InvTxTransferIn   InvTxType = "transfer_in"
	InvTxTransferOut  InvTxType = "transfer_out"
	InvTxTransit      InvTxType = "transit"
	InvTxBroken       InvTxType = "broken"
	InvTxChangeSkuOut InvTxType = "ch_sku_out"
	InvTxChangeSkuIn  InvTxType = "ch_sku_in"
	InvTxAdjIn        InvTxType = "adj_in"
	InvTxAdjout       InvTxType = "adj_out"
	InvTxSysErrIn     InvTxType = "sys_err_in"
	InvTxSysErrOut    InvTxType = "sys_err_out"
)

type InvTxStatus string

func (InvTxStatus) EnumList() []string {
	return []string{
		"waiting",
		"ongoing",
		"cancel",
		"completed",
		"picking",
		"picked",
		"packing",
		"packing_completed",
	}
}

const (
	InvWaiting  InvTxStatus = "waiting"
	InvTxCancel InvTxStatus = "cancel"
	// pada kondisi di gudang
	InvTxOngoing   InvTxStatus = "ongoing"
	InvTxCompleted InvTxStatus = "completed"

	// untuk order dan return
	InvTxProductPick      InvTxStatus = "picking"           // barang proses pengambilan
	InvTxProductPicked    InvTxStatus = "picked"            // barang sudah diambil
	InvTxReadyForPacking  InvTxStatus = "packing"           // siap di packing
	InvTxReadyForCourrier InvTxStatus = "packing_completed" // siap diambil kurir

)

type InvTransaction struct {
	ID          uint  `json:"id" gorm:"primarykey"`
	TeamID      uint  `json:"team_id"`
	WarehouseID uint  `json:"warehouse_id"`
	CreateByID  uint  `json:"create_by_id"` // field bakalan deprecated and removed
	VerifyByID  *uint `json:"verify_by_id"`
	ShippingID  *uint `json:"shipping_id"`

	ExternOrdID string `json:"extern_ord_id" gorm:"index"`

	IsBroken        bool `json:"is_broken"`
	IsBrokenPartial bool `json:"is_broken_partial"`

	Receipt     string `json:"receipt" gorm:"index"`
	ReceiptFile string `json:"receipt_file"`

	Type      InvTxType   `json:"type"`
	Status    InvTxStatus `json:"status"`
	IsShipped bool        `json:"is_shipped"`
	Deleted   bool        `json:"deleted" gorm:"index"`

	Arrived *time.Time `json:"arrived"`
	SendAt  *time.Time `json:"send_at"`
	Created time.Time  `json:"created" gorm:"index"`

	ShippingFee float64 `json:"-"`
	// Deprecated: Shipping fee tidak digunakan, gunakan dari restock cost
	OtherFee float64 `json:"-"`

	Total float64 `json:"total"`

	Shipping  *Shipping      `json:"shipping"`
	Warehouse *Warehouse     `json:"warehouse"`
	Team      *Team          `json:"team"`
	CreatedBy *User          `json:"created_by" gorm:"foreignkey:CreateByID"`
	VerifyBy  *User          `json:"verify_by" gorm:"foreignkey:VerifyByID"`
	Items     InvItemList    `json:"items"`
	InvNotes  []*InvNote     `json:"notes"`
	Cost      []*RestockCost `json:"cost"`
}

func (i *InvTransaction) GetItemCount() int {
	pieces := 0
	for _, item := range i.Items {
		pieces += item.Count
	}
	return pieces
}

func (i *InvTransaction) GetEntityID() string {
	return "inv_transaction"
}

type InvTxItem struct {
	ID               uint    `json:"id" gorm:"primarykey"`
	InvTransactionID uint    `json:"tx_id"`
	SkuID            SkuID   `json:"sku_id"`
	Owned            bool    `json:"owned"`
	Count            int     `json:"count"`
	Price            float64 `json:"price"`
	Total            float64 `json:"total"`

	Sku            *Sku            `json:"sku"`
	InvTransaction *InvTransaction `json:"-"`
}

type InvItemList []*InvTxItem

func (i InvItemList) ProductIDs() ([]uint, error) {
	var err error
	mapprod := map[uint]bool{}
	result := []uint{}

	for _, item := range i {
		skuD, err := item.SkuID.Extract()
		if err != nil {
			return result, err
		}
		mapprod[skuD.ProductID] = true
	}
	for pid := range mapprod {
		result = append(result, pid)
	}

	return result, err
}

func (i InvItemList) TotalCount() (total int) {
	for _, item := range i {
		total += item.Count
	}

	return total
}

func (i InvItemList) Total() (total float64) {
	for _, item := range i {
		total += item.Price * float64(item.Count)
	}

	return total
}

func (i InvItemList) GetSkuCount() int {

	hasil := map[SkuID]bool{}
	for _, item := range i {
		hasil[item.SkuID] = true
	}

	return len(hasil)
}

func (i InvItemList) ItemGroupByTeam() (map[uint][]*InvTxItem, error) {
	hasil := map[uint][]*InvTxItem{}
	// hasil[i.TeamID] = []*InvTxItem{}

	for _, d := range i {
		item := d
		skuData, err := item.SkuID.Extract()
		if err != nil {
			return hasil, err
		}

		if hasil[skuData.TeamID] == nil {
			hasil[skuData.TeamID] = []*InvTxItem{}
		}

		hasil[skuData.TeamID] = append(hasil[skuData.TeamID], item)
	}

	return hasil, nil
}

type PartialItem struct {
	ProductID   uint `json:"product_id"`
	VariationID uint `json:"variation_id"`
	Count       int  `json:"count"`
}

func (ilist InvItemList) GetItemPartials(items []*PartialItem) (InvItemList, error) {
	var err error
	partialMap := map[uint]int{} // mapping dari map[variant_id]count
	itemMap := map[uint]int{}

	hasil := InvItemList{}

	for _, item := range items {
		partialMap[item.VariationID] += item.Count
	}

	for _, item := range ilist {
		skuData, err := item.SkuID.Extract()
		if err != nil {
			return hasil, err
		}
		itemMap[skuData.VariantID] += item.Count
	}

	for _, dd := range ilist {
		invitem := dd
		skuData, err := invitem.SkuID.Extract()
		if err != nil {
			return hasil, err
		}

		// filtering data
		if partialMap[skuData.VariantID] == 0 {
			continue
		}
		if partialMap[skuData.VariantID] > itemMap[skuData.VariantID] {
			return hasil, errors.New("partial count greater")
		}

		count := 0
		if partialMap[skuData.VariantID] > invitem.Count {
			count = invitem.Count
			partialMap[skuData.VariantID] -= invitem.Count
		} else {
			count = partialMap[skuData.VariantID]
			partialMap[skuData.VariantID] = 0
		}

		// log.Println(invitem.SkuID, count, partialMap[skuData.VariantID])

		citem := InvTxItem{
			// ID:               invitem.ID,
			// InvTransactionID: invitem.InvTransactionID,
			SkuID: invitem.SkuID,
			Count: count,
			Price: invitem.Price,
			Total: invitem.Price * float64(count),
		}
		citem.Count = count

		hasil = append(hasil, &citem)

	}

	return hasil, err
}

type NoteType string

func (NoteType) EnumList() []string {
	return []string{
		"common",
		"note_return",
		"note_cancel",
		"broken",
		"problem",
	}
}

const (
	NoteProblem NoteType = "problem"
	NoteCommon  NoteType = "common"
	NoteBroken  NoteType = "broken"
	NoteReturn  NoteType = "note_return"
	NoteCancel  NoteType = "note_cancel"
)

type InvNote struct {
	ID               uint  `json:"id" gorm:"primarykey"`
	InvTransactionID uint  `json:"tx_id"`
	OrderID          *uint `json:"order_id"`

	NoteType NoteType `json:"note_type"`
	NoteText string   `json:"note_text"`
}

type InvertoryHistory struct {
	ID uint `json:"id" gorm:"primarykey"`

	RackID      uint  `json:"rack_id"`
	TxID        *uint `json:"tx_id"`
	InTxID      *uint `json:"in_tx_id"`
	SkuID       SkuID `json:"sku_id"`
	WarehouseID uint  `json:"warehouse_id"`
	TeamID      uint  `json:"team_id"`
	UserID      uint  `json:"user_id"`

	Count    int       `json:"count"`
	Price    float64   `json:"price"`
	ExtPrice float64   `json:"ext_price"`
	Created  time.Time `json:"created" gorm:"index"`

	Rack      *Rack           `json:"rack"`
	Tx        *InvTransaction `json:"-" gorm:"foreignkey:TxID"`
	InTx      *InvTransaction `json:"-" gorm:"foreignkey:InTxID"`
	Sku       *Sku            `json:"-"`
	Warehouse *Warehouse      `json:"-"`
	Team      *Team           `json:"-"`
	User      *User           `json:"-"`
}

func (i *InvertoryHistory) GetFullPrice() float64 {
	return i.Price + i.ExtPrice
}

// GetEntityID implements authorization.Entity.
func (i *InvertoryHistory) GetEntityID() string {
	return "invertory_history"
}

func InvHistorySort(histories []*InvertoryHistory) []*InvertoryHistory {
	hasil := make([]*InvertoryHistory, len(histories))
	mapHist := map[uint][]*InvertoryHistory{}
	mapCount := map[uint]int{}

	for _, dd := range histories {
		item := dd

		if mapHist[item.RackID] == nil { // group by rack
			mapHist[item.RackID] = []*InvertoryHistory{}
		}

		mapHist[item.RackID] = append(mapHist[item.RackID], item)
		mapCount[item.RackID] += item.Count
	}

	listRack := []uint{}
	for rackID, histlist := range mapHist {

		sort.Slice(histlist, func(i, j int) bool {
			return histlist[i].Count > histlist[j].Count
		})

		listRack = append(listRack, rackID)
	}

	sort.Slice(listRack, func(i, j int) bool {
		return mapCount[listRack[i]] > mapCount[listRack[j]]
	})

	c := 0
	for _, rackID := range listRack {
		histlist := mapHist[rackID]

		for _, dd := range histlist {
			item := dd
			hasil[c] = item
			c += 1
		}
	}

	return hasil

}

type ActionType string

const (
	ActionEmpty        ActionType = ""
	ActionEditPrice    ActionType = "edit_price"
	ActionChangeStatus ActionType = "change_status"
)

type InvTimestamp struct {
	ID     uint `json:"id" gorm:"primarykey"`
	TxID   uint `json:"tx_id"`
	UserID uint `json:"user_id"`

	ActionType    ActionType               `json:"action_type"`
	Status        InvTxStatus              `json:"status"`
	Timestamp     time.Time                `json:"timestamp" gorm:"index"`
	From          identity_iface.AgentType `json:"from"`
	BeforeUpdated datatypes.JSONMap        `json:"before_updated"`
	Tx            *InvTransaction          `json:"-"`
	User          *User                    `json:"-"`
}
