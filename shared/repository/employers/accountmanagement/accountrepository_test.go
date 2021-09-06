package accountmanagement_test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"autumnomous.com/bit-jobs-api/shared/database"
	employers "autumnomous.com/bit-jobs-api/shared/repository/employers"
	"autumnomous.com/bit-jobs-api/shared/repository/employers/accountmanagement"
	"autumnomous.com/bit-jobs-api/shared/services/security/encryption"
	"autumnomous.com/bit-jobs-api/shared/testhelper"

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

	result, initialpasswordchanged, publicid, err := repository.AuthenticateEmployerPassword(data["email"], data["password"])

	assert.False(result)
	assert.Nil(err)
	assert.False(initialpasswordchanged)
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

	updatedEmployer, err := repository.UpdateEmployerAccount(data["publicid"], data["firstname"], data["lastname"], data["email"])

	assert.NotNil(updatedEmployer)
	assert.Nil(err)
}

func Test_EmployerRepository_EmployerCreateJob_IncorrectData(t *testing.T) {
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
		"title":             "",
		"streetaddress":     "123 Street Avenue",
		"city":              "City",
		"zipcode":           "00000",
		"tags":              "full-time,remote-friendly",
		"description":       "This is a new job",
		"payperiod":         "year",
		"poststartdatetime": time.Now().String(),
		"postenddatetime":   time.Now().String(),
	}

	minSalary := 10000
	maxSalary := 100000

	job, err := repository.EmployerCreateJob(Employer.PublicID, data["title"], data["streetaddress"], data["city"], data["zipcode"], data["tags"], data["description"], data["poststartdatetime"], data["postenddatetime"], data["payperiod"], minSalary, maxSalary)

	assert.NotNil(err)
	assert.Nil(job)
}

func Test_EmployerRepository_EmployerCreateJob_CorrectData(t *testing.T) {
	assert := assert.New(t)

	repository := employers.NewEmployerRegistry().GetEmployerRepository()

	Employer := testhelper.Helper_RandomEmployer(t)

	data := map[string]string{
		"title":             "Job Title",
		"streetaddress":     "123 Street Avenue",
		"city":              "City",
		"zipcode":           "00000",
		"tags":              "full-time,remote-friendly",
		"description":       "This is a new job",
		"payperiod":         "year",
		"poststartdatetime": time.Now().Format(time.RFC3339),
		"postenddatetime":   time.Now().Format(time.RFC3339),
	}

	minSalary := 10000
	maxSalary := 100000
	job, err := repository.EmployerCreateJob(Employer.PublicID, data["title"], data["streetaddress"], data["city"], data["zipcode"], data["tags"], data["description"], data["poststartdatetime"], data["postenddatetime"], data["payperiod"], minSalary, maxSalary)

	assert.NotNil(job)
	assert.Nil(err)
}

func Test_EmployerRepository_GetJob_IncorrectData(t *testing.T) {

	assert := assert.New(t)

	repository := employers.NewEmployerRegistry().GetEmployerRepository()

	job, err := repository.GetJob("")

	assert.Nil(job)
	assert.NotNil(err)
}

func Test_EmployerRepository_GetJob_CorrectData(t *testing.T) {

	assert := assert.New(t)

	repository := employers.NewEmployerRegistry().GetEmployerRepository()

	employer := testhelper.Helper_RandomEmployer(t)

	testjob := testhelper.Helper_RandomJob(employer, t)

	job, err := repository.GetJob(testjob.PublicID)

	assert.Equal(testjob.City, job.City)
	assert.Equal(testjob.Title, job.Title)
	assert.Equal(testjob.StreetAddress, job.StreetAddress)
	assert.Equal(testjob.EmployerPublicID, job.EmployerPublicID)
	assert.Equal(testjob.Tags, job.Tags)
	assert.Nil(err)

}

func Test_EmployerRepository_GetEmployerJobs_IncorrectData(t *testing.T) {
	assert := assert.New(t)

	repository := employers.NewEmployerRegistry().GetEmployerRepository()

	jobs, err := repository.GetEmployerJobs("")

	assert.Nil(jobs)
	assert.NotNil(err)
}

func Test_EmployerRepository_GetEmployerJobs_CorrectData(t *testing.T) {
	assert := assert.New(t)

	repository := employers.NewEmployerRegistry().GetEmployerRepository()

	employer := testhelper.Helper_RandomEmployer(t)

	testhelper.Helper_RandomJob(employer, t)
	testhelper.Helper_RandomJob(employer, t)
	testhelper.Helper_RandomJob(employer, t)

	jobs, err := repository.GetEmployerJobs(employer.PublicID)

	assert.Equal(len(jobs), 3)
	assert.Nil(err)
}

func Test_EmployerRepository_DeleteJob_IncorrectData_MissingJobPublicID(t *testing.T) {
	assert := assert.New(t)

	repository := employers.NewEmployerRegistry().GetEmployerRepository()

	employer := testhelper.Helper_RandomEmployer(t)

	// job := testhelper.Helper_RandomJob(employer, t)

	job, err := repository.DeleteJob(employer.PublicID, "")

	assert.Nil(job)
	assert.NotNil(err)
}

func Test_EmployerRepository_DeleteJob_IncorrectData_MissingEmployerPublicID(t *testing.T) {
	assert := assert.New(t)

	repository := employers.NewEmployerRegistry().GetEmployerRepository()

	employer := testhelper.Helper_RandomEmployer(t)

	job := testhelper.Helper_RandomJob(employer, t)

	result, err := repository.DeleteJob("", job.PublicID)

	assert.Nil(result)
	assert.NotNil(err)
}

func Test_EmployerRepository_DeleteJob_Correct(t *testing.T) {
	assert := assert.New(t)

	repository := employers.NewEmployerRegistry().GetEmployerRepository()

	employer := testhelper.Helper_RandomEmployer(t)

	job := testhelper.Helper_RandomJob(employer, t)

	result, err := repository.DeleteJob(employer.PublicID, job.PublicID)

	assert.NotNil(result)
	assert.Nil(err)
}

func Test_EmployerRepository_EditJob_MissingJobPublicID(t *testing.T) {
	assert := assert.New(t)

	repository := employers.NewEmployerRegistry().GetEmployerRepository()

	employer := testhelper.Helper_RandomEmployer(t)

	job, err := repository.EditJob(employer.PublicID, "", "", "", "", "", "", "", "", "", "", 0, 0)

	assert.Nil(job)
	assert.NotNil(err)

}

func Test_EmployerRepository_EditJob_MissingEmployerPublicID(t *testing.T) {
	assert := assert.New(t)

	repository := employers.NewEmployerRegistry().GetEmployerRepository()

	job, err := repository.EditJob("", "", "", "", "", "", "", "", "", "", "", 0, 0)

	assert.Nil(job)
	assert.NotNil(err)

}

func Test_EmployerRepository_EditJob_Correct(t *testing.T) {
	assert := assert.New(t)

	repository := employers.NewEmployerRegistry().GetEmployerRepository()

	employer := testhelper.Helper_RandomEmployer(t)

	job := testhelper.Helper_RandomJob(employer, t)

	result, err := repository.EditJob(employer.PublicID, job.PublicID, "A Job", "123 Street", "City", "00000", "full-time,remote", "this is a job", "2021-09-04", "2021-10-01", "hourly", 40, 50)

	assert.NotNil(result)
	assert.Equal(result.Title, "A Job")
	assert.Equal(result.StreetAddress, "123 Street")
	assert.Nil(err)
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
