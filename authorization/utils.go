package authorization

import (
	"context"
	"errors"
)

func GetIdentity(ctx context.Context) (*JwtIdentity, error) {
	val := ctx.Value("identity")
	if val == nil {
		return nil, errors.New("identity not found")
	}

	var identity *JwtIdentity = val.(*JwtIdentity)
	return identity, nil
}
