package employers_test

import (
	"autumnomous-jobs-employer-api/controller/v1/employers"
	"autumnomous-jobs-employer-api/shared/services/security/jwt"
	"autumnomous-jobs-employer-api/shared/testhelper"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Employer_PurchaseJobPackage_IncorrectMethod(t *testing.T) {

	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(employers.PurchaseJobPackage))

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

func Test_Employer_PurchaseJobPackage_IncorrectData(t *testing.T) {

	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(employers.PurchaseJobPackage))

	defer ts.Close()

	tests := map[string]map[string]string{
		"NoJobPackage": {
			"jobpackage": "",
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

func Test_Employer_PurchaseJobPackage_CorrectData(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(employers.PurchaseJobPackage))

	defer ts.Close()

	jobpackage := testhelper.Helper_RandomJobPackage(t)
	data := map[string]string{
		"jobpackage": jobpackage.TypeID,
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

	assert.Equal(int(http.StatusOK), response.StatusCode)

}
