package accountmanagement_test

import (
	"fmt"
	"log"
	"testing"

	"bit-jobs-api/shared/database"
	employers "bit-jobs-api/shared/repository/employers"
	"bit-jobs-api/shared/repository/employers/accountmanagement"
	"bit-jobs-api/shared/services/security/encryption"
	"bit-jobs-api/shared/testhelper"

	"github.com/stretchr/testify/assert"
)

func init() {
	testhelper.Init()
}

func Test_NewEmployerRepository(t *testing.T) {

	assert := assert.New(t)

	result := accountmanagement.NewEmployerRepository(database.DB)
	assert.Equal(database.DB, result.Database)
}

func Test_EmployerRepository_CreateEmployer(t *testing.T) {
	assert := assert.New(t)

	data := map[string]string{
		"firstname": "Kendrick",
		"lastname":  "Lamar",
		"email":     fmt.Sprintf("klamar-%s@damn.com", string(encryption.GeneratePassword(9))),
		"password":  string(encryption.GeneratePassword(9)),
	}

	repository := accountmanagement.NewEmployerRepository(database.DB)

	result, err := repository.CreateEmployer(data["firstname"], data["lastname"], data["email"], data["password"])

	assert.Equal(result.FirstName, data["firstname"])
	assert.Equal(result.LastName, data["lastname"])
	assert.Equal(result.Email, data["email"])
	assert.NotNil(result.PublicID)
	assert.Nil(err)

}

func Test_EmployerRepository_CreateEmployer_Fail_EmptyData(t *testing.T) {
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

		repository := accountmanagement.NewEmployerRepository(database.DB)

		result, err := repository.CreateEmployer(test["firstname"], test["lastname"], test["email"], test["password"])

		assert.NotNil(err)
		assert.Nil(result)

	}

}

func Test_EmployerRepository_GetEmployer(t *testing.T) {
	assert := assert.New(t)

	employer := testhelper.Helper_RandomEmployer(t)

	repository := accountmanagement.NewEmployerRepository(database.DB)

	result, err := repository.GetEmployer(employer.PublicID)

	assert.Equal(result.FirstName, employer.FirstName)
	assert.Equal(result.LastName, employer.LastName)
	assert.Equal(result.Email, employer.Email)
	assert.NotNil(result.PublicID)
	assert.Nil(err)

}

func Test_EmployerRepository_GetEmployer_Fail_EmptyData(t *testing.T) {
	assert := assert.New(t)

	repository := accountmanagement.NewEmployerRepository(database.DB)

	result, err := repository.GetEmployer("")

	assert.NotNil(err)
	assert.Nil(result)

}

func Test_EmployerRepository_AuthenticateEmployerPassword_NoDataReceived(t *testing.T) {

	assert := assert.New(t)

	repository := employers.NewEmployerRegistry().GetEmployerRepository()

	data := map[string]string{
		"email":    "",
		"password": "",
	}

	result, _, publicid, err := repository.AuthenticateEmployerPassword(data["email"], data["password"])

	assert.False(result)
	assert.Nil(err)
	assert.Equal("", publicid)
}

func Test_EmployerRepository_AuthenicateEmployerPassword_CorrectDataReceived(t *testing.T) {

	assert := assert.New(t)

	repository := employers.NewEmployerRegistry().GetEmployerRepository()

	data := map[string]string{
		"firstname": "First",
		"lastname":  "Last",
		"email":     fmt.Sprintf("email-%s@site.com", string(encryption.GeneratePassword(9))),
		"password":  string(encryption.GeneratePassword(9)),
	}

	Employer := &testhelper.TestEmployer{
		FirstName: data["firstname"],
		LastName:  data["lastname"],
		Email:     data["email"],
	}

	hashedPassword, err := encryption.HashPassword([]byte(data["password"]))

	Employer.HashedPassword = hashedPassword
	if err != nil {
		t.Fatal()
	}

	Employer = testhelper.Helper_CreateEmployer(Employer, t)
	result, _, publicid, err := repository.AuthenticateEmployerPassword(data["email"], data["password"])

	assert.True(result)
	assert.Nil(err)
	assert.Equal(Employer.PublicID, publicid)

}

