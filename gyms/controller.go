package gyms

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
)

var Gyms []Academia

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
func (c *Controller) GetAcademias(w http.ResponseWriter, r *http.Request) {
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
