package accountmanagement_test

import (
	"fmt"
	"log"
	"testing"

	"bit-jobs-api/shared/database"
	"bit-jobs-api/shared/repository/applicants"
	"bit-jobs-api/shared/repository/applicants/accountmanagement"
	"bit-jobs-api/shared/services/security/encryption"
	"bit-jobs-api/shared/testhelper"

	"github.com/stretchr/testify/assert"
)

func init() {
	testhelper.Init()
}

func Test_NewApplicantRepository(t *testing.T) {

	assert := assert.New(t)

	result := accountmanagement.NewApplicantRepository(database.DB)
	assert.Equal(database.DB, result.Database)
}

func Test_ApplicantRepository_CreateApplicant(t *testing.T) {
	assert := assert.New(t)

	data := map[string]string{
		"firstname": "Kendrick",
		"lastname":  "Lamar",
		"email":     fmt.Sprintf("klamar-%s@damn.com", string(encryption.GeneratePassword(9))),
		"password":  string(encryption.GeneratePassword(9)),
	}

	repository := accountmanagement.NewApplicantRepository(database.DB)

	result, err := repository.CreateApplicant(data["firstname"], data["lastname"], data["email"], data["password"])

	assert.Equal(result.FirstName, data["firstname"])
	assert.Equal(result.LastName, data["lastname"])
	assert.Equal(result.Email, data["email"])
	assert.NotNil(result.PublicID)
	assert.Nil(err)

}

func Test_ApplicantRepository_CreateApplicant_Fail_EmptyData(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]map[string]string{
		"NoFirstName": {
			"firstname": "",
			"lastname":  "Monáe",
			"email":     "jmonae@dirtycomputer.com",
			"password":  string(encryption.GeneratePassword(9)),
		},
		"NoLastName": {
			"firstname": "Jedenna",
			"lastname":  "",
			"email":     "jedenna@85toafrica.com",
			"password":  string(encryption.GeneratePassword(9)),
		},
		"NoEmail": {
			"firstname": "Kim",
			"lastname":  "Petras",
			"email":     "",
			"password":  string(encryption.GeneratePassword(9)),
		},
		"NoPassword": {
			"firstname": "Kim",
			"lastname":  "Petras",
			"email":     "",
			"password":  "",
		},
	}

	for _, test := range tests {

		repository := accountmanagement.NewApplicantRepository(database.DB)

		result, err := repository.CreateApplicant(test["firstname"], test["lastname"], test["email"], test["password"])

		assert.NotNil(err)
		assert.Nil(result)

	}

}

func Test_ApplicantRepository_AuthenticateApplicantPassword_NoDataReceived(t *testing.T) {

	assert := assert.New(t)

	repository := applicants.NewApplicantRegistry().GetApplicantRepository()

	data := map[string]string{
		"email":    "",
		"password": "",
	}

	result, _, publicid, err := repository.AuthenticateApplicantPassword(data["email"], data["password"])

	assert.False(result)
	assert.Nil(err)
	assert.Equal("", publicid)
}

func Test_ApplicantRepository_AuthenicateApplicantPassword_CorrectDataReceived(t *testing.T) {

	assert := assert.New(t)

	repository := applicants.NewApplicantRegistry().GetApplicantRepository()

	data := map[string]string{
		"firstname": "First",
		"lastname":  "Last",
		"email":     fmt.Sprintf("email-%s@site.com", string(encryption.GeneratePassword(9))),
		"password":  string(encryption.GeneratePassword(9)),
	}

	Applicant := &testhelper.TestUser{
		FirstName: data["firstname"],
		LastName:  data["lastname"],
		Email:     data["email"],
	}

	hashedPassword, err := encryption.HashPassword([]byte(data["password"]))

	Applicant.HashedPassword = hashedPassword
	if err != nil {
		t.Fatal()
	}

	Applicant = testhelper.Helper_CreateApplicant(Applicant, t)
	result, _, publicid, err := repository.AuthenticateApplicantPassword(data["email"], data["password"])

	assert.True(result)
	assert.Nil(err)
	assert.Equal(Applicant.PublicID, publicid)

}

