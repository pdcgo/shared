package spx_tracker

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type DQuery struct {
	Project string `json:"project"`
	Key     string `json:"key"`
}
type TrackQuery struct {
	Bt  string `json:"bt"`
	Sbt string `json:"sbt"`
	Dt  string `json:"dt"`
	R   int    `json:"r"`
	UID string `json:"uid"`
	A   string `json:"a"`
	Pf  string `json:"pf"`
	S   string `json:"s"`
	Env string `json:"env"`
	D   DQuery `json:"d"`
	Ct  int64  `json:"ct"`
}

var spxKey = "MGViZmZmZTYzZDJhNDgxY2Y1N2ZlN2Q1ZWJkYzlmZDY="

type TrackClient struct {
	C *http.Client
}

func (c *TrackClient) generateKey(receipt string) (time.Time, string, error) {

	ts := time.Now()
	key := fmt.Sprintf("%s%d%s", receipt, ts.Unix(), spxKey)

	hasher := sha256.New()
	hasher.Write([]byte(key))

	sha := hasher.Sum(nil)

	return ts, fmt.Sprintf("%x", sha), nil
}

func (c *TrackClient) Track(receipt string) (*TrackResponse, error) {
	resd := TrackResponse{}
	ts, key, err := c.generateKey(receipt)
	if err != nil {
		return &resd, err
	}

	uri := fmt.Sprintf("https://spx.co.id/api/v2/fleet_order/tracking/search?sls_tracking_number=%s|%d%s", receipt, ts.Unix(), key)

	res, err := c.C.Get(uri)
	if err != nil {
		return &resd, err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return &resd, err
	}

	if res.StatusCode != http.StatusOK {
		return &resd, errors.New(string(data))
	}

	err = json.Unmarshal(data, &resd)
	if err != nil {
		return &resd, err
	}

	switch resd.Retcode {
	case 0:
		return &resd, nil
	default:
		return &resd, errors.New(resd.Message)
	}

	return &resd, nil
}
