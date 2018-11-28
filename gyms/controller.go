package gyms

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"firebase.google.com/go"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// firestore client and context
var fireCtx, fireClient = firestoreConnect()

// html default directory
const htmlDir = "pages"

//the timestamp of the last firestore update
var lastUpdate = time.Now()

//firestore credentials
var fireCredentials = Credentials{
	"service_account",
	"tis-iv-221022",
	"235bb4121e370ef619813810b3bf7ffa35767fbf",
	"-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCtWLXXKpPZjATJ\n1jl4isjDWrpi7QaJLeoNWgU9mE0COI9QUJZUsqoDwOGTkkn9NhZCYfjwuYtikjKE\nmftRdIOH51p0dvzRnN+NsJppYGkN1R7a8oxXgz9gIU09Xc9kTux16kwwUzi2Z+pu\nThi4zU9Pr3rU8pL97Z6woZ/Wpk3shdEB+w/BKnJmwNywQMiArlL33dMmUNE4nul4\nvaeuIZL3c8EZ9+UlvhbqbRzXxejMk2XrulP6yj/mZ6TSJZmHOcqU7p7msDnTuLZQ\nLy7uzJ2zi/1ZN8QbpbSwK03p542jEAgbWbMbXew1UwU8KzY3KIbhVWYjpVtgB948\nsL1o/18lAgMBAAECggEAFJJWBiJp6GPa/576NAl4WOKwAuRxJZ96yrnXrF8iKHTQ\n4S2nIYcQcSCXRJ2URUYBb0BUPDEwzvJdp2nkrt5+a4bzr6WdTDzrNNP0BzRnUzpr\nDn3BVPNg9oYNNg6neZZ6LmYvQ37k2JTNd+ixu2C2HFVNZO8irZrg/cI4h082KS7n\nt1wGU+8KCpX/4J+zFFcr3ORX6nmtGwVbDcGExsiAlKEyHR2BiVm6dNKD4xVBj6v2\n+Xr+N4go9/rFGCDDAroo0ah+3h6c8SS5Kki1KdM7FGgRV3O6mrXBbKNJXItYqDCG\n28U2fxatony5IfUzIToIY/s9RidpFLU3cV52GMdcUwKBgQD0mxKuPuWJH0jTeU+Y\n+cioPG4/5pU0IKQY97lD4EThjdW7BsbRZyFcJqU5kZoqBsZ2MflWP+WCgx24ZZbL\nAnkzf3bvVCK+q5Ct2k7wmXGNDNEAKAaUzgZOmX515QNFBLDa/og2klDZTXcI3vtg\nKVvIhgUIDbBeP3gYlzYnWJLRTwKBgQC1a967AeMFgAQsvigEIGgWdl9Wdcbq3c2N\noe4jUBU61Izuvn8MggXJHT74Es01DobBf2VU33d139l3Jk31bMz7N/xI8DKh49GF\nZA5xK5boknuJS9GDnkDxhaBRy6Dh9EI1gz+2KRBXclsVvGJGtaFfpFlcitE36Qo2\nbayA7K3jSwKBgEdbwgxhPvdM0CMZfdYj8Jzb3FH6A8cMSrMZ+ctKbu3aQeLo7DGE\nw5+tioAL8QyXo2gx1gqKY3s6ov37bQ1WcGNMqTbStbwoMvH1ARiBuzWp6oMAKkNZ\nA1AEyXa9U8Hbx3hrzvMUpk9uoO5OlskL58HND0S1MaGdJH0QB/VciqBTAoGBAJmm\nh1gY6/4Pgvml/1wnWiCUFoCydUsbmWi32Wlc/O37cHUPL6kXQfEn7NnLirLB381n\nqRmtvY4+jP6FmYcfo6esreXUUP2dZik0KasdgMzuquQIK6TuVhB33OUJsfNMnPqX\nc1FDDA0T1CLfjthWIhtPpUNkaneQzk50qqHyUf9rAoGBAMYmzr2YA0zPOZW+tMJM\nUWvQcs0m+Z26g9+3/mXETsHvnqorfIRwwqVomohH/VgBwDEvrr/K6ACllyuLQzYb\noSMUcQAqojQbzwItG6ijgYiI0orzqYIkkL8S4fCkzmps2Yg65e309m4iC/Oogwsp\n/+dNdpT0FBIv95snZqkWRlzp\n-----END PRIVATE KEY-----\n",
	"firebase-adminsdk-548v2@tis-iv-221022.iam.gserviceaccount.com",
	"107628466397807120042",
	"https://accounts.google.com/o/oauth2/auth",
	"https://oauth2.googleapis.com/token",
	"https://www.googleapis.com/oauth2/v1/certs",
	"https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-548v2%40tis-iv-221022.iam.gserviceaccount.com",
}

