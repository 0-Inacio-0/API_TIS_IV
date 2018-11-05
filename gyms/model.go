package gyms

import (
	"encoding/json"
	"github.com/pkg/errors"
	"time"
)

type Gym struct {
	EquipType    string      `json:"equip_type"`
	EquipName    string      `json:"equip_name"`
	Address      string      `json:"address"`
	Neighborhood string      `json:"neighborhood"`
	Region       string      `json:"region"`
	Responsible  string      `json:"responsible"`
	Code         string      `json:"code"`
	Theme        string      `json:"theme"`
	Source       string      `json:"source"`
	Lat          float64     `json:"lat"`
	Lon          float64     `json:"lon"`
	Score        float64     `json:"score"`
	NScore       int         `json:"n_score"`
	UsersScores  []UserScore `json:"users_scores"`
}
type UserScore struct {
	Id      string    `json:"id"`
	GymCode string    `json:"gyms_code"`
	Score   float64   `json:"score"`
	Date    time.Time `json:"date"`
}

func AddScore(data []byte) error {
	score := UserScore{}
	err := json.Unmarshal(data, &score)
	if err != nil {
		return errors.Wrap(err, "Error Unmarshalling UserScore")
	}
	for _, ele := range Gyms {
		if ele.Code == score.GymCode {
			ele.UsersScores = append(ele.UsersScores, score)
			ele.NScore = len(ele.UsersScores)
			sum := 0.0
			for _, userScore := range ele.UsersScores {
				sum += userScore.Score
			}
			ele.Score = sum / float64(ele.NScore)
		}
	}
	return err
}
