package employers

import (
	"encoding/json"
	"log"
	"net/http"

	"bit-jobs-api/shared/repository/jobs"
	"bit-jobs-api/shared/response"
	"bit-jobs-api/shared/services/security/jwt"
)

type createJobDetails struct {
	Title             string `json:"title"`
	JobType           string `json:"jobtype"`
	Category          string `json:"category"`
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

	publicID := jwt.GetUserClaim(r)

	repository := jobs.NewJobRegistry().GetJobRepository()

	// TODO: jobDetails.PostEndDatetime = jobDetails.PostStartDatetime + 30 days

	job, err := repository.EmployerCreateJob(publicID, jobDetails.Title, jobDetails.JobType, jobDetails.Category, jobDetails.Description, jobDetails.PostStartDatetime, jobDetails.PostEndDatetime, jobDetails.PayPeriod, jobDetails.MinSalary, jobDetails.MaxSalary)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	response.SendJSON(w, job)

}
