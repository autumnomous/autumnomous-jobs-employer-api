package acl

import (
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"bit-jobs-api/shared/repository/employers"
	"bit-jobs-api/shared/response"
	jwtalso "bit-jobs-api/shared/services/security/jwt"

	jwt "github.com/golang-jwt/jwt"
)

type JWTData struct {
	// Standard claims are the standard jwt claims from the IETF standard
	// https://tools.ietf.org/html/rfc7519
	jwt.StandardClaims
	CustomClaims map[string]string `json:"custom,omitempty"`
}

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
			claims, err := jwtalso.GetClaims(r)

			if err != nil {
				log.Println(err)
				response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
			}

			if !jwtalso.ValidateToken(claims) {
				response.SendJSONMessage(w, http.StatusUnauthorized, response.Unauthorized)
			}

			h.ServeHTTP(w, r)
		} else { // If user is not authenticated, don't allow them to access the page
			response.SendJSONMessage(w, http.StatusInternalServerError, response.Unauthorized)
		}

	})
}

func ValidateMyJWT(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		s := strings.SplitN(req.Header.Get("Authorization"), " ", 2)
		log.Println(s)
		if len(s) == 2 {
			b, err := base64.StdEncoding.DecodeString(s[1])
			log.Println(s[1])
			log.Println(err)
			if err == nil {
				bearerToken := strings.SplitN(string(b), ":", 2)
				log.Println(len(bearerToken))
				if len(bearerToken) == 1 {

					claims, err := jwt.ParseWithClaims(bearerToken[0], &JWTData{}, func(token *jwt.Token) (interface{}, error) {
						if jwt.SigningMethodHS256 != token.Method && !token.Valid {
							return nil, errors.New("invalid signing algorithm")
						}
						return []byte(os.Getenv("KNIT_SIGNING_KEY")), nil
					})

					if err != nil {
						log.Println(err)
						response.SendJSONMessage(w, http.StatusBadRequest, "Couldn't parse JWT")
					}

					data := claims.Claims.(*JWTData)

					if data == nil {
						log.Print("data cannot be nil")
						response.SendJSONMessage(w, http.StatusUnauthorized, "Not authorized")
						//TODO: respond with an httpstatus of not authorized
					} else {
						userId := data.CustomClaims["user"]

						if userId == "" {
							//TODO: respond with an httpstatus of not authorized
							response.SendJSONMessage(w, http.StatusUnauthorized, "User ID is empty")
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