//Api credentials
var apiCredentials = "%cZ=sF,4cz2`&Hbw)suPXJ%0-_#*WUijz2#scf&v(-;+Q(&3R<)mn%oEgGH::W4"

// Controller ...
type Controller struct {
}

// firestoreConnect Initializes the FireStore Connection
func firestoreConnect() (context.Context, *firestore.Client) {
	ctx := context.Background()
	// Get a Firestore client.
	credentialsJson, err := json.Marshal(&fireCredentials)
	if err != nil {
		log.Println("Cant marshal FireStore Credentials")
		log.Fatalln(err)
	}
	sa := option.WithCredentialsJSON(credentialsJson)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Println("Cant load FireStore Credentials")
		log.Fatalln(err)
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Println("Cant load FireStore Credentials")
		log.Fatalln(err)
	}
	log.Println("Firestore Initialized")
	return ctx, client
}

// DetermineListenAddress gets the port of the system
func DetermineListenAddress() string {
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT not set")
		os.Exit(0)
	}
	return ":" + port
}

// GetQuiz GET /gyms
func (c *Controller) GetGyms(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("apiKey") != apiCredentials {
		errorFmt("GetGyms", http.StatusForbidden, "apiKey does not match", errors.New("wrong api_key"), w)
		return
	}
	gyms, err := getGyms()
	if err != nil {
		errorFmt("GetGyms", http.StatusInternalServerError, "an error occurred while getting the gyms from db", err, w)
		return
	}

	data, err := json.Marshal(gyms)
	if err != nil {
		errorFmt("GetGyms", http.StatusInternalServerError, "an error occurred while Marshalling the response", err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

// PostScore POST /
func (c *Controller) PostScore(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("apiKey") != apiCredentials {
		errorFmt("PostScore", http.StatusForbidden, "apiKey does not match", errors.New("wrong api_key"), w)
		return
	}

	score := UserScore{"", "", -1, time.Now()}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) // read the body of the request
	if err != nil {
		errorFmt("PostScore", http.StatusInternalServerError, "an error occurred while reading request", err, w)
		return
	}

	if err := r.Body.Close(); err != nil {
		errorFmt("PostScore", http.StatusInternalServerError, "an error occurred while closing the body", err, w)
		return
	}

	err = json.Unmarshal(body, &score)
	if err != nil {
		errorFmt("PostScore", http.StatusInternalServerError, "an error occurred while unmarshalling the body", err, w)
		return
	}
	if score.Id == "" || score.GymCode == "" || score.Score == -1 {
		err = errors.New("Bad request")
		errorFmt("PostScore", http.StatusUnprocessableEntity, "bad request", err, w)
		log.Printf("Bad request: %s", body)
		return
	}

	err = PostScore(score)
	if err != nil {
		errorFmt("PostScore", http.StatusInternalServerError, "an error occurred while posting the score", err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	return
}

// Home GET /
func (c *Controller) Home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, fmt.Sprintf("./%s/index.html", htmlDir))
	return
}

// UpdateTimeStamp GET /updateTimeStamp
func (c *Controller) UpdateTimeStamp(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("apiKey") != apiCredentials {
		errorFmt("UpdateTimeStamp", http.StatusForbidden, "apiKey does not match", errors.New("wrong api_key"), w)
		return
	}
	data, err := json.Marshal(&lastUpdate)
	if err != nil {
		errorFmt("UpdateTimeStamp", http.StatusInternalServerError, "an error occurred while Marshalling the response", err, w)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

// errorFmt formats controllers reoccurring errors
func errorFmt(errorName string, errorCode int, errorMsg string, err error, w http.ResponseWriter) {
	log.Printf("Error %s: %+v \n", errorName, errors.Wrap(err, errorMsg))
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(errorCode)
}
