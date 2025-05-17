package wallet_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/pdcgo/shared/interfaces/identity_iface"
	"github.com/pdcgo/shared/interfaces/wallet_iface"
)

type teamWalletImpl struct {
	teamID uint
	tipe   wallet_iface.BalanceType

	httpClient *httpClient
	ctx        context.Context
	agent      identity_iface.Agent
	err        error
}

// Err implements wallet_iface.TeamWallet.
func (t *teamWalletImpl) Err() error {
	return t.err
}

// Info implements wallet_iface.TeamWallet.
func (t *teamWalletImpl) Detail() *wallet_iface.WalletInfoRes {
	var res wallet_iface.WalletInfoRes

	req, err := t.httpClient.createReq(t.ctx, t.agent, http.MethodGet, fmt.Sprintf("wallets/%d/details", t.teamID), nil, nil)
	if err != nil {
		res.Err = &ErrSdk
	}

	httpres, err := t.httpClient.Do(req)
	// data, _ := io.ReadAll(httpres.Body)
	// log.Println(string(data))
	if err != nil {
		var reserr *wallet_iface.ResultErr
		if errors.As(err, &reserr) {
			res.Err = reserr
		} else {
			res.Err = &ErrReq
		}

		return &res
	}

	err = json.NewDecoder(httpres.Body).Decode(&res)
	if err != nil {
		// log.Println(err, httpres.Request.URL.String(), "asdasdasd")
		res.Err = &ErrSdk
		return &res
	}
	return &res

}

// SendAmount implements wallet_iface.TeamWallet.
func (t *teamWalletImpl) SendAmount(payload *wallet_iface.SendPayload) *wallet_iface.CreateTransactionRes {
	var res wallet_iface.CreateTransactionRes

	req, err := t.httpClient.createReq(t.ctx, t.agent, http.MethodPost, fmt.Sprintf("wallets/%d/send_amount", t.teamID), nil, nil)
	if err != nil {
		res.Err = &ErrSdk
	}

	httpres, err := t.httpClient.Do(req)
	if err != nil {
		var reserr *wallet_iface.ResultErr
		if errors.As(err, &reserr) {
			res.Err = reserr
		} else {
			res.Err = &ErrReq
		}

		return &res
	}

	// data, _ := io.ReadAll(httpres.Body)
	// log.Println(string(data))
	err = json.NewDecoder(httpres.Body).Decode(&res)
	if err != nil {
		// log.Println(err, httpres.Request.URL.String(), "asdasdasd")
		res.Err = &ErrSdk
		return &res
	}
	return &res
}

func (t *teamWalletImpl) setError(err error) *teamWalletImpl {
	if err != nil {
		t.err = err
	}
	return t
}
