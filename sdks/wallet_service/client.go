package wallet_service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
	"github.com/pdcgo/shared/interfaces/identity_iface"
	"github.com/pdcgo/shared/interfaces/wallet_iface"
)

type httpClient struct {
	endpoint string
	secret   string
}

func (c *httpClient) Do(req *http.Request) (*http.Response, error) {
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return res, &ErrReq
	}

	switch res.StatusCode {
	case 404:
		return res, &ErrDeprecatedOrNotFound
	case 500:
		return res, &wallet_iface.ResultErr{
			Code: "unknown_error",
		}
	case 401:
		return res, &wallet_iface.ResultErr{
			Code:    "not_authorized",
			Message: "pastikan anda sudah login dengan benar karena sesi tidak terdeteksi",
		}
	}

	return res, err
}

func (c *httpClient) createReq(
	ctx context.Context,
	agent identity_iface.Agent,
	method string,
	path string,
	params,
	data interface{},
) (*http.Request, error) {
	var err error
	uri, _ := url.JoinPath(c.endpoint, path)

	buff := bytes.NewBuffer(nil)
	if data != nil {
		err = json.NewEncoder(buff).Encode(data)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, uri, buff)
	if err != nil {
		return req, err
	}

	if params != nil {
		values, err := query.Values(params)
		if err != nil {
			return req, err
		}
		req.URL.RawQuery = values.Encode()
	}

	// generating token agent
	token, err := agent.GetToken("wallet_service", c.secret)
	if err != nil {
		return req, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	return req, nil
}
