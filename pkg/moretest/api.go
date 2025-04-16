package moretest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/schema"
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

type MockJwtIdentity struct {
	jwt.StandardClaims
	UserID     uint
	SuperUser  bool
	ValidUntil int64
	CreatedAt  int64
}

func (j *MockJwtIdentity) Serialize(passphrase string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, j)

	tokenstring, err := token.SignedString([]byte(passphrase))

	return tokenstring, err
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
			jwt := &MockJwtIdentity{
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
