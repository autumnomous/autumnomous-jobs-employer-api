package employers_test

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"bit-jobs-api/controller/v1/employers"
	"bit-jobs-api/shared/repository/jobs"
	"bit-jobs-api/shared/services/security/jwt"
	"bit-jobs-api/shared/testhelper"

	"github.com/stretchr/testify/assert"
)

type JobsResponse struct {
	Jobs             []*jobs.Job `json:"jobs"`
	TotalPostsBought int         `json:"totalpostsbought"`
}

func init() {
	testhelper.Init()
}

func Test_Employer_GetJobs_IncorrectPublicID(t *testing.T) {

	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(employers.GetJobs))

	defer ts.Close()

	request, err := http.NewRequest("GET", ts.URL, nil)

	if err != nil {
		t.Fatal()
	}

	token, err := jwt.GenerateToken("")

	if err != nil {
		t.Fatal()
	}

	token = base64.StdEncoding.EncodeToString([]byte(token))

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

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
		log.Println(err)
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

	token = base64.StdEncoding.EncodeToString([]byte(token))

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	httpClient := &http.Client{}

	response, err := httpClient.Do(request)

	if err != nil {
		log.Println(err)
		t.Fatal()
	}

	decoder := json.NewDecoder(response.Body)
	var result JobsResponse

	decoder.Decode(&result)

	assert.Nil(err)
	assert.Equal(int(http.StatusOK), response.StatusCode)

	assert.Equal(3, len(result.Jobs))
	assert.Equal(0, result.TotalPostsBought)

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

	token = base64.StdEncoding.EncodeToString([]byte(token))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	httpClient := &http.Client{}

	response, err := httpClient.Do(request)

	decoder := json.NewDecoder(response.Body)
	var result []map[string]string

	decoder.Decode(&result)
	assert.Nil(err)
	assert.Equal(int(http.StatusOK), response.StatusCode)
	assert.GreaterOrEqual(len(result), 1)

}

func Test_Employer_GetCompany_Correct(t *testing.T) {

	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(employers.GetEmployerCompany))

	defer ts.Close()

	request, err := http.NewRequest("GET", ts.URL, nil)

	if err != nil {
		log.Println(err)
		t.Fatal()
	}

	employer := testhelper.Helper_RandomEmployer(t)

	company := testhelper.Helper_RandomCompany(t)

	err = testhelper.Helper_SetEmployerCompany(employer.PublicID, company.PublicID)

	if err != nil {
		log.Println(err)
		t.Fatal()
	}

	token, err := jwt.GenerateToken(employer.PublicID)

	if err != nil {
		t.Fatal()
	}

	token = base64.StdEncoding.EncodeToString([]byte(token))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	httpClient := &http.Client{}

	response, err := httpClient.Do(request)

	decoder := json.NewDecoder(response.Body)
	var result []map[string]string

	decoder.Decode(&result)
	assert.Nil(err)
	assert.Equal(int(http.StatusOK), response.StatusCode)

}
