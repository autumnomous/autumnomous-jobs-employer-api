package acl

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"strings"

	"autumnomous.com/bit-jobs-api/shared/response"
	"autumnomous.com/bit-jobs-api/shared/services/security/jwt"
)

// DisallowAuth does not allow authenticated users to access the page
func DisallowAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Authorization") != "" { // If user is authenticated, don't allow them to access the page
			// token, err := jwt.ParseToken(r.Header["Authorization"][0])

			// if err != nil {
			// 	controller.Error500(w, err, []byte("An Error Occurred"))
			// }

			// h.ServeHTTP(w, r)
			response.SendJSONMessage(w, http.StatusInternalServerError, response.Unauthorized)
		} else { // If user is not authenticated, don't allow them to access the page
			response.SendJSONMessage(w, http.StatusInternalServerError, response.Unauthorized)
		}

	})
}

// DisallowAnon does not allow anonymous users to access the page
func DisallowAnon(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Authorization") != "" {
			claims, err := jwt.GetClaims(r)

			if err != nil {
				log.Println(err)
				response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
			}

			if !jwt.ValidateToken(claims) {
				response.SendJSONMessage(w, http.StatusUnauthorized, response.Unauthorized)
			}

			h.ServeHTTP(w, r)
		} else { // If user is not authenticated, don't allow them to access the page
			response.SendJSONMessage(w, http.StatusInternalServerError, response.Unauthorized)
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
