package spx_tracker

type TrackingStatus string

const (
	Delivered TrackingStatus = "Delivered"
	Lost      TrackingStatus = "Lost"
	OnHold    TrackingStatus = "OnHold"
	Cancelled TrackingStatus = "Cancelled"
)

type TrackingListItem struct {
	Timestamp int64  `json:"timestamp"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

type StatusListItem struct {
	Timestamp int    `json:"timestamp"`
	Code      int    `json:"code"`
	Text      string `json:"text"`
	StateLs   string `json:"state_ls"`
	Icon      string `json:"icon"`
}

type TrackResponseData struct {
	SlsTrackingNumber string              `json:"sls_tracking_number"`
	NeedTranslate     int                 `json:"need_translate"`
	DeliveryType      string              `json:"delivery_type"`
	RecipientName     string              `json:"recipient_name"`
	Phone             string              `json:"phone"`
	CurrentStatus     TrackingStatus      `json:"current_status"`
	TrackingList      []*TrackingListItem `json:"tracking_list"`
	StatusList        []*StatusListItem   `json:"status_list"`
}

type TrackResponse struct {
	Retcode int               `json:"retcode"`
	Message string            `json:"message"`
	Data    TrackResponseData `json:"data"`
}