func Test_EmployerRepository_AuthenticateEmployerPassword_IncorrectDataReceived_Email(t *testing.T) {
	assert := assert.New(t)

	repository := employers.NewEmployerRegistry().GetEmployerRepository()

	data := map[string]string{
		"firstname":      "First",
		"lastname":       "Last",
		"email":          fmt.Sprintf("email-%s@site.com", string(encryption.GeneratePassword(9))),
		"Employer-email": fmt.Sprintf("email-%s@site.com", string(encryption.GeneratePassword(9))),
		"password":       string(encryption.GeneratePassword(9)),
	}

	Employer := &testhelper.TestEmployer{
		FirstName: data["firstname"],
		LastName:  data["lastname"],
		Email:     data["Employer-email"],
		Password:  data["password"],
	}

	hashedPassword, err := encryption.HashPassword([]byte(Employer.Password))

	Employer.HashedPassword = hashedPassword
	if err != nil {
		t.Fatal()
	}

	Employer = testhelper.Helper_CreateEmployer(Employer, t)

	result, _, publicID, err := repository.AuthenticateEmployerPassword(data["email"], Employer.Password)

	assert.False(result)
	assert.Nil(err)
	assert.Equal("", publicID)

}

func Test_EmployerRepository_AuthenticateEmployerPassword_IncorrectDataReceived_Password(t *testing.T) {
	assert := assert.New(t)

	repository := employers.NewEmployerRegistry().GetEmployerRepository()

	data := map[string]string{
		"firstname":         "First",
		"lastname":          "Last",
		"email":             fmt.Sprintf("email-%s@site.com", string(encryption.GeneratePassword(9))),
		"employer-password": string(encryption.GeneratePassword(9)),
		"password":          string(encryption.GeneratePassword(9)),
	}

	employer := &testhelper.TestEmployer{
		FirstName: data["firstname"],
		LastName:  data["lastname"],
		Email:     data["email"],
	}

	hashedPassword, err := encryption.HashPassword([]byte(data["employer-password"]))

	employer.HashedPassword = hashedPassword
	if err != nil {
		t.Fatal()
	}

	_ = testhelper.Helper_CreateEmployer(employer, t)

	result, registrationStep, publicid, err := repository.AuthenticateEmployerPassword(data["email"], data["password"])

	assert.False(result)
	assert.Nil(err)
	assert.Equal("", registrationStep)
	assert.Equal("", publicid)

}

func Test_EmployerRepository_UpdateEmployerPassword_NoDataReceived(t *testing.T) {
	assert := assert.New(t)

	repository := employers.NewEmployerRegistry().GetEmployerRepository()

	data := map[string]string{
		"password":    "",
		"newpassword": "",
		"publicid":    "",
	}
	updated, err := repository.UpdateEmployerPassword(data["publicid"], data["password"], data["newpassword"])

	assert.False(updated)
	assert.Nil(err)

}

func Test_EmployerRepository_UpdateEmployerPassword_IncorrectPublicID(t *testing.T) {
	assert := assert.New(t)
	repository := employers.NewEmployerRegistry().GetEmployerRepository()

	data := map[string]string{
		"password":    string(encryption.GeneratePassword(9)),
		"newpassword": string(encryption.GeneratePassword(9)),
		"publicid":    string(encryption.GeneratePassword(10)),
	}

	updated, err := repository.UpdateEmployerPassword(data["publicid"], data["password"], data["newpassword"])

	assert.False(updated)
	assert.Nil(err)
}

