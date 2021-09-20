package employers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"bit-jobs-api/controller/v1/employers"
	"bit-jobs-api/shared/services/security/encryption"
	"bit-jobs-api/shared/services/security/jwt"
	"bit-jobs-api/shared/testhelper"

	"github.com/stretchr/testify/assert"
)

func Test_ApplicantUpdate_UpdatePassword_Success(t *testing.T) {

	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(employers.UpdatePassword))

	data := map[string]string{
		"password":    string(encryption.GeneratePassword(9)),
		"newpassword": string(encryption.GeneratePassword(9)),
		"token":       "",
	}

	applicant := &testhelper.TestUser{FirstName: "First", LastName: "Last", Email: fmt.Sprintf("email-%s@site.com", encryption.GeneratePassword(9))}

	hashedPassword, err := encryption.HashPassword([]byte(data["password"]))

	if err != nil {
		t.Fatal()
	}

	applicant.HashedPassword = hashedPassword

	applicant = testhelper.Helper_CreateApplicant(applicant, t)

	token, err := jwt.GenerateToken(applicant.PublicID)

	if err != nil {
		t.Fatal()
	}

	requestBody, err := json.Marshal(data)

	if err != nil {
		t.Fatal()
	}

	request, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(requestBody))

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
	assert.Equal(int(http.StatusOK), response.StatusCode)
}

func Test_ApplicantUpdate_UpdatePassword_IncorrectMethod(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(employers.UpdatePassword))

	defer ts.Close()

	data := map[string]string{
		"password":    "",
		"newpassword": "",
	}

	requestBody, err := json.Marshal(data)

	if err != nil {
		t.Fatal(err)
	}

	methods := []string{"GET", "PUT", "DELETE"}

	for _, method := range methods {

		request, err := http.NewRequest(method, ts.URL, bytes.NewBuffer(requestBody))

		if err != nil {
			t.Fatal(err)
		}

		client := &http.Client{}

		response, err := client.Do(request)

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(http.StatusMethodNotAllowed, response.StatusCode)
	}
}
func Test_ApplicantUpdate_UpdatePassword_IncorrectDataReceived_NoPassword(t *testing.T) {

	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(employers.UpdatePassword))

	data := map[string]string{
		"password":    "",
		"newpassword": string(encryption.GeneratePassword(9)),
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

	assert.Equal(int(http.StatusBadRequest), response.StatusCode)

}

func Test_ApplicantUpdate_UpdatePassword_IncorrectDataReceived_NoNewPassword(t *testing.T) {

	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(employers.UpdatePassword))

	data := map[string]string{
		"password":    string(encryption.GeneratePassword(9)),
		"newpassword": "",
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

	assert.Equal(int(http.StatusBadRequest), response.StatusCode)

}

func Test_ApplicantUpdate_UpdatePassword_IncorrectDataReceived_NoToken(t *testing.T) {

	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(employers.UpdatePassword))

	data := map[string]string{
		"password":    string(encryption.GeneratePassword(9)),
		"newpassword": string(encryption.GeneratePassword(9)),
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

	assert.Equal(int(http.StatusBadRequest), response.StatusCode)

}

func Test_ApplicantUpdate_UpdatePassword_IncorrectDataReceived_NoData(t *testing.T) {

	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(employers.UpdatePassword))

	data := map[string]string{
		"password":    "",
		"newpassword": "",
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

	assert.Equal(int(http.StatusBadRequest), response.StatusCode)

}

func Test_ApplicantUpdate_UpdateAccount_IncorrectDataReceived_NoToken(t *testing.T) {

	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(employers.UpdateAccount))

	request, err := http.NewRequest("POST", ts.URL, nil)

	if err != nil {
		log.Println(err)
		t.Fatal()
	}

	httpClient := http.Client{}
	response, err := httpClient.Do(request)

	assert.Nil(err)
	assert.Equal(int(http.StatusBadRequest), response.StatusCode)
}

func Test_ApplicantUpdate_UpdateAccount_CorrectDataReceived(t *testing.T) {

	assert := assert.New(t)

	ts := httptest.NewServer((http.HandlerFunc(employers.UpdateAccount)))

	applicant := &testhelper.TestUser{FirstName: "First", LastName: "Last", Email: fmt.Sprintf("email-%s@site.com", encryption.GeneratePassword(9)), Password: string(encryption.GeneratePassword(9))}

	hashedPassword, err := encryption.HashPassword([]byte(applicant.Password))

	if err != nil {
		t.Fatal()
	}

	applicant.HashedPassword = hashedPassword

	applicant = testhelper.Helper_CreateApplicant(applicant, t)

	token, err := jwt.GenerateToken(applicant.PublicID)

	if err != nil {
		t.Fatal()
	}

	tests := map[string]map[string]string{
		"NewBio": {
			"firstname": "First",
			"lastname":  "Last",
			"email":     applicant.Email,
			"bio":       "New Bio",
		},
		"New First Name": {
			"firstname": "NewFirst",
			"lastname":  "Last",
			"email":     applicant.Email,
			"bio":       "New Bio",
		},
		"New Last Name": {
			"firstname": "NewFirst",
			"lastname":  "NewLast",
			"email":     applicant.Email,
			"bio":       "New Bio",
		},
		"New Email": {
			"firstname": "NewFirst",
			"lastname":  "NewLast",
			"email":     fmt.Sprintf("new-email-%s@site.com", encryption.GeneratePassword(9)),
			"bio":       "New Bio",
		},
	}

	for _, test := range tests {

		data, err := json.Marshal(test)

		if err != nil {
			log.Println(err)
			t.Fatal()
		}

		request, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(data))
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

		db_user := testhelper.Helper_GetUser(applicant.PublicID, t)
		assert.Equal(int(http.StatusOK), response.StatusCode)
		assert.Equal(db_user.FirstName, result["firstname"])
		assert.Equal(db_user.LastName, result["lastname"])
		assert.Equal(db_user.Email, result["email"])
		assert.Equal(db_user.Biography, result["bio"])

	}

}
