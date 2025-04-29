package db_models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type RefID string

func (idref RefID) GetProductRefID() (RefID, error) {
	switch CheckRefType(string(idref)) {
	case VariantRef:
		data, _ := idref.ExtractData()
		data.WarehouseID = 0
		data.RefType = ProductRef
		if len(data.RefIDs) > 1 {
			data.RefIDs = []uint{data.RefIDs[0]}
		}

		return NewRefID(data)
	}

	return idref, fmt.Errorf("%s bukan kode variasi", idref)
}

func RefIDQueryStr(unknownRef string) string {
	var queryRef string

	switch CheckRefType(unknownRef) {
	case ProductRef:
		var q RefID = RefID(unknownRef)
		data, _ := q.ExtractData()
		data.WarehouseID = 0
		newq, _ := NewRefID(data)
		queryRef = string(newq)

	case VariantRef:
		var q RefID = RefID(unknownRef)
		data, _ := q.ExtractData()
		data.WarehouseID = 0
		newq, _ := NewRefID(data)
		queryRef = string(newq)

	case UnknownRef:
		querys := strings.Split(unknownRef, "-")
		if len(querys) >= 4 {
			querys[3] = "X"
		}
		queryRef = strings.Join(querys, "-")
	}

	return strings.ToLower(queryRef)
}

type RefType string

const (
	BundleRef  RefType = "B"
	ProductRef RefType = "P"
	VariantRef RefType = "V"
	UnknownRef RefType = ""
)

type TeamCode string

func (code TeamCode) GetTeamID(db *gorm.DB) (uint, error) {
	tim := &Team{}
	err := db.Model(&Team{}).
		Select("id", "team_code").
		Where("team_code = ?", code).
		Find(tim).
		Error

	if err != nil {
		return 0, err
	}

	if tim.ID == 0 {
		return 0, fmt.Errorf("timid of timcode %s not found", code)
	}
	return tim.ID, nil
}

func GetTeamCode(db *gorm.DB, teamID uint) (TeamCode, error) {
	if teamID == 0 {
		return "", errors.New("GetTeamCode teamid cant 0")
	}

	tim := &Team{}
	err := db.Model(&Team{}).
		Select("id", "team_code").
		Where("id = ?", teamID).
		Find(tim).
		Error

	if err != nil {
		return "", err
	}

	if tim.ID == 0 {
		return "", fmt.Errorf("tim %d not found", teamID)
	}

	if tim.TeamCode == "" {
		return "", fmt.Errorf("timcode of tim %d not found", teamID)
	}
	return tim.TeamCode, nil

}

type RefData struct {
	TeamCode    TeamCode `json:"team_code"`
	UserCode    string   `json:"user_code"`
	RefType     RefType  `json:"ref_type"`
	WarehouseID uint     `json:"warehouse_id"`
	RefIDs      []uint   `json:"ref_ids"`
}

func NewRefID(data *RefData) (RefID, error) {
	if len(data.RefIDs) == 0 {
		return "", errors.New("ref_id kosong")
	}
	if data.TeamCode == "" {
		return "", errors.New("ref teamid kosong")
	}

	if data.UserCode == "" {
		return "", errors.New("ref usercode kosong")
	}

	warehouseStr := "X"
	if data.WarehouseID != 0 {
		warehouseStr = fmt.Sprintf("%d", data.WarehouseID)
	}

	hasil := []string{}
	hasil = append(
		hasil,
		strings.ToUpper(string(data.TeamCode)),
		strings.ToUpper(data.UserCode),
		string(data.RefType),
		warehouseStr,
	)

	switch data.RefType {
	case BundleRef, ProductRef:
		if len(data.RefIDs) != 1 {
			return "", errors.New("variant ref id length error")
		}
	case VariantRef:
		if len(data.RefIDs) != 2 {
			return "", errors.New("variant ref id length error")
		}
	}

	for _, id := range data.RefIDs {
		if id == 0 {
			return "", errors.New("ref id have 0 value")
		}

		hasil = append(hasil,
			fmt.Sprintf("%X", id),
		)
	}

	refid := strings.Join(hasil, "-")
	return RefID(refid), nil
}

func (data RefID) String() string {
	return string(data)
}

func (idnya RefID) ExtractData() (*RefData, error) {

	data := RefData{
		RefType: UnknownRef,
	}

	slicedat := strings.Split(string(idnya), "-")
	slicelen := len(slicedat)
	if slicelen <= 4 {
		return &data, fmt.Errorf("length data not enough in ref data %s", idnya)
	}
	teamcode := slicedat[0]

	var warehouseID uint = 0
	if strings.ToUpper(slicedat[3]) != "X" {
		w, err := strconv.ParseUint(slicedat[3], 10, 32)
		if err != nil {
			return &data, err
		}
		warehouseID = uint(w)
	}

	data.TeamCode = TeamCode(teamcode)
	data.UserCode = slicedat[1]
	data.RefType = RefType(slicedat[2])
	data.WarehouseID = warehouseID
	data.RefIDs = []uint{}

	sliceref := []string{}
	switch data.RefType {
	case BundleRef, ProductRef:
		if slicelen < 5 {
			return &data, errors.New("ref not valid for bundle or product")
		}
		sliceref = slicedat[4:5]
	case VariantRef:
		if slicelen < 6 {
			return &data, errors.New("ref not valid for variant")
		}
		sliceref = slicedat[4:6]
	}

	for _, d := range sliceref {
		val, err := hexToUint(d)
		if err != nil {
			return &data, err
		}
		data.RefIDs = append(data.RefIDs, val)
	}

	return &data, nil
}

func CheckRefType(idnya string) RefType {
	refdata, _ := RefID(idnya).ExtractData()
	return refdata.RefType
}
