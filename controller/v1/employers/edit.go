package employers

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"bit-jobs-api/shared/repository/jobs"
	"bit-jobs-api/shared/response"
	"bit-jobs-api/shared/services/security/jwt"
)

type editJobDetails struct {
	Title             string `json:"title"`
	JobType           string `json:"jobtype"`
	Category          string `json:"category"`
	Description       string `json:"description"` // make required?
	PostStartDatetime string `json:"poststartdatetime"`
	Remote            bool   `json:"remote"`
	PublicID          string `json:"publicid"`
}

func EditJob(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		response.SendJSONMessage(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var details editJobDetails

	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&details)

	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)[1]
	authKey, err := base64.StdEncoding.DecodeString(auth)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusBadRequest, response.FriendlyError)
		return
	}

	tokenClaims, err := jwt.ParseToken(string(authKey))

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusBadRequest, response.FriendlyError)
		return
	}

	publicID := tokenClaims.CustomClaims["user"]

	repository := jobs.NewJobRegistry().GetJobRepository()

	job, err := repository.EditJob(publicID, details.PublicID, details.Title, details.JobType, details.Category, details.Description, details.PostStartDatetime, details.Remote)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	response.SendJSON(w, job)
}
