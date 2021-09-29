package employers_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"bit-jobs-api/controller/v1/employers"
	"bit-jobs-api/shared/services/security/jwt"
	"bit-jobs-api/shared/testhelper"

	"github.com/stretchr/testify/assert"
)

func Test_Employer_DeleteJob_IncorrectMethods(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(employers.DeleteJob))

	defer ts.Close()

	methods := []string{"POST", "PUT", "GET"}

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

func Test_Employer_DeleteJob_NoData(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(employers.DeleteJob))

	defer ts.Close()

	requestBody, err := json.Marshal(map[string]string{
		"publicid": "",
	})

	if err != nil {
		t.Fatal()
	}

	request, err := http.NewRequest("DELETE", ts.URL, bytes.NewBuffer(requestBody))

	if err != nil {
		t.Fatal()
	}

	httpClient := &http.Client{}

	response, err := httpClient.Do(request)

	if err != nil {
		t.Fatal()
	}

	assert.Equal(int(http.StatusBadRequest), response.StatusCode)

}

func Test_Employer_DeleteJob_Correct(t *testing.T) {

	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(employers.DeleteJob))

	defer ts.Close()

	employer := testhelper.Helper_RandomEmployer(t)

	job := testhelper.Helper_RandomJob(employer, t)

	requestBody, err := json.Marshal(map[string]string{
		"publicid": job.PublicID,
	})

	if err != nil {
		t.Fatal()
	}

	token, err := jwt.GenerateToken(employer.PublicID)

	if err != nil {
		t.Fatal()
	}

	request, err := http.NewRequest("DELETE", ts.URL, bytes.NewBuffer(requestBody))

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

	assert.Equal(int(http.StatusOK), response.StatusCode)

}
