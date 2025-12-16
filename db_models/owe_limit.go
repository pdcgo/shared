package db_models

type OweLimitConfiguration struct {
	ID        uint64  `gorm:"primarykey" json:"id"`
	TeamID    uint64  `json:"team_id"`
	IsDefault bool    `json:"is_default"`
	ForTeamID *uint64 `json:"for_team_id"`
	Threshold float64 `json:"threshold"`
}

func (o *OweLimitConfiguration) GetEntityID() string {
	return "owe_limit_configration"
}
