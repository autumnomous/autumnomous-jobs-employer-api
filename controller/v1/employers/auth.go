package employers

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"

	"autumnomous-jobs-employer-api/shared/repository/employers"
	"autumnomous-jobs-employer-api/shared/response"
	"autumnomous-jobs-employer-api/shared/services/security/jwt"
)

type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		response.SendJSONMessage(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	decoder := json.NewDecoder(r.Body)
	credentials := LoginCredentials{}
	decoder.Decode(&credentials)

	if credentials.Email == "" || credentials.Password == "" {
		response.SendJSONMessage(w, http.StatusBadRequest, "Bad Request")
		return
	}

	match, registrationStep, publicID, err := AuthenticationFunction(credentials.Email, credentials.Password)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, "Internal Server Error")
	}

	if match {

		// use   to create JWT
		tokenStr, err := jwt.GenerateToken(publicID)

		if err != nil {
			log.Println(err)
			response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
			return
		}

		encodedTokenStr := base64.StdEncoding.EncodeToString([]byte(tokenStr))

		token := map[string]interface{}{
			"token":            encodedTokenStr,
			"registrationstep": registrationStep,
		}
		response.SendJSON(w, token)
		return
	} else {
		response.SendJSONMessage(w, http.StatusUnauthorized, "Login failed")
		return
	}

}

func AuthenticatePassword(email, password string) (bool, string, string, error) {

	repository := employers.NewEmployerRegistry().GetEmployerRepository()

	match, registrationStep, publicID, err := repository.AuthenticateEmployerPassword(email, password)

	if err != nil {
		log.Println(err)
		return false, "", "", err
	}

	return match, registrationStep, publicID, nil

}

var AuthenticationFunction = AuthenticatePassword
