package employers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"bit-jobs-api/controller/v1/employers"
	"bit-jobs-api/shared/services/security/encryption"
	"bit-jobs-api/shared/testhelper"

	"github.com/stretchr/testify/assert"
)

func init() {
	testhelper.Init()
}

func Test_ApplicantLogin_Success(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(employers.Login))

	defer ts.Close()

	data := map[string]string{
		"email":    fmt.Sprintf("lavernecox-%s@amazing.com", string(encryption.GeneratePassword(9))),
		"password": string(encryption.GeneratePassword(9)),
	}

	requestBody, err := json.Marshal(data)

	if err != nil {
		t.Fatal()
	}

	request, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(requestBody))

	if err != nil {
		t.Fatal()
	}

	employers.AuthenticationFunction = func(email, password string) (bool, string, string, error) {
		return true, "", "", nil
	}

	defer func() {
		employers.AuthenticationFunction = employers.AuthenticatePassword
	}()

	httpClient := &http.Client{}

	response, err := httpClient.Do(request)

	if err != nil {
		t.Fatal()
	}

	var result map[string]interface{}
	decoder := json.NewDecoder(response.Body)

	err = decoder.Decode(&result)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(int(http.StatusOK), response.StatusCode)
	assert.NotNil(result)
}

func Test_ApplicantLogin_NoDataReceived(t *testing.T) {

	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(employers.Login))

	defer ts.Close()

	data := map[string]string{
		"email":    "",
		"password": "",
	}

	requestBody, err := json.Marshal(data)

	if err != nil {
		t.Fatal()
	}

	request, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(requestBody))

	if err != nil {
		t.Fatal()
	}

	httpClient := &http.Client{}

	response, err := httpClient.Do(request)

	if err != nil {
		t.Fatal()
	}

	var result map[string]interface{}

	decoder := json.NewDecoder(response.Body)
	decoder.Decode(&result)

	assert.Equal(int(http.StatusBadRequest), response.StatusCode)

}

func Test_ApplicantLogin_IncorrectMethod(t *testing.T) {

	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(employers.Login))

	defer ts.Close()

	data := map[string]string{
		"email":    "",
		"password": "",
	}

	requestBody, err := json.Marshal(data)

	if err != nil {
		t.Fatal()
	}

	methods := []string{
		"GET",
		"PUT",
		"DELETE",
	}

	for _, method := range methods {

		request, err := http.NewRequest(method, ts.URL, bytes.NewBuffer(requestBody))

		if err != nil {
			t.Fatal()
		}

		httpClient := &http.Client{}

		response, err := httpClient.Do(request)

		if err != nil {
			t.Fatal()
		}

		assert.Equal(int(http.StatusMethodNotAllowed), response.StatusCode)

	}

}

func Test_ApplicantLogin_IncorrectPassword(t *testing.T) {

	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(employers.Login))

	data := map[string]string{
		"firstname":          "First",
		"lastname":           "Last",
		"email":              fmt.Sprintf("email-%s@test.com", string(encryption.GeneratePassword(9))),
		"password":           string(encryption.GeneratePassword(9)),
		"applicant-password": string(encryption.GeneratePassword(9)),
	}

	applicant := testhelper.TestUser{
		FirstName: data["firstname"],
		LastName:  data["lastname"],
		Email:     data["email"],
	}

	hashedPassword, err := encryption.HashPassword([]byte(data["applicant-password"]))

	if err != nil {
		t.Fatal()
	}

	applicant.HashedPassword = hashedPassword

	requestBody, err := json.Marshal(data)

	if err != nil {
		t.Fatal()
	}

	request, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(requestBody))

	if err != nil {
		t.Fatal()
	}

	httpClient := &http.Client{}

	response, err := httpClient.Do(request)

	if err != nil {
		t.Fatal()
	}

	assert.Equal(int(http.StatusUnauthorized), response.StatusCode)

}
