package employers

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"autumnomous-jobs-employer-api/shared/repository/jobs"
	"autumnomous-jobs-employer-api/shared/response"
	"autumnomous-jobs-employer-api/shared/services/security/jwt"
)

type DeleteJobDetails struct {
	PublicID string `json:"publicid"`
}

func DeleteJob(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodDelete {
		response.SendJSONMessage(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var details DeleteJobDetails

	decoder.Decode(&details)

	if details.PublicID == "" {
		response.SendJSONMessage(w, http.StatusBadRequest, response.MissingRequiredValue)
		return
	}

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

	if publicID == "" {
		response.SendJSONMessage(w, http.StatusBadRequest, response.FriendlyError)
		return
	}

	repository := jobs.NewJobRegistry().GetJobRepository()

	job, err := repository.DeleteJob(publicID, details.PublicID)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusBadRequest, response.FriendlyError)
		return
	}

	response.SendJSON(w, job)
}
