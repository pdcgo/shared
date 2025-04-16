package moretest_test

type StatWtAsset struct {
	ATeamID int `json:"team_id"`
}

func (StatWtAsset) TableName() string {
	return "stat_wt_asset"
}

// func TestSkipBigquery(t *testing.T) {
// 	var db gorm.DB

// 	moretest.Suite(t, "testing skipping",
// 		moretest.SetupListFunc{
// 			moretest_mock.MockBigqueryDatabase(&db),
// 		},
// 		moretest.SkipGcpNotLogin(t, func(t *testing.T) {
// 			hasil := []*StatWtAsset{}
// 			db.Model(&StatWtAsset{}).Limit(5).Find(&hasil)
// 			debugtool.LogJson(hasil, ":asd")
// 		}),
// 	)
// }
