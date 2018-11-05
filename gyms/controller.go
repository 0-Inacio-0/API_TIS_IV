package gyms

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

var Gyms []Gym

const htmlDir = "pages"

//Controller ...
type Controller struct {
}

func Init() {
	log.Print("Reading data...")
	data, err := ioutil.ReadFile("data/output.json")
	err = json.Unmarshal(data, &Gyms)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Done!")
}

// GetQuiz GET /
func (c *Controller) GetGyms(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(Gyms)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode("Error Marshaling json"); err != nil {
			log.Printf("Error GetAcademia: %+v \n", errors.Wrap(err, "json failed to encode"))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Printf("Error GetAcademia: %+v \n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

// GetQuiz GET /
func (c *Controller) PostScore(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) // read the body of the request
	if err != nil {
		log.Printf("Error PostScore: %+v \n", errors.Wrap(err, "Error reading request "))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := r.Body.Close(); err != nil {
		log.Printf("Error PostScore: %+v \n", errors.Wrap(err, "Error closing the body of the request"))
	}
	log.Printf("Json request: %s\n", body)
	err = AddScore(body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Printf("Error AddQuiz: %+v \n", errors.Wrap(err, "Error encoding json error "))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Printf("Error AddQuiz: %+v \n", errors.Wrap(err, "Error encoding json error "))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	return
}

func (c *Controller) Home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, fmt.Sprintf("./%s/index.html", htmlDir))
	return
}
