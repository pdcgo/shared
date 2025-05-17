package wallet_service_test

import (
	"testing"

	"github.com/pdcgo/gudang/src/identity/mock_identity"
	"github.com/pdcgo/shared/interfaces/wallet_iface"
	"github.com/pdcgo/shared/sdks/wallet_service"
	"github.com/stretchr/testify/assert"
)

func TestCreateWallet(t *testing.T) {

	agent := mock_identity.NewMockAgent(1, "test")

	srv := wallet_service.NewWalletService(&wallet_iface.WalletServiceConfig{
		AppSecret: "test",
		Endpoint:  "http://localhost:8080/wallet_service/v1/",
	})

	t.Run("test create wallet", func(t *testing.T) {
		t.Run("test dengan team id 0", func(t *testing.T) {
			res := srv.CreateWallet(t.Context(), agent, 0)
			assert.NotNil(t, res.Err)
		})

		t.Run("test dengan team id ada", func(t *testing.T) {
			res := srv.CreateWallet(t.Context(), agent, 1)
			assert.Equal(t, "create_error", res.Err.Code)
		})

	})

	t.Run("test team wallet", func(t *testing.T) {
		wallet := srv.TeamWallet(t.Context(), agent, 1, wallet_iface.OpsBalance)
		t.Run("test get info wallet", func(t *testing.T) {
			res := wallet.Detail()
			assert.Nil(t, res.Err)
		})

		t.Run("testing send balance", func(t *testing.T) {

			res := wallet.SendAmount(
				&wallet_iface.SendPayload{
					ToTeamID: 2,
					RefID:    "asdasdasdasd",
					Type:     wallet_iface.TxTypeBrokenGood,
					Amount:   20000,
				},
			)

			assert.NotEqual(t, "sdk_error", res.Err.Code)
		})
	})

}
