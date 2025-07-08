package authorization

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pdcgo/shared/db_models"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

func muxGetToken(r *http.Request) string {
	token := r.Header.Get("Authorization")
	if token == "" {
		token = r.URL.Query().Get("token_key")
	}
	token, _ = strings.CutPrefix(token, "Bearer ")
	return token
}

var ErrAppIDNotSet = errors.New("app id not set")

func muxGetAppID(r *http.Request) (uint, error) {
	dd := r.Header.Get("X-App-Id")
	if dd == "" {
		return 0, ErrAppIDNotSet
	}
	appID, err := strconv.ParseUint(dd, 10, 32)
	return uint(appID), err
}

func muxSetError(w http.ResponseWriter, code int, data any) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	raw, err := json.Marshal(data)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(raw)
}

func muxGetAppKey(db *gorm.DB, appID uint) (*db_models.AppKey, error) {

	if appID == 0 {
		return nil, ErrAppIDNotSet
	}

	appKey := db_models.AppKey{}
	err := db.
		Model(&db_models.AppKey{}).
		Where("id = ?", appID).
		Find(&appKey).
		Error

	if err != nil {
		return &appKey, err
	}

	if appKey.ID == 0 {
		return &appKey, errors.New("appkey not found")
	}

	return &appKey, nil
}

func MuxAuthMiddleware(db *gorm.DB, next http.Handler, defaultphrase string) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := muxGetToken(r)
		appID, err := muxGetAppID(r)

		appkey, err := muxGetAppKey(db, appID)
		if err != nil {
			if errors.Is(err, ErrAppIDNotSet) {
				appkey = &db_models.AppKey{
					AppName: "default",
					Key:     defaultphrase,
				}

			} else {
				muxSetError(w, http.StatusUnauthorized, &AuthError{
					ErrMsg: err.Error(),
				})
				return
			}
		}

		identity := JwtIdentity{}
		err = identity.Deserialize(appkey.Key, token)
		if err != nil {
			muxSetError(w, http.StatusUnauthorized, &AuthError{
				ErrMsg: err.Error(),
			})
			return
		}
		// identity.UserID = appkey.UserID
		// identity.From = db_models.RootTeamType
		// identity.UserAgent = identity_iface.ThirdAppAgent

		expTime := time.Unix(int64(identity.StandardClaims.ExpiresAt), 0)
		if time.Now().After(expTime) {
			muxSetError(w, http.StatusUnauthorized, &AuthError{
				ErrMsg: "token expired",
			})
			return
		}
		pctx := r.Context()
		ctx := context.WithValue(pctx, "identity", &identity)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
