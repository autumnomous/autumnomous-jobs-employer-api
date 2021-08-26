package employers

import (
	"encoding/json"
	"log"
	"net/http"

	"autumnomous.com/bit-jobs-api/shared/repository/employers"
	"autumnomous.com/bit-jobs-api/shared/response"
	"autumnomous.com/bit-jobs-api/shared/services/security/jwt"
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

	tokenClaims, err := jwt.GetStrClaims(r)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusBadRequest, response.FriendlyError)
		return
	}

	publicID := tokenClaims["user"]

	if publicID == "" {
		response.SendJSONMessage(w, http.StatusBadRequest, response.FriendlyError)
		return
	}

	repository := employers.NewEmployerRegistry().GetEmployerRepository()

	job, err := repository.DeleteJob(publicID, details.PublicID)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusBadRequest, response.FriendlyError)
		return
	}

	response.SendJSON(w, job)
}
