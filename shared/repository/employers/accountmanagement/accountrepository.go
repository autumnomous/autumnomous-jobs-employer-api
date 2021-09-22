package accountmanagement

import (
	"database/sql"
	"errors"
	"log"

	"bit-jobs-api/shared/services/security/encryption"

	_ "github.com/lib/pq"
)

type EmployerRepository struct {
	Database *sql.DB
}

type Employer struct {
	FirstName              string `json:"firstname"`
	LastName               string `json:"lastname"`
	Email                  string `json:"email"`
	Password               string
	CompanyPublicID        string `json:"companypublicid"`
	PublicID               string `json:"publicid"`
	InitialPasswordChanged bool   `json:"initialpasswordchanged"`
}

type Job struct {
	PublicID          string `json:"publicid"`
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

type JobPackage struct {
	ID           int     `json:"id"`
	TypeID       string  `json:"typeid"`
	IsActive     bool    `json:"isactive"`
	Title        string  `json:"title"`
	NumberOfJobs int     `json:"numberofjobs"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
}

func NewEmployerRepository(db *sql.DB) *EmployerRepository {
	return &EmployerRepository{Database: db}
}

func (repository *EmployerRepository) CreateEmployer(firstName, lastName, email, password string) (*Employer, error) {

	if firstName == "" || lastName == "" || email == "" || password == "" {
		return nil, errors.New("data cannot be empty")
	}

	employer := &Employer{FirstName: firstName, LastName: lastName, Email: email}

	stmt, err := repository.Database.Prepare(`INSERT INTO employers(email, firstname, lastname, password) VALUES ($1, $2, $3, $4) RETURNING publicid;`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = stmt.QueryRow(email, firstName, lastName, password).Scan(&employer.PublicID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return employer, nil
}

func (repository *EmployerRepository) GetEmployer(userID string) (*Employer, error) {
	// TODO

	return &Employer{}, nil
}

func (repository *EmployerRepository) AuthenticateEmployerPassword(email, password string) (bool, bool, string, error) {

	if email == "" || password == "" {
		return false, false, "", nil
	}

	var databasePassword, publicID string
	var initialPasswordChanged bool
	stmt, err := repository.Database.Prepare(`SELECT password, initialpasswordchanged, publicid FROM employers WHERE email=$1;`)

	if err != nil {
		log.Println(err)
		return false, false, "", err
	}

	err = stmt.QueryRow(email).Scan(&databasePassword, &initialPasswordChanged, &publicID)

	if err != nil {

		if err.Error() == "sql: no rows in result set" {
			return false, false, "", nil
		} else {
			log.Println(err)
			return false, false, "", err
		}

	}

	if encryption.CompareHashes([]byte(databasePassword), []byte(password)) {
		return true, initialPasswordChanged, publicID, nil
	}

	return false, false, "", nil
}

func (repository *EmployerRepository) UpdateEmployerPassword(publicID, password, newPassword string) (bool, error) {

	if publicID == "" || password == "" || newPassword == "" {
		return false, nil
	}
	var databasePassword string

	stmt, err := repository.Database.Prepare(`SELECT password FROM employers WHERE publicid=$1;`)

	if err != nil {
		log.Println(err)
		return false, err
	}

	err = stmt.QueryRow(publicID).Scan(&databasePassword)

	if err != nil {

		if err.Error() == "sql: no rows in result set" {
			return false, nil
		} else {
			log.Println(err)
			return false, err
		}

	}

	if encryption.CompareHashes([]byte(databasePassword), []byte(password)) {
		stmt, err = repository.Database.Prepare(`UPDATE employers SET password=$1, initialpasswordchanged=true WHERE publicid=$2;`)

		if err != nil {
			log.Println(err)
			return false, err
		}

		hashedNewPassword, err := encryption.HashPassword([]byte(newPassword))

		if err != nil {
			log.Println(err)
			return false, err
		}

		_, err = stmt.Exec(hashedNewPassword, publicID)

		if err != nil {
			log.Println(err)
			return false, err
		}

		return true, nil
	} else {
		return false, nil
	}
}

func (repository *EmployerRepository) UpdateEmployerAccount(publicID, firstName, lastName, email string) (*Employer, error) {

	Employer := &Employer{}

	stmt, err := repository.Database.Prepare(`SELECT firstname, lastname, email FROM employers WHERE publicid=$1;`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = stmt.QueryRow(publicID).Scan(&Employer.FirstName, &Employer.LastName, &Employer.Email)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if firstName != "" {
		Employer.FirstName = firstName
	}

	if lastName != "" {
		Employer.LastName = lastName
	}

	if email != "" {
		Employer.Email = email
	}

	Employer.PublicID = publicID
	stmt, err = repository.Database.Prepare(`UPDATE employers SET firstname=$1, lastname=$2, email=$3 WHERE publicid=$4;`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	_, err = stmt.Exec(Employer.FirstName, Employer.LastName, Employer.Email, Employer.PublicID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return Employer, nil
}

func (repository *EmployerRepository) EmployerCreateJob(employerPublicID, jobTitle, jobStreetAddress, jobCity, jobZipCode, jobTags, jobDescription, postStartDatetime, postEndDatetime, payPeriod string, minSalary, maxSalary int) (*Job, error) {

	if jobTitle == "" {
		return nil, errors.New("data cannot be empty")
	}

	var job Job

	stmt, err := repository.Database.Prepare(`
		INSERT INTO 
		jobs(title, streetaddress, city, zipcode, tags, description, minsalary, 
				maxsalary, payperiod, poststartdatetime, postenddatetime, employerid) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, (SELECT id FROM employers WHERE publicid=$12)) 
		RETURNING publicid;`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = stmt.QueryRow(jobTitle, jobStreetAddress, jobCity, jobZipCode, jobTags, jobDescription, minSalary, maxSalary, payPeriod, postStartDatetime, postEndDatetime, employerPublicID).Scan(&job.PublicID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return repository.GetJob(job.PublicID)

}

func (repository *EmployerRepository) GetJob(jobPublicID string) (*Job, error) {

	if jobPublicID == "" {
		return nil, errors.New("missing required value")
	}
	var job Job

	stmt, err := repository.Database.Prepare(`
		SELECT jobs.title, jobs.city, jobs.streetaddress, jobs.zipcode, jobs.tags, jobs.description, 
			jobs.minsalary, jobs.maxsalary, jobs.payperiod, jobs.poststartdatetime,
			jobs.postenddatetime, employers.publicid
		FROM jobs
		JOIN employers ON employers.id=jobs.employerid
		WHERE jobs.publicid=$1;`,
	)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = stmt.QueryRow(jobPublicID).Scan(&job.Title, &job.City, &job.StreetAddress, &job.ZipCode, &job.Tags, &job.Description, &job.MinSalary, &job.MaxSalary, &job.PayPeriod, &job.PostStartDatetime, &job.PostEndDatetime, &job.EmployerPublicID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	job.PublicID = jobPublicID

	return &job, nil
}

func (repository *EmployerRepository) GetEmployerJobs(employerPublicID string) ([]*Job, int, error) {

	if employerPublicID == "" {
		return nil, -1, errors.New("missing required value")
	}

	var jobs []*Job
	var employerTotalPostsBought int

	stmt, err := repository.Database.Prepare(`
			SELECT jobs.title, jobs.tags, jobs.description, 
				jobs.minsalary, jobs.maxsalary, jobs.payperiod, jobs.poststartdatetime, jobs.postenddatetime, jobs.publicid,
				employers.totalpostsbought
			FROM jobs
			JOIN employers ON employers.id=jobs.employerid
			WHERE jobs.employerid=(SELECT id FROM employers WHERE publicid=$1);`)

	if err != nil {
		log.Println(err)
		return nil, -1, err
	}

	rows, err := stmt.Query(employerPublicID)

	if err != nil {
		log.Println(err)
		return nil, -1, err
	}

	defer rows.Close()
	for rows.Next() {
		job := &Job{}

		err := rows.Scan(&job.Title, &job.Tags, &job.Description, &job.MinSalary, &job.MaxSalary, &job.PayPeriod, &job.PostStartDatetime, &job.PostEndDatetime, &job.PublicID, &employerTotalPostsBought)

		if err != nil {
			log.Println(err)
			return nil, -1, err
		}
		job.EmployerPublicID = employerPublicID

		jobs = append(jobs, job)
	}

	return jobs, employerTotalPostsBought, nil
}

func (repository *EmployerRepository) DeleteJob(employerPublicID, jobPublicID string) (*Job, error) {

	if employerPublicID == "" || jobPublicID == "" {
		return nil, errors.New("missing required value")
	}

	var job Job

	stmt, err := repository.Database.Prepare(`DELETE FROM jobs WHERE publicid=$1 AND employerid=(SELECT id FROM employers WHERE publicid=$2) RETURNING title, city, streetaddress, zipcode,tags, description;`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = stmt.QueryRow(jobPublicID, employerPublicID).Scan(&job.Title, &job.City, &job.StreetAddress, &job.ZipCode, &job.Tags, &job.Description)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &job, nil
}

func (repository *EmployerRepository) EditJob(employerPublicID, jobPublicID, jobTitle, jobStreetAddress, jobCity, jobZipCode, jobTags, jobDescription, postStartDatetime, postEndDatetime, payPeriod string, minSalary, maxSalary int) (*Job, error) {

	if employerPublicID == "" || jobPublicID == "" {
		return nil, errors.New("missing required value")
	}

	job, err := repository.GetJob(jobPublicID)

	if err != nil {
		return nil, err
	}

	if jobTitle != "" {
		job.Title = jobTitle
	}

	if jobStreetAddress != "" {
		job.StreetAddress = jobStreetAddress
	}

	if jobCity != "" {
		job.City = jobCity
	}

	if jobZipCode != "" {
		job.ZipCode = jobZipCode
	}

	if jobTags != "" {
		job.Tags = jobTags
	}

	if jobDescription != "" {
		job.Description = jobDescription
	}

	if postStartDatetime != "" {
		job.PostStartDatetime = postStartDatetime
	}

	if postEndDatetime != "" {
		job.PostEndDatetime = postEndDatetime
	}

	if minSalary != 0 {
		job.MinSalary = minSalary
	}

	if maxSalary != 0 {
		job.MaxSalary = maxSalary
	}

	if payPeriod != "" {
		job.PayPeriod = payPeriod
	}
	stmt, err := repository.Database.Prepare(`UPDATE jobs SET title=$1, streetaddress=$2, city=$3, zipcode=$4, tags=$5, description=$6 , poststartdatetime=$7, postenddatetime=$8, minsalary=$9, maxsalary=$10, payperiod=$11 WHERE publicid=$12 AND employerid=(SELECT id FROM employers WHERE publicid=$13);`)

	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec(job.Title, job.StreetAddress, job.City, job.ZipCode, job.Tags, job.Description, job.PostStartDatetime, job.PostEndDatetime, job.MinSalary, job.MaxSalary, job.PayPeriod, job.PublicID, employerPublicID)

	if err != nil {
		return nil, err
	}

	return job, nil
}

func (repository *EmployerRepository) GetActiveJobPackages() ([]*JobPackage, error) {

	var packages []*JobPackage

	stmt, err := repository.Database.Prepare(`
			SELECT id, typeid, isactive, title, numberofjobs, description, price
			FROM jobpackages
			WHERE isactive=TRUE;`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	rows, err := stmt.Query()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		jobPackage := &JobPackage{}

		err := rows.Scan(&jobPackage.ID, &jobPackage.TypeID, &jobPackage.IsActive, &jobPackage.Title, &jobPackage.NumberOfJobs, &jobPackage.Description, &jobPackage.Price)

		if err != nil {
			log.Println(err)
			return nil, err
		}

		packages = append(packages, jobPackage)
	}

	return packages, nil
}
