package employers

import (
	"encoding/json"
	"log"
	"net/http"

	"autumnomous.com/bit-jobs-api/shared/repository/employers"
	"autumnomous.com/bit-jobs-api/shared/response"
	"autumnomous.com/bit-jobs-api/shared/services/security/jwt"
)

// "jobtitle":          "",
// 			"jobstreetaddress":  "",
// 			"jobcity":           "",
// 			"jobzipcode":        "",
// 			"jobtype":           "",
// 			"jobremotefriendly": "",
// 			"jobdescription":    "",

type createJobDetails struct {
	JobTitle         string `json:"jobtitle"`
	JobStreetAddress string `json:"jobstreetaddress"`
	JobCity          string `json:"jobcity"`
	JobZipCode       string `json:"jobzipcode"`
	JobTags          string `json:"jobtags"`
	JobDescription   string `json:"jobdescription"` // make required?
}

func CreateJob(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		response.SendJSONMessage(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	var jobDetails createJobDetails

	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&jobDetails)

	if jobDetails.JobTitle == "" {
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

	repository := employers.NewEmployerRegistry().GetEmployerRepository()

	job, err := repository.EmployerCreateJob(publicID, jobDetails.JobTitle, jobDetails.JobStreetAddress, jobDetails.JobCity, jobDetails.JobZipCode, jobDetails.JobTags, jobDetails.JobDescription)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	response.SendJSON(w, job)

}
