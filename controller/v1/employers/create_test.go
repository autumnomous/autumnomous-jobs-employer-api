package employers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"autumnomous.com/bit-jobs-api/controller/v1/employers"
	"autumnomous.com/bit-jobs-api/shared/services/security/encryption"
	"autumnomous.com/bit-jobs-api/shared/services/security/jwt"
	"autumnomous.com/bit-jobs-api/shared/testhelper"

	"github.com/stretchr/testify/assert"
)

func Test_Employer_CreateJob_IncorrectMethod(t *testing.T) {

	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(employers.CreateJob))

	defer ts.Close()

	methods := []string{"GET", "PUT", "DELETE"}

	for _, method := range methods {

		request, err := http.NewRequest(method, ts.URL, nil)

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

func Test_Employer_CreateJob_IncorrectData(t *testing.T) {

	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(employers.CreateJob))

	defer ts.Close()

	tests := map[string]map[string]string{
		"NoJobTitle": {
			"jobtitle":          "",
			"jobstreetaddress":  "",
			"jobcity":           "",
			"jobzipcode":        "",
			"jobtype":           "",
			"jobremotefriendly": "",
			"jobdescription":    "",
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

		response, err := httpClient.Do(request)

		assert.Equal(int(http.StatusBadRequest), response.StatusCode)
		assert.Nil(err)

	}

}

func Test_Employer_CreateJob_CorrectData(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(employers.CreateJob))

	defer ts.Close()

	data := map[string]string{
		"jobtitle":          fmt.Sprintf("New Job %s", encryption.GeneratePassword(9)),
		"jobstreetaddress":  "123 Street Avenue",
		"jobcity":           "City",
		"jobzipcode":        "00000",
		"jobtype":           "full-time",
		"jobremotefriendly": "yes",
		"jobdescription":    "This is a new job",
	}

	requestBody, err := json.Marshal(data)

	if err != nil {
		t.Fatal()
	}

	request, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(requestBody))

	if err != nil {
		t.Fatal()
	}
	employer := testhelper.Helper_RandomEmployer(t)

	token, err := jwt.GenerateToken(employer.PublicID)

	if err != nil {
		t.Fatal()
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	httpClient := &http.Client{}

	response, err := httpClient.Do(request)

	if err != nil {
		t.Fatal()
	}

	var result map[string]interface{}

	decoder := json.NewDecoder(response.Body)

	err = decoder.Decode(&result)

	if err != nil {
		t.Fatal()
	}

	assert.Contains(result, "publicid")
	assert.Equal(int(http.StatusOK), response.StatusCode)

}
