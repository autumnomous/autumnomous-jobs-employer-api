package employers_test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"autumnomous.com/bit-jobs-api/controller/v1/employers"
	"autumnomous.com/bit-jobs-api/shared/services/security/jwt"
	"autumnomous.com/bit-jobs-api/shared/testhelper"
	"github.com/stretchr/testify/assert"
)

func Test_Employer_GetJobs_IncorrectPublicID(t *testing.T) {

	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(employers.GetJobs))

	defer ts.Close()

	request, err := http.NewRequest("GET", ts.URL, nil)

	if err != nil {
		t.Fatal()
	}

	httpClient := &http.Client{}

	response, err := httpClient.Do(request)

	assert.Nil(err)
	assert.Equal(int(http.StatusBadRequest), response.StatusCode)

}

func Test_Employer_GetJobs_IncorrectMethod(t *testing.T) {

	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(employers.GetJobs))

	defer ts.Close()

	methods := []string{"POST", "DELETE", "PUT"}

	for _, method := range methods {
		request, err := http.NewRequest(method, ts.URL, nil)

		if err != nil {
			t.Fatal()
		}

		httpClient := &http.Client{}

		response, err := httpClient.Do(request)

		assert.Nil(err)
		assert.Equal(int(http.StatusMethodNotAllowed), response.StatusCode)
	}

}

func Test_Employer_GetJobs_Correct(t *testing.T) {

	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(employers.GetJobs))

	defer ts.Close()

	request, err := http.NewRequest("GET", ts.URL, nil)

	if err != nil {
		t.Fatal()
	}

	employer := testhelper.Helper_RandomEmployer(t)

	testhelper.Helper_RandomJob(employer, t)
	testhelper.Helper_RandomJob(employer, t)
	testhelper.Helper_RandomJob(employer, t)

	token, err := jwt.GenerateToken(employer.PublicID)

	if err != nil {
		t.Fatal()
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	httpClient := &http.Client{}

	response, err := httpClient.Do(request)

	decoder := json.NewDecoder(response.Body)
	var result []map[string]string

	decoder.Decode(&result)

	assert.Nil(err)
	assert.Equal(int(http.StatusOK), response.StatusCode)
	assert.Equal(len(result), 3)

}

func Test_Employer_GetJobPackages_Correct(t *testing.T) {

	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(employers.GetActiveJobPackages))

	defer ts.Close()

	request, err := http.NewRequest("GET", ts.URL, nil)

	if err != nil {
		t.Fatal()
	}

	employer := testhelper.Helper_RandomEmployer(t)

	testhelper.Helper_RandomJobPackage(t)

	token, err := jwt.GenerateToken(employer.PublicID)

	if err != nil {
		t.Fatal()
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	httpClient := &http.Client{}

	response, err := httpClient.Do(request)

	decoder := json.NewDecoder(response.Body)
	var result []map[string]string

	decoder.Decode(&result)
	log.Println(result)
	assert.Nil(err)
	assert.Equal(int(http.StatusOK), response.StatusCode)
	assert.GreaterOrEqual(len(result), 1)

}
