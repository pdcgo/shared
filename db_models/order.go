package db_models

import (
	"errors"
	"time"

	"github.com/pdcgo/shared/interfaces/identity_iface"
	"gorm.io/datatypes"
)

type AdjustmentType string

const (
	AdjCommision AdjustmentType = "aff_commission"
	// Penyesuaian Saldo Penjual untuk biaya premi Pesanan yang Gagal Terkirim: 241006SDVJR41U
	AdjPremi AdjustmentType = "premi"
	// Kompensasi Biaya Kemasan Program Garansi Bebas Pengembalian 2410311PXCW62F
	AdjPackaging AdjustmentType = "packaging"
	// [Penambahan Wallet] Pengembalian Dana dari Order Return 250106Q36VBD4V
	// Penyesuaian saldo dilakukan untuk pesanan 24103110QE0X03 karena terdapat Pengembalian Barang/Dana setelah dana dilepaskan
	AdjReturn AdjustmentType = "return_adj"
	// Penyesuaian Ongkos Kirim Bebas Pengembalian 241224MKUP9X4U
	AdjShipping     AdjustmentType = "shipping_adj"
	AdjCompensation AdjustmentType = "compensation"
	// [Penambahan Wallet] Penggantian Dana Penuh Barang Hilang 2412063KHB09UF
	// [Penambahan Wallet] Penggantian Dana Sebagian Barang Hilang 241214SG9CB1JB
	AdjLostCompensation AdjustmentType = "lost_compensation"
	AdjUnknown          AdjustmentType = "unknown"
	AdsPayment          AdjustmentType = "ads_payment"

	AdjFund       AdjustmentType = "fund" // jarang digunakan, untuk wd
	AdjOrderFund  AdjustmentType = "order_fund"
	AdjUnknownAdj AdjustmentType = "unknown_adj"
)

func (AdjustmentType) EnumList() []string {
	return []string{
		"aff_commission",
		"premi",
		"packaging",
		"return_adj",
		"shipping_adj",
		"compensation",
		"lost_compensation",
		"unknown",
		"order_fund",
		"unknown_adj",
	}
}

type OrdStatus string

// DefaultField implements api_model.SortingField.
func (o OrdStatus) DefaultField() string {
	return "created"
}

// EnumList implements api_model.SortingField.
func (o OrdStatus) EnumList() []string {
	return []string{
		"created",
		"process",
		"picking",
		"packing",
		"packing_completed",
		"shipped",
		"courrier_shipped",
		"completed",
		"cancel",
		"problem",
		"return",
		"return_problem",
		"return_completed",
	}
}

// IsEmpty implements api_model.SortingField.
func (o OrdStatus) IsEmpty() bool {
	return o.String() == ""
}

// String implements api_model.SortingField.
func (o OrdStatus) String() string {
	return string(o)
}

const (
	OrdCreated         OrdStatus = "created" // order ketika di request tapi belum diproses sama kang paket gudang
	OrdProcess         OrdStatus = "process" // order ketika di request tapi sudah diproses sama kang paket gudang
	OrdShipped         OrdStatus = "shipped" // order ketika di request tapi sudah dikirim sama kang paket gudang
	OrdCourrierShipped OrdStatus = "courrier_shipped"
	OrdCompleted       OrdStatus = "completed" // order ketika di request tapi sudah diproses sama kang paket gudang dan dibayar
	OrdCancel          OrdStatus = "cancel"
	OrdProblem         OrdStatus = "problem"
	OrdReturn          OrdStatus = "return"
	OrdReturnProblem   OrdStatus = "return_problem"
	OrdReturnCompleted OrdStatus = "return_completed"

	// status untuk proses packing di gudang
	OrdProductPick      OrdStatus = "picking"           // barang proses pengambilan
	OrdReadyForPacking  OrdStatus = "packing"           // siap di packing
	OrdReadyForCourrier OrdStatus = "packing_completed" // siap diambil kurir
)

type OrderMpType string

func (OrderMpType) EnumList() []string {
	return []string{
		"tokopedia",
		"shopee",
		"tiktok",
		"lazada",
		"custom",
		"mengantar",
	}
}

