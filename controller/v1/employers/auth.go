package employers

import (
	"encoding/json"
	"log"
	"net/http"

	"bit-jobs-api/shared/repository/employers"
	"bit-jobs-api/shared/response"
	"bit-jobs-api/shared/services/security/jwt"
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

	match, initialPasswordChanged, publicID, err := AuthenticationFunction(credentials.Email, credentials.Password)

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

		token := map[string]interface{}{
			"token":                  tokenStr,
			"initialpasswordchanged": initialPasswordChanged,
		}
		response.SendJSON(w, token)
		return
	} else {
		response.SendJSONMessage(w, http.StatusUnauthorized, "Login failed")
		return
	}

}

func AuthenticatePassword(email, password string) (bool, bool, string, error) {

	repository := employers.NewEmployerRegistry().GetEmployerRepository()

	match, initialPasswordChanged, publicID, err := repository.AuthenticateEmployerPassword(email, password)

	if err != nil {
		log.Println(err)
		return false, initialPasswordChanged, "", err
	}

	return match, initialPasswordChanged, publicID, nil

}

var AuthenticationFunction = AuthenticatePassword
