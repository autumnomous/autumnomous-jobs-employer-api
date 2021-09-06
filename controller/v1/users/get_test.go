package user_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	users "autumnomous.com/bit-jobs-api/controller/v1/users"
	"autumnomous.com/bit-jobs-api/shared/services/security/encryption"
	"autumnomous.com/bit-jobs-api/shared/services/security/jwt"
	"autumnomous.com/bit-jobs-api/shared/testhelper"
	"github.com/stretchr/testify/assert"
)

func init() {
	testhelper.Init()
}

func Test_GetUser_Applicant_Success(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(users.GetUser))

	client := &testhelper.TestUser{FirstName: "First", LastName: "Last", Email: fmt.Sprintf("email-%s@site.com", encryption.GeneratePassword(9))}

	hashedPassword, err := encryption.HashPassword([]byte(encryption.GeneratePassword(9)))

	if err != nil {
		t.Fatal()
	}

	client.HashedPassword = hashedPassword

	client = testhelper.Helper_CreateApplicant(client, t)

	token, err := jwt.GenerateToken(client.PublicID)

	if err != nil {
		t.Fatal()
	}

	request, err := http.NewRequest("GET", ts.URL, nil)

	if err != nil {
		t.Fatal()
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	httpClient := http.Client{}

	response, err := httpClient.Do(request)

	if err != nil {
		t.Fatal()
	}

	var result map[string]string
	decoder := json.NewDecoder(response.Body)

	err = decoder.Decode(&result)

	if err != nil {
		t.Fatal(err)
	}

	assert.Contains(result, "accounttype")
	assert.Contains(result["accounttype"], "applicant")
	assert.Equal(http.StatusOK, response.StatusCode)

}

func Test_GetUser_Employer_Success(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(users.GetUser))

	employer := &testhelper.TestEmployer{FirstName: "First", LastName: "Last", Email: fmt.Sprintf("email-%s@site.com", encryption.GeneratePassword(9))}

	hashedPassword, err := encryption.HashPassword([]byte(encryption.GeneratePassword(9)))

	if err != nil {
		t.Fatal()
	}

	employer.HashedPassword = hashedPassword

	employer = testhelper.Helper_CreateEmployer(employer, t)

	token, err := jwt.GenerateToken(employer.PublicID)

	if err != nil {
		t.Fatal()
	}

	request, err := http.NewRequest("GET", ts.URL, nil)

	if err != nil {
		t.Fatal()
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	httpClient := http.Client{}

	response, err := httpClient.Do(request)

	if err != nil {
		t.Fatal()
	}

	var result map[string]string
	decoder := json.NewDecoder(response.Body)

	err = decoder.Decode(&result)

	if err != nil {
		t.Fatal(err)
	}

	assert.Contains(result, "accounttype")
	assert.Contains(result["accounttype"], "employer")
	assert.Equal(http.StatusOK, response.StatusCode)

}
