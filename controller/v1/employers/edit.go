package employers

import (
	"encoding/json"
	"log"
	"net/http"

	"bit-jobs-api/shared/repository/employers"
	"bit-jobs-api/shared/response"
	"bit-jobs-api/shared/services/security/jwt"
)

type editJobDetails struct {
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
	PublicID          string `json:"publicid"`
}

func EditJob(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		response.SendJSONMessage(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var details editJobDetails

	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&details)

	tokenClaims, err := jwt.GetStrClaims(r)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusBadRequest, response.FriendlyError)
		return
	}

	publicID := tokenClaims["user"]

	repository := employers.NewEmployerRegistry().GetEmployerRepository()

	job, err := repository.EditJob(publicID, details.PublicID, details.Title, details.StreetAddress, details.City, details.ZipCode, details.Tags, details.Description, details.PostStartDatetime, details.PostEndDatetime, details.PayPeriod, details.MinSalary, details.MaxSalary)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	response.SendJSON(w, job)
}
