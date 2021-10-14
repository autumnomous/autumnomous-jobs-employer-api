package jobs_test

import (
	"bit-jobs-api/shared/repository/jobs"
	"bit-jobs-api/shared/services/security/encryption"
	"bit-jobs-api/shared/testhelper"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_EmployerRepository_EmployerCreateJob_IncorrectData(t *testing.T) {
	assert := assert.New(t)

	repository := jobs.NewJobRegistry().GetJobRepository()

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
		"jobtype":           "full-time",
		"category":          "full-stack",
		"description":       "This is a new job",
		"payperiod":         "year",
		"poststartdatetime": time.Now().String(),
		"postenddatetime":   time.Now().String(),
	}

	minSalary := 10000
	maxSalary := 100000

	job, err := repository.EmployerCreateJob(Employer.PublicID, data["title"], data["jobtype"], data["category"], data["description"], data["poststartdatetime"], data["postenddatetime"], data["payperiod"], minSalary, maxSalary)

	assert.NotNil(err)
	assert.Nil(job)
}

func Test_EmployerRepository_EmployerCreateJob_CorrectData(t *testing.T) {
	assert := assert.New(t)

	repository := jobs.NewJobRegistry().GetJobRepository()

	Employer := testhelper.Helper_RandomEmployer(t)

	data := map[string]string{
		"title":             "Job Title",
		"jobtype":           "full-time",
		"category":          "full-stack",
		"description":       "This is a new job",
		"payperiod":         "year",
		"poststartdatetime": time.Now().Format(time.RFC3339),
		"postenddatetime":   time.Now().Format(time.RFC3339),
	}

	minSalary := 10000
	maxSalary := 100000
	job, err := repository.EmployerCreateJob(Employer.PublicID, data["title"], data["jobtype"], data["category"], data["description"], data["poststartdatetime"], data["postenddatetime"], data["payperiod"], minSalary, maxSalary)

	assert.NotNil(job)
	assert.Nil(err)
}

func Test_EmployerRepository_GetJob_IncorrectData(t *testing.T) {

	assert := assert.New(t)

	repository := jobs.NewJobRegistry().GetJobRepository()

	job, err := repository.GetJob("")

	assert.Nil(job)
	assert.NotNil(err)
}

func Test_EmployerRepository_GetJob_CorrectData(t *testing.T) {

	assert := assert.New(t)

	repository := jobs.NewJobRegistry().GetJobRepository()

	employer := testhelper.Helper_RandomEmployer(t)

	testjob := testhelper.Helper_RandomJob(employer, t)

	job, err := repository.GetJob(testjob.PublicID)

	assert.Equal(testjob.Title, job.Title)
	assert.Equal(testjob.EmployerPublicID, job.EmployerPublicID)
	assert.Nil(err)

}

func Test_EmployerRepository_GetEmployerJobs_IncorrectData(t *testing.T) {
	assert := assert.New(t)

	repository := jobs.NewJobRegistry().GetJobRepository()
	var totalpostsbought int
	jobs, totalpostsbought, err := repository.GetEmployerJobs("")

	assert.Nil(jobs)
	assert.Equal(totalpostsbought, -1)
	assert.NotNil(err)
}

func Test_EmployerRepository_GetEmployerJobs_CorrectData(t *testing.T) {
	assert := assert.New(t)

	repository := jobs.NewJobRegistry().GetJobRepository()

	employer := testhelper.Helper_RandomEmployer(t)
	var totalpostsbought int
	testhelper.Helper_RandomJob(employer, t)
	testhelper.Helper_RandomJob(employer, t)
	testhelper.Helper_RandomJob(employer, t)

	jobs, totalpostsbought, err := repository.GetEmployerJobs(employer.PublicID)

	employer.TotalPostsBought = totalpostsbought
	assert.Equal(len(jobs), 3)
	assert.Equal(employer.TotalPostsBought, 0)
	assert.Nil(err)
}

func Test_EmployerRepository_DeleteJob_IncorrectData_MissingJobPublicID(t *testing.T) {
	assert := assert.New(t)

	repository := jobs.NewJobRegistry().GetJobRepository()

	employer := testhelper.Helper_RandomEmployer(t)

	// job := testhelper.Helper_RandomJob(employer, t)

	job, err := repository.DeleteJob(employer.PublicID, "")

	assert.Nil(job)
	assert.NotNil(err)
}

func Test_EmployerRepository_DeleteJob_IncorrectData_MissingEmployerPublicID(t *testing.T) {
	assert := assert.New(t)

	repository := jobs.NewJobRegistry().GetJobRepository()

	employer := testhelper.Helper_RandomEmployer(t)

	job := testhelper.Helper_RandomJob(employer, t)

	result, err := repository.DeleteJob("", job.PublicID)

	assert.Nil(result)
	assert.NotNil(err)
}

func Test_EmployerRepository_DeleteJob_Correct(t *testing.T) {
	assert := assert.New(t)

	repository := jobs.NewJobRegistry().GetJobRepository()

	employer := testhelper.Helper_RandomEmployer(t)

	job := testhelper.Helper_RandomJob(employer, t)

	result, err := repository.DeleteJob(employer.PublicID, job.PublicID)

	assert.NotNil(result)
	assert.Nil(err)
}

func Test_EmployerRepository_EditJob_MissingJobPublicID(t *testing.T) {
	assert := assert.New(t)

	repository := jobs.NewJobRegistry().GetJobRepository()

	employer := testhelper.Helper_RandomEmployer(t)

	job, err := repository.EditJob(employer.PublicID, "", "", "", "", "", "", "", "", 0, 0)

	assert.Nil(job)
	assert.NotNil(err)

}

func Test_EmployerRepository_EditJob_MissingEmployerPublicID(t *testing.T) {
	assert := assert.New(t)

	repository := jobs.NewJobRegistry().GetJobRepository()

	job, err := repository.EditJob("", "", "", "", "", "", "", "", "", 0, 0)

	assert.Nil(job)
	assert.NotNil(err)

}

func Test_EmployerRepository_EditJob_Correct(t *testing.T) {
	assert := assert.New(t)

	repository := jobs.NewJobRegistry().GetJobRepository()

	employer := testhelper.Helper_RandomEmployer(t)

	job := testhelper.Helper_RandomJob(employer, t)

	result, err := repository.EditJob(employer.PublicID, job.PublicID, "A Job", "full-time", "full-stack", "this is a job", "2021-09-04", "2021-10-01", "hourly", 40, 50)

	assert.NotNil(result)
	assert.Equal(result.Title, "A Job")
	assert.Nil(err)
}
