package raja_ongkir

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
)

type ApiKey struct {
	c       int32
	keys    []string
	ckey    int32
	limiter *time.Timer
}

func NewApiKey(key []string) *ApiKey {
	ckey := len(key)
	t := time.NewTimer(time.Minute / time.Duration(120*ckey))

	return &ApiKey{
		c:       0,
		keys:    key,
		ckey:    int32(ckey),
		limiter: t,
	}
}

func (k *ApiKey) Key() string {
	<-k.limiter.C
	defer k.limiter.Reset(time.Minute / time.Duration(120*k.ckey))

	newc := atomic.AddInt32(&k.c, 1)
	if newc >= k.ckey {
		atomic.SwapInt32(&k.c, 0)
		newc = 0
	}

	return k.keys[newc]
}

func KomerceTrack(listKey *ApiKey, receipt, courrier string) (*KWaybillRes, error) {
	hasil := KWaybillRes{}
	url := fmt.Sprintf("https://rajaongkir.komerce.id/api/v1/track/waybill?awb=%s&courier=%s", receipt, courrier)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return &hasil, err
	}

	req.Header.Add("key", listKey.Key())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return &hasil, err
	}

	err = json.NewDecoder(res.Body).Decode(&hasil)
	if err != nil {
		return &hasil, err
	}

	switch hasil.Meta.Code {
	case 200:
		return &hasil, err
	default:
		err = errors.New(hasil.Meta.Message)
		return &hasil, err
	}

}

func Track(listKey *ApiKey, courrier string, receipt string) (*WaybillRes, error) {
	hasil := WaybillRes{}

	url := "https://pro.rajaongkir.com/api/waybill"
	payload := strings.NewReader(fmt.Sprintf("waybill=%s&courier=%s", receipt, courrier))
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return &hasil, err
	}

	req.Header.Add("key", listKey.Key())
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return &hasil, err
	}

	err = json.NewDecoder(res.Body).Decode(&hasil)
	if err != nil {
		return &hasil, err
	}
	return &hasil, err

}
