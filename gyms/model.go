package gyms

import "time"

type Academia struct {
	EquipType    string `json:"equip_type"`
	EquipName    string `json:"equip_name"`
	Address      string `json:"address"`
	Neighborhood string `json:"neighborhood"`
	Region       string `json:"region"`
	Responsible  string `json:"responsible"`
	Code         string `json:"code"`
	Theme        string `json:"theme"`
	Source       string `json:"source"`
	Lat          float64 `json:"lat"`
	Lon          float64 `json:"lon"`
	Score		 float64 `json:"score"`
	NScore		 int `json:"n_score"`
	UsersScores  []UserScore `json:"users_scores"`
}
type UserScore struct {
	Id string `json:"id"`
	Score float64 `json:"score"`
	Date time.Time `json:"date"`
}