func Test_ApplicantRepository_AuthenticateApplicantPassword_IncorrectDataReceived_Email(t *testing.T) {
	assert := assert.New(t)

	repository := applicants.NewApplicantRegistry().GetApplicantRepository()

	data := map[string]string{
		"firstname":       "First",
		"lastname":        "Last",
		"email":           fmt.Sprintf("email-%s@site.com", string(encryption.GeneratePassword(9))),
		"Applicant-email": fmt.Sprintf("email-%s@site.com", string(encryption.GeneratePassword(9))),
		"password":        string(encryption.GeneratePassword(9)),
	}

	Applicant := &testhelper.TestUser{
		FirstName: data["firstname"],
		LastName:  data["lastname"],
		Email:     data["Applicant-email"],
		Password:  data["password"],
	}

	hashedPassword, err := encryption.HashPassword([]byte(Applicant.Password))

	Applicant.HashedPassword = hashedPassword
	if err != nil {
		t.Fatal()
	}

	Applicant = testhelper.Helper_CreateApplicant(Applicant, t)

	result, _, publicID, err := repository.AuthenticateApplicantPassword(data["email"], Applicant.Password)

	assert.False(result)
	assert.Nil(err)
	assert.Equal("", publicID)

}

func Test_ApplicantRepository_AuthenticateApplicantPassword_IncorrectDataReceived_Password(t *testing.T) {
	assert := assert.New(t)

	repository := applicants.NewApplicantRegistry().GetApplicantRepository()

	data := map[string]string{
		"firstname":          "First",
		"lastname":           "Last",
		"email":              fmt.Sprintf("email-%s@site.com", string(encryption.GeneratePassword(9))),
		"Applicant-password": string(encryption.GeneratePassword(9)),
		"password":           string(encryption.GeneratePassword(9)),
	}

	applicant := &testhelper.TestUser{
		FirstName: data["firstname"],
		LastName:  data["lastname"],
		Email:     data["email"],
	}

	hashedPassword, err := encryption.HashPassword([]byte(data["Applicant-password"]))

	applicant.HashedPassword = hashedPassword
	if err != nil {
		t.Fatal()
	}

	_ = testhelper.Helper_CreateApplicant(applicant, t)

	result, initialpasswordchanged, publicid, err := repository.AuthenticateApplicantPassword(data["email"], data["password"])

	assert.False(result)
	assert.Nil(err)
	assert.False(initialpasswordchanged)
	assert.Equal("", publicid)

}

func Test_ApplicantRepository_UpdateApplicantPassword_NoDataReceived(t *testing.T) {
	assert := assert.New(t)

	repository := applicants.NewApplicantRegistry().GetApplicantRepository()

	data := map[string]string{
		"password":    "",
		"newpassword": "",
		"publicid":    "",
	}
	updated, err := repository.UpdateApplicantPassword(data["publicid"], data["password"], data["newpassword"])

	assert.False(updated)
	assert.Nil(err)

}

func Test_ApplicantRepository_UpdateApplicantPassword_IncorrectPublicID(t *testing.T) {
	assert := assert.New(t)
	repository := applicants.NewApplicantRegistry().GetApplicantRepository()

	data := map[string]string{
		"password":    string(encryption.GeneratePassword(9)),
		"newpassword": string(encryption.GeneratePassword(9)),
		"publicid":    string(encryption.GeneratePassword(10)),
	}

	updated, err := repository.UpdateApplicantPassword(data["publicid"], data["password"], data["newpassword"])

	assert.False(updated)
	assert.Nil(err)
}

func Test_ApplicantRepository_UpdateApplicantPassword_CorrectData(t *testing.T) {

	assert := assert.New(t)

	repository := applicants.NewApplicantRegistry().GetApplicantRepository()

	Applicant := &testhelper.TestUser{
		FirstName: "First",
		LastName:  "Last",
		Email:     fmt.Sprintf("email-%s@site.com", encryption.GeneratePassword(9)),
	}

	password := encryption.GeneratePassword(9)
	hashedPassword, err := encryption.HashPassword([]byte(password))

	if err != nil {
		t.Fatal()
	}

	Applicant.HashedPassword = hashedPassword

	Applicant = testhelper.Helper_CreateApplicant(Applicant, t)

	data := map[string]string{
		"password":    string(password),
		"newpassword": string(encryption.GeneratePassword(9)),
		"publicid":    Applicant.PublicID,
	}

	updated, err := repository.UpdateApplicantPassword(data["publicid"], data["password"], data["newpassword"])

	assert.True(updated)
	assert.Nil(err)
}

