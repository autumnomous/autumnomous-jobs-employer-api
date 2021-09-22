package jwt

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt"
)

var signingKey = []byte(os.Getenv("KNIT_SIGNING_KEY"))

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

func  GenerateToken(userId string) string {

	claims := JWTData{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
		CustomClaims: map[string]string{
			"user": userId,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("KNIT_SIGNING_KEY")))

	if err != nil {
		return ""
	}

	return tokenString
}

// ParseToken parses a given JWT token
func ParseToken(inputTokenString string) (*jwt.Token, jwt.MapClaims, error) {

	token, err := jwt.Parse(inputTokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return signingKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return token, claims, nil
	} else {
		return nil, nil, err
	}

}

// ValidateToken confirms that a token is valid
func ValidateToken(claims jwt.MapClaims) bool {

	if claims.Valid() == nil {
		return true
	}

	return false
}

// GetClaims returns the jwt claims as a map
func GetClaims(r *http.Request) (jwt.MapClaims, error) {
	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)[1]
	// log.Println(auth)
	// authKey, err := base64.StdEncoding.DecodeString(auth)

	// if err != nil {
	// 	log.Println(err)
	// 	return nil, err
	// }

	_, claims, err := ParseToken(auth)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return claims, nil
}

// GetStrClaims returns the jwt claims as a map of strings
func GetStrClaims(r *http.Request) (map[string]string, error) {

	if r.Header.Get("Authorization") == "" {
		return nil, errors.New("a token is required for this action")
	}

	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)[1]
	// authKey, err := base64.StdEncoding.DecodeString(auth)
	// if err != nil {
	// 	log.Println(err)
	// 	return nil, err
	// }

	_, claims, err := ParseToken(auth)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	strClaims := make(map[string]string)

	strClaims["user"] = fmt.Sprintf("%v", claims["user"])
	// strClaims["exp"] = fmt.Sprintf("%v", claims["exp"])
	// strClaims["iat"] = fmt.Sprintf("%v", claims["iat"])
	// strClaims["nbf"] = fmt.Sprintf("%v", claims["nbf"])

	return strClaims, nil
}
