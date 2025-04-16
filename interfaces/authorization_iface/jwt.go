package authorization_iface

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pdcgo/gudang/services/ware_db"
	"github.com/pdcgo/shared/interfaces/identity_iface"
	"gorm.io/gorm"
)

type JwtIdentity struct {
	jwt.StandardClaims
	UserID     uint
	SuperUser  bool
	ValidUntil int64
	From       ware_db.TeamType
	UserAgent  identity_iface.AgentType
	CreatedAt  int64
}

// GetAgentType implements Identity.
func (j *JwtIdentity) GetAgentType() identity_iface.AgentType {
	return j.UserAgent
}

// IsTokenExpired implements Identity.
func (j *JwtIdentity) IsTokenExpired(tx *gorm.DB) (bool, error) {
	expTime := time.Unix(int64(j.StandardClaims.ExpiresAt), 0)
	if time.Now().After(expTime) {
		return true, nil
	}
	return false, nil
}

// GetExpired implements authorization_iface.Identity.
func (j *JwtIdentity) GetExpired(tx *gorm.DB) (*ExpiredToken, error) {
	var err error
	usr := ware_db.User{}

	err = tx.Model(&ware_db.User{}).First(&usr, j.UserID).Error

	return &ExpiredToken{
		LastCreated: usr.LastCreated,
		LastReset:   usr.LastReset,
	}, err
}

// GetUserID implements Identity.
func (j *JwtIdentity) GetUserID() uint {
	return j.UserID
}

// IdentityID implements Identity.
func (j *JwtIdentity) IdentityID() uint {
	return j.UserID
}

// IsSuperUser implements Identity.
func (j *JwtIdentity) IsSuperUser() bool {
	return j.SuperUser
}

// HasRole implements Identity.
func (u *JwtIdentity) HasRole(db *gorm.DB, domainID uint, keyname string) (bool, error) {
	var roleCount int64

	roles := db.Model(&Role{}).
		Select("id").
		Where("domain_id = ?", domainID).
		Where("key = ?", keyname)

	err := db.Model(&UserRole{}).
		Where("role_id = (?)", roles).
		Where("user_id = ?", u.IdentityID()).
		Count(&roleCount).Error
	if err != nil {
		return false, err
	}

	hasRole := roleCount > 0
	return hasRole, err
}

func (j *JwtIdentity) Serialize(passphrase string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, j)

	tokenstring, err := token.SignedString([]byte(passphrase))

	return tokenstring, err
}

func (j *JwtIdentity) Deserialize(passphrase string, tokenstring string) error {

	token, err := jwt.ParseWithClaims(tokenstring, &JwtIdentity{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", errors.New("gagal validate algorithm")
		}

		return []byte(passphrase), nil
	})

	if err != nil {
		return err
	}

	claims, ok := token.Claims.(*JwtIdentity)

	if ok && token.Valid {
		*j = *claims

		return nil
	}

	return errors.New("identity error")
}
