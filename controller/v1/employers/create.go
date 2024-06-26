package employers

import (
	"encoding/json"
	"log"
	"net/http"

	"autumnomous-jobs-employer-api/shared/repository/jobs"
	"autumnomous-jobs-employer-api/shared/response"
	"autumnomous-jobs-employer-api/shared/services/security/jwt"
)

type createJobDetails struct {
	Title       string `json:"title"`
	JobType     string `json:"jobtype"`
	Category    string `json:"category"`
	Description string `json:"description"` // make required?
	Remote      bool   `json:"remote"`
	VisibleDate string `json:"visibledate"`
	PayPeriod   string `json:"payperiod"`
	MinSalary   int64  `json:"minsalary"`
	MaxSalary   int64  `json:"maxsalary"`
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

	job, err := repository.EmployerCreateJob(publicID, jobDetails.Title, jobDetails.JobType, jobDetails.Category, jobDetails.Description, jobDetails.VisibleDate, jobDetails.PayPeriod, jobDetails.Remote, jobDetails.MinSalary, jobDetails.MaxSalary)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	response.SendJSON(w, job)

}
