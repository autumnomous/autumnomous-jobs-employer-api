package jwt

import (
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt"
)

// Token represents a JWT token
type Token struct {
	JWToken string `json:"token"`
}

// GenerateToken generates a JWT Token
/*func GenerateToken(userID string) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user"] = userID
	// claims["exp"] = time.Now().Add(time.Minute * 60 * 24 * 7).Unix()
	// claims["iat"] = time.Now()
	// claims["nbf"] = time.Now()

	tokenStr, err := token.SignedString(signingKey)

	return tokenStr, err
} */

type JWTData struct {
	// Standard claims are the standard jwt claims from the IETF standard
	// https://tools.ietf.org/html/rfc7519
	jwt.StandardClaims
	CustomClaims map[string]string `json:"custom,omitempty"`
}

func GenerateToken(userId string) (string, error) {

	claims := JWTData{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 168).Unix(), // week
			IssuedAt:  time.Now().Unix(),
		},
		CustomClaims: map[string]string{
			"user": userId,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("KNIT_SIGNING_KEY")))

	if err != nil {
		return "", err
	}

	return tokenString, err
}

// ParseToken parses a given JWT token
func ParseToken(inputTokenString string) (*JWTData, error) {

	claims, err := jwt.ParseWithClaims(inputTokenString, &JWTData{}, func(token *jwt.Token) (interface{}, error) {
		if jwt.SigningMethodHS256 != token.Method && !token.Valid {
			return nil, errors.New("invalid signing algorithm")
		}
		return []byte(os.Getenv("KNIT_SIGNING_KEY")), nil
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return claims.Claims.(*JWTData), nil

}

func GetUserClaim(r *http.Request) string {

	if r.Header.Get("Authorization") == "" {
		log.Println(errors.New("problem with bearer token"))
		return ""
	}

	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

	authKey, err := base64.StdEncoding.DecodeString(auth[1])

	if err != nil {
		log.Println(err)
		return ""
	}

	tokenClaims, err := ParseToken(string(authKey))

	if err != nil {
		log.Println(err)
		return ""
	}

	return tokenClaims.CustomClaims["user"]

}