func Test_EmployerRepository_UpdateEmployerPassword_CorrectData(t *testing.T) {

	assert := assert.New(t)

	repository := employers.NewEmployerRegistry().GetEmployerRepository()

	Employer := &testhelper.TestEmployer{
		FirstName: "First",
		LastName:  "Last",
		Email:     fmt.Sprintf("email-%s@site.com", encryption.GeneratePassword(9)),
	}

	password := encryption.GeneratePassword(9)
	hashedPassword, err := encryption.HashPassword([]byte(password))

	if err != nil {
		t.Fatal()
	}

	Employer.HashedPassword = hashedPassword

	Employer = testhelper.Helper_CreateEmployer(Employer, t)

	data := map[string]string{
		"password":    string(password),
		"newpassword": string(encryption.GeneratePassword(9)),
		"publicid":    Employer.PublicID,
	}

	updated, err := repository.UpdateEmployerPassword(data["publicid"], data["password"], data["newpassword"])

	assert.True(updated)
	assert.Nil(err)
}

func Test_EmployerRepository_UpdateEmployerAccount_CorrectData(t *testing.T) {

	assert := assert.New(t)

	repository := employers.NewEmployerRegistry().GetEmployerRepository()

	Employer := &testhelper.TestEmployer{
		FirstName: "First",
		LastName:  "Last",
		Email:     fmt.Sprintf("email-%s@site.com", encryption.GeneratePassword(9)),
	}

	password := encryption.GeneratePassword(9)
	hashedPassword, err := encryption.HashPassword([]byte(password))

	if err != nil {
		log.Println(err)
		t.Fatal()
	}

	Employer.HashedPassword = hashedPassword

	Employer = testhelper.Helper_CreateEmployer(Employer, t)

	data := map[string]string{
		"firstname": "NewFirst",
		"lastname":  "NewLast",
		"email":     fmt.Sprintf("new-email-%s@site.com", encryption.GeneratePassword(9)),
		"publicid":  Employer.PublicID,
	}

	updatedEmployer, err := repository.UpdateEmployerAccount(data["publicid"], data["firstname"], data["lastname"], data["email"], "", "", "", "", "", "")

	assert.NotNil(updatedEmployer)
	assert.Nil(err)
}

func Test_EmployerRepository_SetEmployerCompany_Correct(t *testing.T) {
	assert := assert.New(t)

	repository := employers.NewEmployerRegistry().GetEmployerRepository()

	employer := testhelper.Helper_RandomEmployer(t)
	company := testhelper.Helper_RandomCompany(t)

	err := repository.SetEmployerCompany(employer.PublicID, company.PublicID)

	result := testhelper.Helper_GetEmployer(employer.PublicID, t)

	assert.Equal(result.CompanyPublicID, company.PublicID)
	assert.Nil(err)

}

func Test_EmployerRepository_SetEmployerCompany_Incorrect_NoEmployerPublicID(t *testing.T) {
	assert := assert.New(t)

	repository := employers.NewEmployerRegistry().GetEmployerRepository()

	err := repository.SetEmployerCompany("", "notvalid")

	assert.NotNil(err)

}

func Test_EmployerRepository_SetEmployerCompany_Incorrect_NoCompanyPublicID(t *testing.T) {
	assert := assert.New(t)

	repository := employers.NewEmployerRegistry().GetEmployerRepository()

	err := repository.SetEmployerCompany("notvalid", "")

	assert.NotNil(err)

}

func Test_EmployerRepository_GetEmployerCompany_Incorrect_NoEmployerPublicID(t *testing.T) {

	assert := assert.New(t)

	repository := employers.NewEmployerRegistry().GetEmployerRepository()

	_, err := repository.GetEmployerCompany("")

	assert.NotNil(err)

}

func Test_EmployerRepository_GetEmployerCompany_Correct(t *testing.T) {

	assert := assert.New(t)

	repository := employers.NewEmployerRegistry().GetEmployerRepository()

	employer := testhelper.Helper_RandomEmployer(t)
	company := testhelper.Helper_RandomCompany(t)

	err := testhelper.Helper_SetEmployerCompany(employer.PublicID, company.PublicID)

	if err != nil {
		log.Println(err)
		t.Fatal()
	}

	result, err := repository.GetEmployerCompany(employer.PublicID)

	assert.Nil(err)
	assert.Equal(company.Domain, result.Domain)
	assert.Equal(company.Name, result.Name)

}

