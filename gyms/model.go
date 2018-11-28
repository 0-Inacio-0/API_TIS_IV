package gyms

import (
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
	"log"
	"time"
)

type Credentials struct {
	AccType                 string `json:"type"`
	ProjectId               string `json:"project_id"`
	PrivateKeyId            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientId                string `json:"client_id"`
	AuthUri                 string `json:"auth_uri"`
	TokenUri                string `json:"token_uri"`
	AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url"`
	ClientX509CertUrl       string `json:"client_x509_cert_url"`
}

// Gym is all info about the gym
type Gym struct {
	EquipType    string      `json:"equip_type" firestore:"equip_type,omitempty"`
	EquipName    string      `json:"equip_name" firestore:"equip_name,omitempty"`
	Address      string      `json:"address" firestore:"address,omitempty"`
	Neighborhood string      `json:"neighborhood" firestore:"neighborhood,omitempty"`
	Region       string      `json:"region" firestore:"region,omitempty"`
	Responsible  string      `json:"responsible" firestore:"responsible,omitempty"`
	Code         string      `json:"code" firestore:"code,omitempty"`
	Theme        string      `json:"theme" firestore:"theme,omitempty"`
	Source       string      `json:"source" firestore:"source,omitempty"`
	Lat          float64     `json:"lat" firestore:"lat,omitempty"`
	Lon          float64     `json:"lon" firestore:"lon,omitempty"`
	Score        float64     `json:"score" firestore:"score,omitempty"`
	NScore       int         `json:"n_score" firestore:"n_score,omitempty"`
	UsersScores  []UserScore `json:"users_scores" firestore:"users_scores,omitempty"`
}

// UserScore is the info about the score of a gym from a user
type UserScore struct {
	Id      string    `json:"id" firestore:"id,omitempty"`
	GymCode string    `json:"gym_code" firestore:"gym_code,omitempty"`
	Score   int       `json:"score" firestore:"score,omitempty"`
	Date    time.Time `json:"date" firestore:"date,omitempty"`
}

// getGyms returns a slice of gym fetched from firestore
func getGyms() ([]Gym, error) {
	var gyms []Gym
	gym := Gym{}
	iter := fireClient.Collection("Gyms").Documents(fireCtx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			errors.Wrap(err, "failed to iterate while feting gyms from firestore")
			return nil, err
		}
		doc.DataTo(&gym)
		gyms = append(gyms, gym)
	}
	return gyms, nil
}

// PostScore checks if the user can set a score and push it to
// firestore returns an error in case of the gym does note exists
// or the user cant vote.
func PostScore(score UserScore) error {
	gym := Gym{}

	iter := fireClient.Collection("Gyms").Where("code", "==", score.GymCode).Documents(fireCtx)
	found := false
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			if !found {
				err = errors.New("no gym found with gym_code:" + score.GymCode)
				return err
			}
			break
		}
		if err != nil {
			return errors.Wrap(err, "error while getting gym in firestore")
		}
		found = true
		doc.DataTo(&gym)
	}

	score.Date = time.Now()
	gymScore := 0
	gymsScoreLen := len(gym.UsersScores) + 1
	for _, ele := range gym.UsersScores {
		if ele.Id == score.Id {
			log.Println("bateu")
			date := ele.Date
			if date.AddDate(0, 3, 0).After(time.Now()) {
				//todo set error types for comp outside this func
				err := errors.New("this user(user_id:" + score.Id + ") already voted in the last 3 months")
				return err
			} else {
				gymsScoreLen--
			}
		} else {
			gymScore += ele.Score
		}
	}
	gymScore += score.Score
	gym.UsersScores = append(gym.UsersScores, score)

	gym.Score = float64(gymScore / gymsScoreLen)
	_, err := fireClient.Collection("Gyms").Doc(score.GymCode).Set(fireCtx, gym)
	if err != nil {
		return errors.Wrap(err, "error while adding the updated gym in firestore")
	}
	lastUpdate = time.Now()
	return nil
}
