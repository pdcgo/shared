package raja_ongkir_test

import (
	"testing"

	"github.com/pdcgo/shared/pkg/raja_ongkir"
	"github.com/stretchr/testify/assert"
)

var listKey *raja_ongkir.ApiKey = raja_ongkir.NewApiKey([]string{})

func TestRajaOngkir(t *testing.T) {
	t.Skip("ganti api")
	res, err := raja_ongkir.Track(listKey, "jnt", "JX3589413017")
	assert.Nil(t, err)

	t.Error(res.Rajaongkir.Status.Code)

	// debugtool.LogJson(d)

	t.Error(res.Rajaongkir.Result.DeliveryStatus)
}

func TestVersiKomerce(t *testing.T) {
	t.Skip("belum impelmentasi")
	// res, err := raja_ongkir.KomerceTrack("JX3589413017", "jnt")
	res, err := raja_ongkir.KomerceTrack(listKey, "JX3576753600", "jnt")

	assert.Nil(t, err)
	assert.NotEmpty(t, res)
}

func TestJneCm(t *testing.T) {
	t.Skip("belum impelmentasi")
	// res, err := raja_ongkir.KomerceTrack("JX3589413017", "jnt")
	res, err := raja_ongkir.KomerceTrack(listKey, "CM53517086145", "jne")

	assert.Nil(t, err)
	assert.NotEmpty(t, res)

	// debugtool.LogJson(res)

	t.Run("testing parsing tanggal", func(t *testing.T) {
		for _, item := range res.Data.Manifest {
			ts, err := item.GetTimestamp()
			assert.Nil(t, err)
			assert.NotEqual(t, 0, ts)
		}
	})

}
