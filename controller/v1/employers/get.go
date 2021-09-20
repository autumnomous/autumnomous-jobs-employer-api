package employers

import (
	"encoding/json"
	"log"
	"net/http"

	"autumnomous.com/bit-jobs-api/shared/repository/employers"
	"autumnomous.com/bit-jobs-api/shared/repository/employers/accountmanagement"
	"autumnomous.com/bit-jobs-api/shared/response"
	"autumnomous.com/bit-jobs-api/shared/services/security/jwt"
)

type JobsResponse struct {
	Jobs             []*accountmanagement.Job `json:"jobs"`
	TotalPostsBought int                      `json:"totalpostsbought"`
}

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

	jobs, totalPostsBought, err := repository.GetEmployerJobs(publicID)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	var resp JobsResponse

	resp.Jobs = jobs
	resp.TotalPostsBought = totalPostsBought

	response.SendJSON(w, resp)

}

func GetJob(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		response.SendJSONMessage(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// tokenClaims, err := jwt.GetStrClaims(r)

	// if err != nil {
	// 	log.Println(err)
	// 	response.SendJSONMessage(w, http.StatusBadRequest, response.FriendlyError)
	// 	return
	// }

	// publicID := tokenClaims["user"]

	// if publicID == "" {
	// 	response.SendJSONMessage(w, http.StatusBadRequest, response.FriendlyError)
	// 	return
	// }

	var details map[string]string
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&details)

	if details["publicid"] == "" {
		response.SendJSONMessage(w, http.StatusBadRequest, response.FriendlyError)
		return
	}
	repository := employers.NewEmployerRegistry().GetEmployerRepository()

	job, err := repository.GetJob(details["publicid"])

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	response.SendJSON(w, job)
}

func GetActiveJobPackages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.SendJSONMessage(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	repository := employers.NewEmployerRegistry().GetEmployerRepository()

	packages, err := repository.GetActiveJobPackages()

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	response.SendJSON(w, packages)

}

func GetEmployer(w http.ResponseWriter, r *http.Request) {

}