const (
	OrderMpTokopedia OrderMpType = "tokopedia"
	OrderMpShopee    OrderMpType = "shopee"
	OrderMpTiktok    OrderMpType = "tiktok"
	OrderMpLazada    OrderMpType = "lazada"
	OrderMengantar   OrderMpType = "mengantar"
	OrderMpCustom    OrderMpType = "custom"
)

func (data OrderMpType) Validate() error {
	listtype := []OrderMpType{
		OrderMpTokopedia,
		OrderMpShopee,
		OrderMpTiktok,
		OrderMpLazada,
		OrderMpCustom,
		OrderMengantar,
	}

	for _, tipe := range listtype {
		if tipe == data {
			return nil
		}
	}
	return errors.New("order marketplace type error")
}

type ProductSourceType string

const (
	SupplierProdSource  ProductSourceType = "supplier"
	WarehouseProdSource ProductSourceType = "warehouse"
	DummyProdSource     ProductSourceType = "dummy"
)

func (ProductSourceType) EnumList() []string {
	return []string{
		"supplier",
		"warehouse",
		"dummy",
	}
}

// untuk dummy sebelum order dipindah semua ke ware db

type CustomerAddress struct {
	ID         uint   `json:"id" gorm:"primarykey"`
	OrderID    uint   `json:"order_id"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Province   string `json:"province"`
	City       string `json:"city"`
	District   string `json:"district"`
	PostalCode string `json:"postal_code"`
	Address    string `json:"address"`

	Order *Order `json:"-"`
}

type OrderAdditionalCostType string

const (
	FakeOrderCost OrderAdditionalCostType = "fake_order_cost"
)

func (OrderAdditionalCostType) EnumList() []string {
	return []string{
		"fake_order_cost",
	}
}

type OrderAdditionalCost struct {
	Type        OrderAdditionalCostType `json:"type"`
	PaymentType PaymentType             `json:"payment_type"`
	Amount      float64                 `json:"amount"`
}

type Order struct {
	ID                  uint  `json:"id" gorm:"primarykey"`
	TeamID              uint  `json:"team_id"`
	CreatedByID         uint  `json:"created_by_id"`
	InvertoryTxID       *uint `json:"invertory_tx_id"`
	InvertoryReturnTxID *uint `json:"invertory_ret_tx_id"`
	DoubleOrder         bool  `json:"double_order"`

	OrderRefID   string      `json:"order_ref_id" gorm:"index"`
	OrderFrom    OrderMpType `json:"order_from"`
	OrderMpTotal int         `json:"order_mp_total"`

	// bagian tipe2 order
	ProductSourceType ProductSourceType `json:"product_source"`
	ParentPartialID   *uint             `json:"parent_partial_id"`
	IsPartial         bool              `json:"is_partial"`
	IsOrderFake       bool              `json:"is_order_fake"`

	// perkara wd fund
	WdTotal float64   `json:"wd_total"`
	WdTime  time.Time `json:"wd_time"`
	WdSet   bool      `json:"wd_set"`

	WdFundAt time.Time `json:"wd_fund_at"`
	WdFund   bool      `json:"wd_fund"`

	Adjustment float64 `json:"adjustment"`

	OrderTime time.Time `json:"order_time" gorm:"index"`
	OrderMpID uint      `json:"order_mp_id"`

	Receipt           string `json:"receipt"`
	ReceiptFile       string `json:"receipt_file"`
	ReceiptReturn     string `json:"receipt_return"`
	ReceiptReturnFile string `json:"receipt_return_file"`

	Status          OrdStatus                                 `json:"status"`
	WarehouseFee    float64                                   `json:"warehouse_fee"`
	ShipmentFee     float64                                   `json:"shipping_fee"`
	ItemCount       int                                       `json:"item_count"`
	Total           float64                                   `json:"total"`
	CreatedAt       time.Time                                 `json:"created_at" gorm:"index"`
	AdditionalCosts datatypes.JSONSlice[*OrderAdditionalCost] `json:"additional_costs"`

	Address           *CustomerAddress  `json:"address"`
	Invoices          []*Invoice        `json:"-"`
	Items             []*OrderItem      `json:"items"`
	InvertoryTx       *InvTransaction   `json:"invertory_tx" gorm:"foreignkey:InvertoryTxID"`
	InvertoryReturnTx *InvTransaction   `json:"invertory_return_tx" gorm:"foreignkey:InvertoryReturnTxID"`
	OrderMp           *Marketplace      `json:"order_mp,omitempty"`
	CreatedBy         *User             `json:"user,omitempty"`
	OrderBundle       []*OrderBundle    `json:"order_bundles,omitempty"`
	Notes             []*InvNote        `json:"notes"`
	OrderTimestamp    []*OrderTimestamp `json:"timestamp"`
	Tags              []*OrderTag       `json:"tags,omitempty" gorm:"many2many:order_tag_relations;"`
}

func (o *Order) HaveNoteText(text string) bool {
	var foundNote bool
	for _, ordn := range o.Notes {
		if ordn.NoteText == text {
			foundNote = true
		}
	}

	return foundNote
}

func (o *Order) IsWarehouseProcessed() bool {

	switch o.Status {
	case OrdCreated:
		return false
	default:
		return true
	}
}

func (o *Order) IsGiveToCourrier() bool {
	switch o.Status {
	case OrdShipped:
		return true
	}

	return false
}

func (o *Order) IsHaveSent() bool {
	switch o.Status {
	case OrdCompleted, OrdShipped, OrdCourrierShipped, OrdReturn, OrdReturnCompleted:
		return true
	}

	return false
}

func (o *Order) GetEntityID() string {
	return "order"
}

type OrderItem struct {
	ID             uint `json:"id" gorm:"primarykey"`
	OrderID        uint `json:"order_id"`
	ProductID      uint `json:"product_id"`
	VariationID    uint `json:"variation_id"`
	RetVariationID uint `json:"ret_variation_id"`

	Markup      MarkupValue `json:"markup"`
	Owned       bool        `json:"owned"`
	ProductName string      `json:"name"`
	Price       float64     `json:"price"`
	Count       int         `json:"count"`
	Total       float64     `json:"total"`

	Product   *Product        `json:"product"`
	Variation *VariationValue `json:"variation,omitempty"`
}

type OrderTimestamp struct {
	ID          uint                     `json:"id" gorm:"primarykey"`
	OrderID     uint                     `json:"order_id"`
	UserID      uint                     `json:"user_id"`
	From        identity_iface.AgentType `json:"from"`
	OrderStatus OrdStatus                `json:"order_status"`
	Timestamp   time.Time                `json:"timestamp" gorm:"index"`

	Order *Order `json:"-"`
	User  *User  `json:"-"`
}

type OrderTag struct {
	ID   uint   `json:"id" gorm:"primarykey"`
	Name string `json:"name" gorm:"index:order_tag_unique,unique" binding:"required,lte=100"`
}

type RelationFrom string

const (
	RelationFromWarehouse RelationFrom = "warehouse"
	RelationFromTracking  RelationFrom = "tracking"
	RelationFromUser      RelationFrom = "user"
	RelationFromUnknown   RelationFrom = ""
)

type OrderTagRelation struct {
	OrderID      uint   `json:"order_id" gorm:"primaryKey"`
	OrderTagID   uint   `json:"order_tag_id" gorm:"primaryKey"`
	RelationFrom string `json:"relation_from"`

	Order    *Order    `json:"product"`
	OrderTag *OrderTag `json:"tag"`
}

type InvoItem struct {
	MpFrom          OrderMpType
	ExternalOrderID string
	Type            AdjustmentType
	TransactionDate time.Time
	Description     string
	Amount          float64
	BalanceAfter    float64
}

type InvoItemList []*InvoItem

func (in InvoItemList) TotalAmount() float64 {
	var amount float64 = 0
	for _, item := range in {
		amount += item.Amount
	}

	return amount
}

type MarketplaceOrderItem struct {
	ID            uint   `json:"id" gorm:"primaryKey"`
	OrderID       uint   `json:"order_id"`
	MpProductName string `json:"mp_product_name"`
	Count         int    `json:"count"`

	Order *Order `json:"-" gorm:"foreignKey:OrderID"`
}
