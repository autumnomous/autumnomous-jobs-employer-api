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
	// Bio          string `json:"bio"`
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

	employer, err := repository.UpdateEmployerAccount(publicID, data.FirstName, data.LastName, data.Email, data.PhoneNumber, data.MobileNumber, data.Role)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	response.SendJSON(w, employer)
}