func Test_ApplicantRepository_UpdateApplicantAccount_CorrectData(t *testing.T) {

	assert := assert.New(t)

	repository := applicants.NewApplicantRegistry().GetApplicantRepository()

	applicant := &testhelper.TestUser{
		FirstName: "First",
		LastName:  "Last",
		Email:     fmt.Sprintf("email-%s@site.com", encryption.GeneratePassword(9)),
		Biography: "",
	}

	password := encryption.GeneratePassword(9)
	hashedPassword, err := encryption.HashPassword([]byte(password))

	if err != nil {
		log.Println(err)
		t.Fatal()
	}

	applicant.HashedPassword = hashedPassword

	applicant = testhelper.Helper_CreateApplicant(applicant, t)

	data := map[string]string{
		"firstname": "NewFirst",
		"lastname":  "NewLast",
		"email":     fmt.Sprintf("new-email-%s@site.com", encryption.GeneratePassword(9)),
		"bio":       "New Bio",
		"publicid":  applicant.PublicID,
	}

	updatedApplicant, err := repository.UpdateApplicantAccount(data["publicid"], data["firstname"], data["lastname"], data["email"], data["bio"])

	assert.NotNil(updatedApplicant)
	assert.Nil(err)
}

// func Test_ApplicantRepository_GetApplicantRegistrationStep_CorrectData_NewApplicant(t *testing.T) {

// 	assert := assert.New(t)

// 	repository := applicants.NewApplicantRegistry().GetApplicantRepository()

// 	Applicant := testhelper.Helper_RandomApplicant(t)

// 	result, err := repository.GetApplicantRegistrationStep(Applicant.PublicID)

// 	assert.Nil(err)
// 	assert.NotNil(result)
// 	assert.Equal("1", result.RegistrationStep)

// }

// func Test_ApplicantRepository_GetApplicantRegistrationStep_CorrectData_StepTwo(t *testing.T) {

// 	assert := assert.New(t)

// 	repository := applicants.NewApplicantRegistry().GetApplicantRepository()

// 	Applicant := testhelper.Helper_RandomApplicant(t)

// 	testhelper.Helper_ChangeRegistrationStep("2", Applicant, t)
// 	result, err := repository.GetApplicantRegistrationStep(Applicant.PublicID)

// 	assert.NotNil(result)
// 	assert.Equal("2", result.RegistrationStep)
// 	assert.Nil(err)
// }

// func Test_ApplicantRepository_GetApplicantRegistrationStep_CorrectData_StepThree(t *testing.T) {

// 	assert := assert.New(t)

// 	repository := applicants.NewApplicantRegistry().GetApplicantRepository()

// 	Applicant := testhelper.Helper_RandomApplicant(t)

// 	testhelper.Helper_ChangeRegistrationStep("3", Applicant, t)
// 	result, err := repository.GetApplicantRegistrationStep(Applicant.PublicID)

// 	assert.NotNil(result)
// 	assert.Equal("3", result.RegistrationStep)
// 	assert.Nil(err)
// }

// func Test_ApplicantRepository_GetApplicantRegistrationStep_IncorrectData_Step(t *testing.T) {

// 	assert := assert.New(t)

// 	Applicant := testhelper.Helper_RandomApplicant(t)

// 	err := testhelper.Helper_ChangeRegistrationStep("not-acceptable", Applicant, t)

// 	assert.NotNil(err)

// }

// func Test_ApplicantRepository_GetApplicantRegistrationStep_IncorrectData_PublicID(t *testing.T) {

// 	assert := assert.New(t)

// 	repository := applicants.NewApplicantRegistry().GetApplicantRepository()

// 	result, err := repository.GetApplicantRegistrationStep(string(encryption.GeneratePassword(9)))

// 	assert.Nil(result)
// 	assert.NotNil(err)
// }

// func Test_ApplicantRepository_SetApplicantRegistrationStep_IncorrectData_PublicID(t *testing.T) {
// 	assert := assert.New(t)

// 	repository := applicants.NewApplicantRegistry().GetApplicantRepository()

// 	result, err := repository.SetApplicantRegistrationStep("3", string(encryption.GeneratePassword(9)))

// 	assert.Nil(result)
// 	assert.NotNil(err)

// }

// func Test_ApplicantRepository_SetApplicantRegistrationStep_IncorrectData_Step(t *testing.T) {
// 	assert := assert.New(t)

// 	Applicant := testhelper.Helper_RandomApplicant(t)

// 	err := testhelper.Helper_ChangeRegistrationStep("not-acceptable", Applicant, t)

// 	assert.NotNil(err)
// }

// func Test_ApplicantRepository_SetApplicantRegistrationStep_CorrectData(t *testing.T) {
// 	assert := assert.New(t)

// 	repository := applicants.NewApplicantRegistry().GetApplicantRepository()

// 	Applicant := testhelper.Helper_RandomApplicant(t)

// 	testhelper.Helper_ChangeRegistrationStep("2", Applicant, t)
// 	result, err := repository.SetApplicantRegistrationStep("3", Applicant.PublicID)

// 	assert.NotNil(result)
// 	assert.Equal("3", result.RegistrationStep)
// 	assert.Nil(err)
// }
