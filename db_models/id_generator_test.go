package db_models_test

import (
	"testing"

	"github.com/pdcgo/shared/db_models"
	"github.com/stretchr/testify/assert"
)

func TestIdGenerator(t *testing.T) {
	t.Run("testing bundle sku", func(t *testing.T) {
		refb, err := db_models.NewRefID(&db_models.RefData{
			TeamCode: "h2O",
			UserCode: "gh",
			RefType:  db_models.BundleRef,
			RefIDs:   []uint{12},
		})

		assert.Equal(t, db_models.RefID("H2O-GH-B-X-C"), refb)
		assert.Nil(t, err)

		// ---------------- baru -------------

		t.Run("extract bundle sku id", func(t *testing.T) {
			ref := db_models.RefID("H2O-GH-B-X-C")

			dat, err := ref.ExtractData()
			assert.Nil(t, err)
			assert.Equal(t, db_models.TeamCode("H2O"), dat.TeamCode)
			assert.Equal(t, "GH", dat.UserCode)
			assert.Equal(t, db_models.BundleRef, dat.RefType)
			assert.Equal(t, uint(12), dat.RefIDs[0])

			// ---------------- baru -------------

		})
	})

	t.Run("testing product id", func(t *testing.T) {

		refp, err := db_models.NewRefID(&db_models.RefData{
			TeamCode: "h2O",
			UserCode: "gh",
			RefType:  db_models.ProductRef,
			RefIDs:   []uint{12},
		})

		assert.Equal(t, db_models.RefID("H2O-GH-P-X-C"), refp)
		assert.Nil(t, err)

		// ---------------- baru -------------

		t.Run("extract product sku id", func(t *testing.T) {
			ref := db_models.RefID("H2O-GH-P-X-C")

			dat, err := ref.ExtractData()
			assert.Nil(t, err)
			assert.Equal(t, db_models.TeamCode("H2O"), dat.TeamCode)
			assert.Equal(t, "GH", dat.UserCode)
			assert.Equal(t, db_models.ProductRef, dat.RefType)
			assert.Equal(t, uint(12), dat.RefIDs[0])

		})
	})

	t.Run("testing variant sku id", func(t *testing.T) {

		refv, _ := db_models.NewRefID(&db_models.RefData{
			TeamCode: "h2O",
			UserCode: "gh",
			RefType:  db_models.VariantRef,
			RefIDs:   []uint{12, 13},
		})

		assert.Equal(t, db_models.RefID("H2O-GH-V-X-C-D"), refv)

		t.Run("testing extract data variant skuid", func(t *testing.T) {
			ref := db_models.RefID("H2O-GH-V-X-C-D")

			dat, err := ref.ExtractData()
			assert.Nil(t, err)
			assert.Equal(t, db_models.TeamCode("H2O"), dat.TeamCode)
			assert.Equal(t, "GH", dat.UserCode)
			assert.Equal(t, db_models.VariantRef, dat.RefType)
			assert.Equal(t, uint(12), dat.RefIDs[0])
			assert.Equal(t, uint(13), dat.RefIDs[1])

		})

	})

	t.Run("cek getting ref type", func(t *testing.T) {
		tipe := db_models.CheckRefType("1-GH-V-X-C-D")
		assert.Equal(t, db_models.VariantRef, tipe)

		t.Run("testing cek dengan string len kurang dari tipe", func(t *testing.T) {
			tipe := db_models.CheckRefType("va")
			assert.Equal(t, db_models.UnknownRef, tipe)
		})

		t.Run("testing cek dengan string len kurang dari tipe", func(t *testing.T) {
			tipe := db_models.CheckRefType("var")
			assert.Equal(t, db_models.UnknownRef, tipe)
		})
	})

	t.Run("test get ref type bundle", func(t *testing.T) {
		refTypeBundle := db_models.CheckRefType("1-GH-B-X-D")
		assert.Equal(t, db_models.BundleRef, refTypeBundle)

		bundlesku, err := db_models.NewRefID(&db_models.RefData{
			TeamCode: "h2O",
			UserCode: "as",
			RefType:  db_models.BundleRef,
			RefIDs:   []uint{1},
		})
		assert.Nil(t, err)

		refType := db_models.CheckRefType(string(bundlesku))
		assert.Equal(t, db_models.BundleRef, refType)
	})

	t.Run("testing variant sku gudang id", func(t *testing.T) {

		refv, _ := db_models.NewRefID(&db_models.RefData{
			TeamCode:    "h2O",
			UserCode:    "gh",
			WarehouseID: 2,
			RefType:     db_models.VariantRef,
			RefIDs:      []uint{12, 13},
		})

		assert.Equal(t, db_models.RefID("H2O-GH-V-2-C-D"), refv)

		t.Run("testing extract data gudang skuid", func(t *testing.T) {
			ref := db_models.RefID("H2O-GH-V-2-C-D")

			dat, err := ref.ExtractData()
			assert.Nil(t, err)
			assert.Equal(t, db_models.TeamCode("H2O"), dat.TeamCode)
			assert.Equal(t, "GH", dat.UserCode)
			assert.Equal(t, db_models.VariantRef, dat.RefType)
			assert.Equal(t, uint(12), dat.RefIDs[0])
			assert.Equal(t, uint(13), dat.RefIDs[1])

		})

	})

	t.Run("cek getting ref type", func(t *testing.T) {
		tipe := db_models.CheckRefType("H2O-GH-V-2-C-D")
		assert.Equal(t, db_models.VariantRef, tipe)

		t.Run("testing cek dengan string len kurang dari tipe", func(t *testing.T) {
			tipe := db_models.CheckRefType("va")
			assert.Equal(t, db_models.UnknownRef, tipe)
		})

		t.Run("testing cek dengan string len kurang dari tipe", func(t *testing.T) {
			tipe := db_models.CheckRefType("var")
			assert.Equal(t, db_models.UnknownRef, tipe)
		})
	})
}
