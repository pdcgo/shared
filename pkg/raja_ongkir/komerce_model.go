package raja_ongkir

import "time"

type KWaybillRes struct {
	Meta Meta `json:"meta"`
	Data Data `json:"data"`
}
type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}
type KSummary struct {
	CourierCode   string        `json:"courier_code"`
	CourierName   string        `json:"courier_name"`
	WaybillNumber string        `json:"waybill_number"`
	ServiceCode   string        `json:"service_code"`
	WaybillDate   string        `json:"waybill_date"`
	ShipperName   string        `json:"shipper_name"`
	ReceiverName  string        `json:"receiver_name"`
	Origin        string        `json:"origin"`
	Destination   string        `json:"destination"`
	Status        DeliverStatus `json:"status"`
}
type KDetails struct {
	WaybillNumber    string `json:"waybill_number"`
	WaybillDate      string `json:"waybill_date"`
	WaybillTime      string `json:"waybill_time"`
	Weight           string `json:"weight"`
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
type KDeliveryStatus struct {
	Status      string `json:"status"`
	PodReceiver string `json:"pod_receiver"`
	PodDate     string `json:"pod_date"`
	PodTime     string `json:"pod_time"`
}
type Manifest struct {
	ManifestCode        string `json:"manifest_code,omitempty"`
	ManifestDescription string `json:"manifest_description,omitempty"`
	ManifestDate        string `json:"manifest_date"`
	ManifestTime        string `json:"manifest_time"`
	CityName            string `json:"city_name,omitempty"`
}

func (m *Manifest) GetTimestamp() (int64, error) {
	// Combine date and time
	datetimeStr := m.ManifestDate + " " + m.ManifestTime
	layout := "2006-01-02 15:04:05" // Go's reference layout

	parsedTime, err := time.ParseInLocation(layout, datetimeStr, time.Local)
	if err != nil {
		parsedTime, err := time.ParseInLocation("2006-01-02 15:04", datetimeStr, time.Local)
		if err != nil {
			return 0, err
		}
		return parsedTime.Unix(), nil
	}

	return parsedTime.Unix(), nil
}

type Data struct {
	Delivered      bool            `json:"delivered"`
	Summary        KSummary        `json:"summary"`
	Details        KDetails        `json:"details"`
	DeliveryStatus KDeliveryStatus `json:"delivery_status"`
	Manifest       []Manifest      `json:"manifest"`
}
