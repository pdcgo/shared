package authorization

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pdcgo/shared/db_models"
	"github.com/pdcgo/shared/interfaces/identity_iface"
	"gorm.io/gorm"
)

type AuthError struct {
	ErrMsg string `json:"err_msg"`
}

func AppKeyAuthRequired(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := getToken(ctx)
		appID, err := getAppID(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, &AuthError{
				ErrMsg: err.Error(),
			})
			return
		}
		appkey, err := getAppKey(db, appID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, &AuthError{
				ErrMsg: err.Error(),
			})
			return
		}

		identity := JwtIdentity{}
		err = identity.Deserialize(appkey.Key, token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, &AuthError{
				ErrMsg: err.Error(),
			})
			return
		}
		identity.UserID = appkey.UserID
		identity.From = db_models.RootTeamType
		identity.UserAgent = identity_iface.ThirdAppAgent

		expTime := time.Unix(int64(identity.StandardClaims.ExpiresAt), 0)
		if time.Now().After(expTime) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, &AuthError{
				ErrMsg: "token expired",
			})
			return
		}

		ctx.Set("identity", &identity)
		ctx.Next()

	}
}

func getAppKey(db *gorm.DB, appID uint) (*db_models.AppKey, error) {
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

func getAppID(ctx *gin.Context) (uint, error) {
	data := ctx.GetHeader("X-App-Id")
	if data == "" {
		return 0, errors.New("application id not found in header")
	}
	appID, err := strconv.ParseUint(data, 10, 32)
	return uint(appID), err
}

func getToken(ctx *gin.Context) string {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		token = ctx.Query("token_key")
	}

	token, _ = strings.CutPrefix(token, "Bearer ")
	return token
}
