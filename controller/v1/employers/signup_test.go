package employers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"bit-jobs-api/controller/v1/employers"
	"bit-jobs-api/shared/repository/employers/accountmanagement"
	"bit-jobs-api/shared/services/security/encryption"
	"bit-jobs-api/shared/testhelper"

	"github.com/stretchr/testify/assert"
)

func init() {
	testhelper.Init()
}

func Test_Employer_SignUp_IncorrectRequestMethod(t *testing.T) {

	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(employers.SignUp))

	defer ts.Close()

	data := map[string]string{
		"firstname": "",
		"lastname":  "",
		"email":     "",
	}

	requestBody, err := json.Marshal(data)

	if err != nil {
		t.Fatal()
	}

	methods := []string{"GET", "PUT", "DELETE"}

	for _, method := range methods {

		request, err := http.NewRequest(method, ts.URL, bytes.NewBuffer(requestBody))

		if err != nil {
			t.Fatal()
		}

		httpClient := &http.Client{}

		result, err := httpClient.Do(request)

		assert.Nil(err)
		assert.Equal(int(http.StatusMethodNotAllowed), result.StatusCode)
	}
}

func Test_Employer_SignUp_IncorrectData(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(employers.SignUp))

	defer ts.Close()

	tests := map[string]interface{}{
		"NoFirstName": map[string]string{
			"lastname": "hooks",
			"email":    fmt.Sprintf("bell-%s@power.com", string(encryption.GeneratePassword(9))),
		},
		"NoLastName": map[string]string{
			"firstname": "Assata",
			"email":     fmt.Sprintf("shakur-%s@power.com", string(encryption.GeneratePassword(9))),
		},
		"NoEmail": map[string]string{
			"firstname": "Fred",
			"lastname":  "Hampton",
		},
	}

	for _, test := range tests {

		requestBody, err := json.Marshal(test)

		if err != nil {
			t.Fatal()
		}

		request, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(requestBody))

		if err != nil {
			t.Fatal()
		}

		httpClient := &http.Client{}

		result, err := httpClient.Do(request)

		assert.Nil(err)
		assert.Equal(int(http.StatusBadRequest), result.StatusCode)

	}
}

func Test_Employer_SignUp_CorrectData(t *testing.T) {

	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(employers.SignUp))

	defer ts.Close()

	data, err := json.Marshal(map[string]string{
		"firstname": "First",
		"lastname":  "Last",
		"email":     fmt.Sprintf("email-%s@site.com", encryption.GeneratePassword(9)),
	})

	if err != nil {
		t.Fatal()
	}

	employers.SendWelcomeMessageFunction = func(domain, apiKey, password string, employer *accountmanagement.Employer) (string, error) {
		return "", nil
	}

	defer func() {
		employers.SendWelcomeMessageFunction = employers.SendWelcomeMessage
	}()

	request, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(data))

	if err != nil {
		t.Fatal()
	}

	httpClient := &http.Client{}

	response, err := httpClient.Do(request)

	if err != nil {
		t.Fatal(err)
	}

	var result string
	decoder := json.NewDecoder(response.Body)

	err = decoder.Decode(&result)

	if err != nil {
		t.Fatal(err)
	}

	assert.Nil(err)

}
