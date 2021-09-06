package testhelper

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"autumnomous.com/bit-jobs-api/shared/database"
	"autumnomous.com/bit-jobs-api/shared/services/security/encryption"

	"github.com/joho/godotenv"
)

// Applicant
type TestUser struct {
	FirstName      string
	LastName       string
	Email          string
	Password       string
	Biography      string
	HashedPassword []byte
	PublicID       string
}

type TestEmployer struct {
	FirstName       string
	LastName        string
	Email           string
	Password        string
	CompanyPublicID string
	HashedPassword  []byte
	PublicID        string
}

type TestJob struct {
	PublicID          string `jaon:"publicid"`
	Title             string `json:"title"`
	City              string `json:"city"`
	StreetAddress     string `json:"streetaddress"`
	ZipCode           string `json:"zipcode"`
	Tags              string `json:"tags"`
	Description       string `json:"description"` // make required?
	EmployerPublicID  string `json:"employerpublicid"`
	MinSalary         int    `json:"minsalary"`
	MaxSalary         int    `json:"maxsalary"`
	PayPeriod         string `json:"payperiod"`
	PostStartDatetime string `json:"poststartdatetime"`
	PostEndDatetime   string `json:"postenddatetime"`
}

type TestApplication struct {
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	Email       string `json:"email"`
	ID          int
	ApplicantID string `json:"applicantid"`
	PublicID    string `json:"publicid"`
}

func Init() {
	os.Clearenv()

	err := godotenv.Load("test.env")

	if err != nil {
		log.Println(err)
		log.Fatal("Error loading test.env file:", err)
	}

	database.Connect("HEROKU_POSTGRESQL_JADE_URL")

}

func Helper_CreateApplicant(applicant *TestUser, t *testing.T) *TestUser {

	stmt, err := database.DB.Prepare(`INSERT INTO users (firstname, lastname, email, password, accounttype) VALUES ($1, $2, $3, $4, 'applicant') RETURNING publicid;`)

	if err != nil {
		log.Println(err)
		t.Fatal()
	}

	err = stmt.QueryRow(applicant.FirstName, applicant.LastName, applicant.Email, applicant.HashedPassword).Scan(&applicant.PublicID)

	if err != nil {
		log.Println(err)
		t.Fatal()
	}

	return applicant

}

func Helper_GetUser(publicID string, t *testing.T) *TestUser {

	stmt, err := database.DB.Prepare(`SELECT firstname, lastname, email, password, biography, publicid FROM users WHERE publicid=$1;`)

	if err != nil {
		log.Println(err)
		t.Fatal()
	}

	user := &TestUser{}
	err = stmt.QueryRow(publicID).Scan(&user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Biography, &user.PublicID)

	if err != nil {
		log.Println(err)
		t.Fatal()
	}

	return user
}

func Helper_CreateEmployer(employer *TestEmployer, t *testing.T) *TestEmployer {

	stmt, err := database.DB.Prepare(`INSERT INTO employers (firstname, lastname, email, password) VALUES ($1, $2, $3, $4) RETURNING publicid;`)

	if err != nil {
		log.Println(err)
		t.Fatal()
	}

	err = stmt.QueryRow(employer.FirstName, employer.LastName, employer.Email, employer.HashedPassword).Scan(&employer.PublicID)

	if err != nil {
		log.Println(err)
		t.Fatal()
	}

	return employer

}
func Helper_RandomEmployer(t *testing.T) *TestEmployer {
	employer := &TestEmployer{FirstName: string(encryption.GeneratePassword(5)),
		LastName: string(encryption.GeneratePassword(5)),
		Email:    fmt.Sprintf("email-%s@site.com", string(encryption.GeneratePassword(9))),
		Password: string(encryption.GeneratePassword(9))}

	hashedPassword, err := encryption.HashPassword([]byte(employer.Password))

	if err != nil {
		t.Fatal()
	}

	employer.HashedPassword = hashedPassword

	return Helper_CreateEmployer(employer, t)

}

func Helper_RandomApplicant(t *testing.T) *TestUser {
	applicant := &TestUser{FirstName: string(encryption.GeneratePassword(5)),
		LastName: string(encryption.GeneratePassword(5)),
		Email:    fmt.Sprintf("email-%s@site.com", string(encryption.GeneratePassword(9))),
		Password: string(encryption.GeneratePassword(9))}

	hashedPassword, err := encryption.HashPassword([]byte(applicant.Password))

	if err != nil {
		t.Fatal()
	}

	applicant.HashedPassword = hashedPassword

	return Helper_CreateApplicant(applicant, t)

}

func Helper_CreateJob(job *TestJob, t *testing.T) *TestJob {
	stmt, err := database.DB.Prepare(`INSERT INTO 
											jobs(title, streetaddress, city, zipcode, tags, description,minsalary, maxsalary, payperiod, poststartdatetime, postenddatetime, employerid) 
											VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, (SELECT id FROM employers WHERE publicid=$12)) 
											RETURNING publicid;`)
	if err != nil {
		log.Println(err)
		return nil
	}

	err = stmt.QueryRow(job.Title, job.StreetAddress, job.City, job.ZipCode, job.Tags, job.Description, job.MinSalary, job.MaxSalary, job.PayPeriod, job.PostStartDatetime, job.PostEndDatetime, job.EmployerPublicID).Scan(&job.PublicID)

	if err != nil {
		log.Println(err)
		return nil
	}

	return job
}

func Helper_RandomJob(employer *TestEmployer, t *testing.T) *TestJob {

	job := &TestJob{
		Title:             string(encryption.GeneratePassword(5)),
		City:              string(encryption.GeneratePassword(5)),
		StreetAddress:     string(encryption.GeneratePassword(5)),
		ZipCode:           string(encryption.GeneratePassword(5)),
		Tags:              strings.Join([]string{string(encryption.GeneratePassword(5)), string(encryption.GeneratePassword(5)), string(encryption.GeneratePassword(5))}, ","),
		Description:       string(encryption.GeneratePassword(5)),
		MinSalary:         0,
		MaxSalary:         100,
		PayPeriod:         "hour",
		PostStartDatetime: time.Now().Format(time.RFC3339),
		PostEndDatetime:   time.Now().Format(time.RFC3339),
		EmployerPublicID:  employer.PublicID,
	}

	return Helper_CreateJob(job, t)
}

// func Helper_ChangeRegistrationStep(step string, Applicant *TestUser, t *testing.T) error {

// 	stmt, err := database.DB.Prepare(`UPDATE applications SET registrationstep=$1 WHERE Applicantid=(SELECT id FROM Applicants WHERE publicid=$2);`)

// 	if err != nil {
// 		t.Fatal()
// 	}

// 	err = stmt.QueryRow(step, Applicant.PublicID).Err()

// 	return err

// }
