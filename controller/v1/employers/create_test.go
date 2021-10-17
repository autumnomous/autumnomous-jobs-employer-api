package employers_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bit-jobs-api/controller/v1/employers"
	"bit-jobs-api/shared/services/security/encryption"
	"bit-jobs-api/shared/services/security/jwt"
	"bit-jobs-api/shared/testhelper"

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

	data := map[string]interface{}{
		"title":             fmt.Sprintf("New Job %s", encryption.GeneratePassword(9)),
		"jobtype":           "full-time",
		"category":          "full-stack",
		"description":       "This is a new job",
		"remote":            false,
		"poststartdatetime": time.Now().Format(time.RFC3339),
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

	token = base64.StdEncoding.EncodeToString([]byte(token))
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
