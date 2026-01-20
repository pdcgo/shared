package db_models

// type DraftOrder struct { // bakalan ada di database
// 	ID uint `json:"id"`

// 	TeamID uint `json:"team_id"`
// 	UserID uint `json:"user_id"`

// 	OrderRefID   string                              `json:"order_ref_id" gorm:"index"`
// 	OrderMpID    uint                                `json:"order_mp_id"`
// 	OrderTotal   int                                 `json:"order_total"`
// 	OrderFrom    OrderMpType                         `json:"order_from"`
// 	OrderPayload database.JSONType[*OrderPayload]    `json:"order_payload"`
// 	MpProducts   datatypes.JSONSlice[*MpProductItem] `json:"mp_products"`
// 	Created      time.Time                           `json:"created"`

// 	OrderMp *Marketplace `json:"order_mp"`
// 	Team    *Team        `json:"team"`
// 	User    *User        `json:"user"`
// }

// func (d *DraftOrder) GetEntityID() string {
// 	return "draft_order"
// }