// func Test_EmployerRepository_GetEmployerRegistrationStep_CorrectData_NewEmployer(t *testing.T) {

// 	assert := assert.New(t)

// 	repository := Employers.NewEmployerRegistry().GetEmployerRepository()

// 	Employer := testhelper.Helper_RandomEmployer(t)

// 	result, err := repository.GetEmployerRegistrationStep(Employer.PublicID)

// 	assert.Nil(err)
// 	assert.NotNil(result)
// 	assert.Equal("1", result.RegistrationStep)

// }

// func Test_EmployerRepository_GetEmployerRegistrationStep_CorrectData_StepTwo(t *testing.T) {

// 	assert := assert.New(t)

// 	repository := Employers.NewEmployerRegistry().GetEmployerRepository()

// 	Employer := testhelper.Helper_RandomEmployer(t)

// 	testhelper.Helper_ChangeRegistrationStep("2", Employer, t)
// 	result, err := repository.GetEmployerRegistrationStep(Employer.PublicID)

// 	assert.NotNil(result)
// 	assert.Equal("2", result.RegistrationStep)
// 	assert.Nil(err)
// }

// func Test_EmployerRepository_GetEmployerRegistrationStep_CorrectData_StepThree(t *testing.T) {

// 	assert := assert.New(t)

// 	repository := Employers.NewEmployerRegistry().GetEmployerRepository()

// 	Employer := testhelper.Helper_RandomEmployer(t)

// 	testhelper.Helper_ChangeRegistrationStep("3", Employer, t)
// 	result, err := repository.GetEmployerRegistrationStep(Employer.PublicID)

// 	assert.NotNil(result)
// 	assert.Equal("3", result.RegistrationStep)
// 	assert.Nil(err)
// }

// func Test_EmployerRepository_GetEmployerRegistrationStep_IncorrectData_Step(t *testing.T) {

// 	assert := assert.New(t)

// 	Employer := testhelper.Helper_RandomEmployer(t)

// 	err := testhelper.Helper_ChangeRegistrationStep("not-acceptable", Employer, t)

// 	assert.NotNil(err)

// }

// func Test_EmployerRepository_GetEmployerRegistrationStep_IncorrectData_PublicID(t *testing.T) {

// 	assert := assert.New(t)

// 	repository := Employers.NewEmployerRegistry().GetEmployerRepository()

// 	result, err := repository.GetEmployerRegistrationStep(string(encryption.GeneratePassword(9)))

// 	assert.Nil(result)
// 	assert.NotNil(err)
// }

// func Test_EmployerRepository_SetEmployerRegistrationStep_IncorrectData_PublicID(t *testing.T) {
// 	assert := assert.New(t)

// 	repository := Employers.NewEmployerRegistry().GetEmployerRepository()

// 	result, err := repository.SetEmployerRegistrationStep("3", string(encryption.GeneratePassword(9)))

// 	assert.Nil(result)
// 	assert.NotNil(err)

// }

// func Test_EmployerRepository_SetEmployerRegistrationStep_IncorrectData_Step(t *testing.T) {
// 	assert := assert.New(t)

// 	Employer := testhelper.Helper_RandomEmployer(t)

// 	err := testhelper.Helper_ChangeRegistrationStep("not-acceptable", Employer, t)

// 	assert.NotNil(err)
// }

// func Test_EmployerRepository_SetEmployerRegistrationStep_CorrectData(t *testing.T) {
// 	assert := assert.New(t)

// 	repository := Employers.NewEmployerRegistry().GetEmployerRepository()

// 	Employer := testhelper.Helper_RandomEmployer(t)

// 	testhelper.Helper_ChangeRegistrationStep("2", Employer, t)
// 	result, err := repository.SetEmployerRegistrationStep("3", Employer.PublicID)

// 	assert.NotNil(result)
// 	assert.Equal("3", result.RegistrationStep)
// 	assert.Nil(err)
// }
