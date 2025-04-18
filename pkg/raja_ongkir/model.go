package raja_ongkir

type DeliverStatus string

const (
	Delivered     DeliverStatus = "DELIVERED"
	OnProcess     DeliverStatus = "ON PROCESS"
	ReturnProcess DeliverStatus = "RETURN TO SHIPPER"
)

type Query struct {
	Waybill string `json:"waybill"`
	Courier string `json:"courier"`
}

type Status struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

type Summary struct {
	CourierCode   string `json:"courier_code"`
	CourierName   string `json:"courier_name"`
	WaybillNumber string `json:"waybill_number"`
	ServiceCode   string `json:"service_code"`
	WaybillDate   string `json:"waybill_date"`
	ShipperName   string `json:"shipper_name"`
	ReceiverName  string `json:"receiver_name"`
	Origin        string `json:"origin"`
	Destination   string `json:"destination"`
	Status        string `json:"status"`
}

type DeliveryStatus struct {
	Status      DeliverStatus `json:"status"`
	PodReceiver string        `json:"pod_receiver"`
	PodDate     string        `json:"pod_date"`
	PodTime     string        `json:"pod_time"`
}

type ManifestItem struct {
	ManifestCode        string `json:"manifest_code,omitempty"`
	ManifestDescription string `json:"manifest_description,omitempty"`
	ManifestDate        string `json:"manifest_date"`
	ManifestTime        string `json:"manifest_time"`
	CityName            string `json:"city_name"`
}

type Details struct {
	WaybillNumber string `json:"waybill_number"`
	WaybillDate   string `json:"waybill_date"`
	WaybillTime   string `json:"waybill_time"`
	// Weight           int    `json:"weight"`
	Origin           string `json:"origin"`
	Destination      string `json:"destination"`
	ShipperName      string `json:"shipper_name"`
	ShipperAddress1  string `json:"shipper_address1"`
	ShipperAddress2  string `json:"shipper_address2"`
	ShipperAddress3  string `json:"shipper_address3"`
	ShipperCity      string `json:"shipper_city"`
	ReceiverName     string `json:"receiver_name"`
	ReceiverAddress1 string `json:"receiver_address1"`
	ReceiverAddress2 string `json:"receiver_address2"`
	ReceiverAddress3 string `json:"receiver_address3"`
	ReceiverCity     string `json:"receiver_city"`
}

type Result struct {
	Delivered      bool           `json:"delivered"`
	Summary        Summary        `json:"summary"`
	Details        Details        `json:"details"`
	DeliveryStatus DeliveryStatus `json:"delivery_status"`
	Manifest       []ManifestItem `json:"manifest"`
}

type RajaOngkir struct {
	Query  Query  `json:"query"`
	Status Status `json:"status"`
	Result Result `json:"result"`
}

type WaybillRes struct {
	Rajaongkir RajaOngkir `json:"rajaongkir"`
}
