package moretest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/schema"
	"github.com/pdcgo/shared/interfaces/authorization_iface"
	"github.com/pdcgo/v2_gots_sdk"
	"github.com/pdcgo/v2_gots_sdk/pdc_api"
	"github.com/stretchr/testify/assert"
)

var encoder = schema.NewEncoder()

type RequesterFunc func(api *pdc_api.Api) *httptest.ResponseRecorder

func SetupApiRequester(group *v2_gots_sdk.SdkGroup, handler *RequesterFunc) SetupFunc {

	return func(t *testing.T) func() error {

		*handler = func(api *pdc_api.Api) *httptest.ResponseRecorder {
			r := group.GetGinEngine()
			w := httptest.NewRecorder()

			data := bytes.NewBuffer(nil)
			if api.Payload != nil {
				err := json.NewEncoder(data).Encode(api.Payload)
				assert.Nil(t, err)
			}
			req, err := http.NewRequest(api.Method, api.RelativePath, data)
			assert.Nil(t, err)

			if api.Query != nil {
				q := req.URL.Query()
				encoder.Encode(api.Query, q)
				req.URL.RawQuery = q.Encode()
			}

			r.ServeHTTP(w, req)
			return w
		}

		return nil
	}
}

type ConvertableToken interface {
	GetUserID() uint
}

func SetupAuthApiRequester(user ConvertableToken, group *v2_gots_sdk.SdkGroup, handler *RequesterFunc) SetupFunc {
	return func(t *testing.T) func() error {

		*handler = func(api *pdc_api.Api) *httptest.ResponseRecorder {
			r := group.GetGinEngine()
			w := httptest.NewRecorder()

			data := bytes.NewBuffer(nil)
			if api.Payload != nil {
				err := json.NewEncoder(data).Encode(api.Payload)
				assert.Nil(t, err)
			}
			req, err := http.NewRequest(api.Method, api.RelativePath, data)
			assert.Nil(t, err)

			if api.Query != nil {
				q := req.URL.Query()
				encoder.Encode(api.Query, q)
				req.URL.RawQuery = q.Encode()
			}

			now := time.Now().Add(time.Hour * 24)
			// setting authorization
			jwt := &authorization_iface.JwtIdentity{
				UserID:     user.GetUserID(),
				ValidUntil: now.UnixMicro(),
			}

			token, err := jwt.Serialize("test_phrase")
			assert.Nil(t, err)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

			r.ServeHTTP(w, req)
			return w
		}

		return nil
	}
}
