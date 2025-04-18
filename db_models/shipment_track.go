package db_models

import (
	"time"

	"gorm.io/datatypes"
)

// usahakan tidak terlalu dependend biar enak jika dibuat microservice

type ThirdPartyKey uint

const (
	RajaOngkir ThirdPartyKey = iota
	SpxUnderGround
)

type TrackStatus string

const (
	TrackDelivered     TrackStatus = "delivered"
	TrackReturnProcess TrackStatus = "return_process"
	TrackReturned      TrackStatus = "returned"
	ShipmentProcess    TrackStatus = "shipment_process"
	TrackLost          TrackStatus = "lost"
	DeliveryFailed     TrackStatus = "delivery_failed"
	TrackCancel        TrackStatus = "cancel"
	Created            TrackStatus = "created"
	Unknown            TrackStatus = "unknown"
)

func (TrackStatus) EnumList() []string {
	return []string{
		"delivered",
		"return_process",
		"returned",
		"shipment_process",
		"lost",
		"delivery_failed",
		"cancel",
		"created",
		"unknown",
	}
}

type TrackHistory struct {
	Name      string `json:"name"`
	Desc      string `json:"desc"`
	Timestamp int64  `json:"timestamp"`
}

type TrackInfo struct {
	ID           uint          `gorm:"primarykey" json:"id"`
	ShippingID   uint          `json:"shipping_id"`
	Receipt      string        `gorm:"index:receipt_unique,unique" json:"receipt"`
	Status       TrackStatus   `json:"status"`
	CheckWith    ThirdPartyKey `json:"check_with"`
	CheckSuccess bool          `json:"check_success"`
	ErrMsg       string        `json:"err_msg"`
	LastUpdated  time.Time     `json:"last_updated"`

	History datatypes.JSONSlice[TrackHistory] `json:"history"`
	// Extras  string                            `json:"extras"`

	Shipping Shipping `json:"shipping"`
}
