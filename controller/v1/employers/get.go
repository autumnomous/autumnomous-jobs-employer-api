package employers

import (
	"log"
	"net/http"

	"autumnomous.com/bit-jobs-api/shared/repository/employers"
	"autumnomous.com/bit-jobs-api/shared/response"
	"autumnomous.com/bit-jobs-api/shared/services/security/jwt"
)

func GetJobs(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		response.SendJSONMessage(w, http.StatusMethodNotAllowed, "Method not allowed")
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

	jobs, err := repository.GetEmployerJobs(publicID)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	response.SendJSON(w, jobs)

}
