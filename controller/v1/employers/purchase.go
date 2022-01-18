package employers

import (
	"encoding/json"
	"log"
	"net/http"

	"autumnomous-jobs-employer-api/shared/repository/jobpackages"
	"autumnomous-jobs-employer-api/shared/repository/jobs"
	"autumnomous-jobs-employer-api/shared/response"
	"autumnomous-jobs-employer-api/shared/services/security/jwt"
)

type purchaseJobPackageDetails struct {
	JobPackage string `json:"jobpackage"`
}

func PurchaseJobPackage(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		response.SendJSONMessage(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	var jobDetails purchaseJobPackageDetails

	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&jobDetails)

	if jobDetails.JobPackage == "" {

		response.SendJSONMessage(w, http.StatusBadRequest, response.MissingRequiredValue)
		return
	}

	publicID := jwt.GetUserClaim(r)

	repository := jobpackages.NewJobPackageRegistry().GetJobPackageRepository()

	// TODO: stripe payment

	jobPackage, err := repository.GetJobPackage(jobDetails.JobPackage)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	jobRepository := jobs.NewJobRegistry().GetJobRepository()

	for i := 0; i < jobPackage.NumberOfJobs; i++ {
		// TODO: jobDetails.PostEndDatetime = jobDetails.PostStartDatetime + 30 days

		_, err := jobRepository.EmployerCreateJob(publicID, "Edit", "", "", "", "", "", false, 0, 0)

		if err != nil {
			log.Println(err)
			response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
			return
		}

	}

	response.SendJSONMessage(w, http.StatusOK, "success")

}
