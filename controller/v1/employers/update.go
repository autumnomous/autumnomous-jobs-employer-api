package employers

import (
	"encoding/json"
	"log"
	"net/http"

	"bit-jobs-api/shared/repository/employers"
	"bit-jobs-api/shared/response"
	"bit-jobs-api/shared/services/security/jwt"
	// stripe "github.com/stripe/stripe-go/v72"
)

type updatePasswordCredentials struct {
	Password    string `json:"password"`
	NewPassword string `json:"newpassword"`
}

type updateAccountData struct {
	FirstName    string `json:"firstname"`
	LastName     string `json:"lastname"`
	Email        string `json:"email"`
	PhoneNumber  string `json:"phonenumber"`
	MobileNumber string `json:"mobilenumber"`
	Role         string `json:"role"`
	Facebook     string `json:"facebook"`
	Twitter      string `json:"twitter"`
	Instagram    string `json:"instagram"`
	// Bio          string `json:"bio"`
}

type updateCompanyData struct {
	Name         string  `json:"name"`
	Location     string  `json:"location"`
	Longitude    float64 `json:"longitude"`
	Latitude     float64 `json:"latitude"`
	URL          string  `json:"url"`
	Facebook     string  `json:"facebook"`
	Twitter      string  `json:"twitter"`
	Instagram    string  `json:"instagram"`
	Description  string  `json:"description"`
	Logo         string  `json:"logo"`
	ExtraDetails string  `json:"extradetails"`
}

type updatePaymentMethodData struct {
	PaymentMethod string `json:"paymentmethod"`
}

type updatePaymentDetailsData struct {
	PaymentDetails string `json:"paymentdetails"`
}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		response.SendJSONMessage(w, http.StatusMethodNotAllowed, response.FriendlyError)
		return
	}

	var credentials updatePasswordCredentials
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&credentials)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	if credentials.Password == "" || credentials.NewPassword == "" {
		response.SendJSONMessage(w, http.StatusBadRequest, response.FriendlyError)
		return
	}

	publicID := jwt.GetUserClaim(r)

	if publicID == "" {
		response.SendJSONMessage(w, http.StatusBadRequest, response.FriendlyError)
		return
	}

	repository := employers.NewEmployerRegistry().GetEmployerRepository()

	updated, err := repository.UpdateEmployerPassword(publicID, credentials.Password, credentials.NewPassword)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	if updated {
		response.SendJSONMessage(w, http.StatusOK, response.Success)
		return
	} else {
		response.SendJSONMessage(w, http.StatusBadRequest, response.FriendlyError)
		return
	}

}

func UpdateAccount(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		response.SendJSONMessage(w, http.StatusMethodNotAllowed, response.FriendlyError)
		return
	}

	publicID := jwt.GetUserClaim(r)

	if publicID == "" {
		response.SendJSONMessage(w, http.StatusBadRequest, response.FriendlyError)
		return
	}

	var data updateAccountData
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	repository := employers.NewEmployerRegistry().GetEmployerRepository()

	employer, err := repository.UpdateEmployerAccount(publicID, data.FirstName, data.LastName, data.Email, data.PhoneNumber, data.MobileNumber, data.Role, data.Facebook, data.Twitter, data.Instagram)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	response.SendJSON(w, employer)
}

func UpdateCompany(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		response.SendJSONMessage(w, http.StatusMethodNotAllowed, response.FriendlyError)
		return
	}

	publicID := jwt.GetUserClaim(r)

	if publicID == "" {
		response.SendJSONMessage(w, http.StatusBadRequest, response.FriendlyError)
		return
	}

	var data updateCompanyData
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	repository := employers.NewEmployerRegistry().GetEmployerRepository()

	company, err := repository.UpdateEmployerCompany(publicID, data.Name, data.Location, data.URL, data.Facebook, data.Twitter, data.Instagram, data.Description, data.Logo, data.ExtraDetails, data.Longitude, data.Latitude)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	response.SendJSON(w, company)

}

func UpdatePaymentMethod(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		response.SendJSONMessage(w, http.StatusMethodNotAllowed, response.FriendlyError)
		return
	}

	publicID := jwt.GetUserClaim(r)

	if publicID == "" {
		response.SendJSONMessage(w, http.StatusBadRequest, response.FriendlyError)
		return
	}

	var method updatePaymentMethodData
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&method)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	repository := employers.NewEmployerRegistry().GetEmployerRepository()

	err = repository.UpdateEmployerPaymentMethod(publicID, method.PaymentMethod)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	response.SendJSON(w, nil)

}

func UpdatePaymentDetails(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		response.SendJSONMessage(w, http.StatusMethodNotAllowed, response.FriendlyError)
		return
	}

	publicID := jwt.GetUserClaim(r)

	if publicID == "" {
		response.SendJSONMessage(w, http.StatusBadRequest, response.FriendlyError)
		return
	}

	var details updatePaymentDetailsData
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&details)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	repository := employers.NewEmployerRegistry().GetEmployerRepository()

	err = repository.UpdateEmployerPaymentDetails(publicID, details.PaymentDetails)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	response.SendJSON(w, nil)

}
