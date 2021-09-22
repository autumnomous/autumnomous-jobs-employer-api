package acl

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"strings"

	"bit-jobs-api/shared/repository/employers"
	"bit-jobs-api/shared/response"
	jwt "bit-jobs-api/shared/services/security/jwt"
)

func ValidateMyJWT(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		s := strings.SplitN(req.Header.Get("Authorization"), " ", 2)
		if len(s) == 2 {
			b, err := base64.StdEncoding.DecodeString(s[1])

			if err == nil {
				bearerToken := strings.SplitN(string(b), ":", 2)

				if len(bearerToken) == 1 {

					data, err := jwt.ParseToken(bearerToken[0])

					if err != nil {
						log.Println(err)
						response.SendJSONMessage(w, http.StatusBadRequest, "Couldn't parse JWT")
					}

					if data == nil {
						log.Print("data cannot be nil")
						response.SendJSONMessage(w, http.StatusUnauthorized, "Unauthorized")
						//TODO: respond with an httpstatus of not authorized
					} else {
						userId := data.CustomClaims["user"]

						if userId == "" {
							//TODO: respond with an httpstatus of not authorized
							response.SendJSONMessage(w, http.StatusUnauthorized, "Unauthorized")
						} else {
							//TODO: get valid/active user by userid
							repository := employers.NewEmployerRegistry().GetEmployerRepository()
							_, err := repository.GetEmployer(userId)

							//validate user below

							if err != nil {
								log.Println("bad keys")
								//todo: respond with not authorized
								response.SendJSONMessage(w, http.StatusUnauthorized, "Bad keys")
							} else {

								h.ServeHTTP(w, req)
							}
						}
					}
				}
			}
		} else {
			log.Println("bearer token bad")
			//respond with invalid token
		}
	})
}

//AllowAPIKey allows authentication if an API key is present
func AllowAPIKey(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Authorization") != "" {

			auth := strings.Split(r.Header["Authorization"][0], " ")[1]
			authKey, err := base64.StdEncoding.DecodeString(string(auth))

			if err != nil {
				log.Println(err)
				response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
				return
			}

			apiKey := os.Getenv("API_KEY")
			if strings.Compare(strings.TrimSpace(string(authKey)), apiKey) == 0 && apiKey != "" {
				h.ServeHTTP(w, r)
			} else {
				response.SendJSONMessage(w, http.StatusUnauthorized, response.InvalidAPIKey)
			}

			// check for key in db

		} else { // If user is not authenticated, don't allow them to access the page
			response.SendJSONMessage(w, http.StatusUnauthorized, response.InvalidAPIKey)
		}

	})
}
