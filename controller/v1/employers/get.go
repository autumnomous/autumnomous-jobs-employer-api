package employers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"autumnomous-jobs-employer-api/shared/repository/employers"
	"autumnomous-jobs-employer-api/shared/repository/jobpackages"
	"autumnomous-jobs-employer-api/shared/repository/jobs"
	"autumnomous-jobs-employer-api/shared/response"
	"autumnomous-jobs-employer-api/shared/services/security/jwt"
	"autumnomous-jobs-employer-api/shared/services/zipcode"
)

type JobsResponse struct {
	Jobs []*jobs.Job `json:"jobs"`
}

type AutocompleteLocationData struct {
	Characters string `json:"chars"`
}

func GetJobs(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		response.SendJSONMessage(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	publicID := jwt.GetUserClaim(r)

	if publicID == "" {
		response.SendJSONMessage(w, http.StatusBadRequest, response.FriendlyError)
		return
	}

	repository := jobs.NewJobRegistry().GetJobRepository()

	jobs, err := repository.GetEmployerJobs(publicID)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	response.SendJSON(w, jobs)

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
	repository := jobs.NewJobRegistry().GetJobRepository()

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

	repository := jobpackages.NewJobPackageRegistry().GetJobPackageRepository()

	packages, err := repository.GetActiveJobPackages()

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	response.SendJSON(w, packages)

}

func GetEmployer(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		response.SendJSONMessage(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	publicID := jwt.GetUserClaim(r)

	repository := employers.NewEmployerRegistry().GetEmployerRepository()

	employer, err := repository.GetEmployer(publicID)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	response.SendJSON(w, employer)
}

func GetEmployerCompany(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		response.SendJSONMessage(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	publicID := jwt.GetUserClaim(r)

	repository := employers.NewEmployerRegistry().GetEmployerRepository()

	company, err := repository.GetEmployerCompany(publicID)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	response.SendJSON(w, company)

}

func GetAutocompleteLocationData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.SendJSONMessage(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var details AutocompleteLocationData

	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&details)

	if details.Characters == "" {
		response.SendJSONMessage(w, http.StatusBadRequest, response.MissingRequiredValue)
		return
	}

	gateway := zipcode.NewZipCodeGateway(os.Getenv("ZIPCODESERVICES_API_KEY"))

	data, err := gateway.GetAutoComplete(details.Characters)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	response.SendJSON(w, data)

}
