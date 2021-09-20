package user

import (
	"log"
	"net/http"

	"autumnomous.com/bit-jobs-api/shared/repository/users"
	"autumnomous.com/bit-jobs-api/shared/response"
	"autumnomous.com/bit-jobs-api/shared/services/security/jwt"
)

func GetUser(w http.ResponseWriter, r *http.Request) {

	tokenClaims, err := jwt.GetStrClaims(r)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusBadRequest, response.FriendlyError)
		return
	}

	publicID := tokenClaims["user"]

	// check if applicant
	user_repository := users.NewUserRegistry().GetUserRepository()

	user, err := user_repository.GetUser(publicID)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}
	response.SendJSON(w, user)
	// if applicant_user != nil {

	// }

	// check if employer

	// if both not nil

}
