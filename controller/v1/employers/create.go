package employers

import (
	"encoding/json"
	"log"
	"net/http"

	"bit-jobs-api/shared/repository/employers"
	"bit-jobs-api/shared/response"
	"bit-jobs-api/shared/services/security/jwt"
)

// "jobtitle":          "",
// 			"jobstreetaddress":  "",
// 			"jobcity":           "",
// 			"jobzipcode":        "",
// 			"jobtype":           "",
// 			"jobremotefriendly": "",
// 			"jobdescription":    "",

type createJobDetails struct {
	Title             string `json:"title"`
	StreetAddress     string `json:"streetaddress"`
	City              string `json:"city"`
	ZipCode           string `json:"zipcode"`
	Tags              string `json:"tags"`
	Description       string `json:"description"` // make required?
	MinSalary         int    `json:"minsalary"`
	MaxSalary         int    `json:"maxsalary"`
	PayPeriod         string `json:"payperiod"`
	PostStartDatetime string `json:"poststartdatetime"`
	PostEndDatetime   string `json:"postenddatetime"`
}

func CreateJob(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		response.SendJSONMessage(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	var jobDetails createJobDetails

	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&jobDetails)

	if jobDetails.Title == "" {
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

	job, err := repository.EmployerCreateJob(publicID, jobDetails.Title, jobDetails.StreetAddress, jobDetails.City, jobDetails.ZipCode, jobDetails.Tags, jobDetails.Description, jobDetails.PostStartDatetime, jobDetails.PostEndDatetime, jobDetails.PayPeriod, jobDetails.MinSalary, jobDetails.MaxSalary)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	response.SendJSON(w, job)

}
