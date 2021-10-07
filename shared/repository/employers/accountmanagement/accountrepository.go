package accountmanagement

import (
	"database/sql"
	"errors"
	"log"

	"bit-jobs-api/shared/repository/companies"
	"bit-jobs-api/shared/services/security/encryption"

	_ "github.com/lib/pq"
)

type EmployerRepository struct {
	Database *sql.DB
}

type Employer struct {
	FirstName        string `json:"firstname"`
	LastName         string `json:"lastname"`
	Email            string `json:"email"`
	PhoneNumber      string `json:"phonenumber"`
	MobileNumber     string `json:"mobilenumber"`
	Role             string `json:"role"`
	TotalPostsBought int    `json:"totalpostsbought"`
	RegistrationStep string `json:"registrationstep"`
	Password         string
	CompanyPublicID  string `json:"companypublicid"`
	PublicID         string `json:"publicid"`
}

type Job struct {
	PublicID          string `json:"publicid"`
	Title             string `json:"title"`
	JobType           string `json:"jobtype"`
	Category          string `json:"category"`
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

// RegistrationStep represents which stage in the registration process the user is in
type RegistrationStep int64

const (
	// ChangePassword Registration Step 1
	ChangePassword RegistrationStep = iota

	// PersonalInformation Registration Step 2
	PersonalInformation

	// CompanyDetails Registration Step 3
	CompanyDetails

	// PaymentMethod Registration Step 4
	PaymentMethod

	// PaymentDetails Registration Step 5
	PaymentDetails

	// Complete Registration Step 6
	RegistrationComplete
)

func (rs RegistrationStep) String() string {
	return [...]string{"change-password", "personal-information", "company-details", "payment-method", "payment-details", "registration-complete"}[rs]
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

	if userID == "" {
		return nil, errors.New("missing required value")
	}
	var employer Employer

	stmt, err := repository.Database.Prepare(`
		SELECT firstname, lastname, email, totalpostsbought, registrationstep, mobilenumber, phonenumber, role
		FROM employers
		WHERE publicid=$1;`,
	)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = stmt.QueryRow(userID).Scan(&employer.FirstName, &employer.LastName, &employer.Email, &employer.TotalPostsBought, &employer.RegistrationStep, &employer.MobileNumber, &employer.PhoneNumber, &employer.Role)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	employer.PublicID = userID

	return &employer, nil
}

func (repository *EmployerRepository) AuthenticateEmployerPassword(email, password string) (bool, string, string, error) {

	if email == "" || password == "" {
		return false, "", "", nil
	}

	var databasePassword, publicID, registrationStep string
	stmt, err := repository.Database.Prepare(`SELECT password, registrationStep, publicid FROM employers WHERE email=$1;`)

	if err != nil {
		log.Println(err)
		return false, "", "", err
	}

	err = stmt.QueryRow(email).Scan(&databasePassword, &registrationStep, &publicID)

	if err != nil {

		if err.Error() == "sql: no rows in result set" {
			return false, "", "", nil
		} else {
			log.Println(err)
			return false, "", "", err
		}

	}

	if encryption.CompareHashes([]byte(databasePassword), []byte(password)) {
		return true, registrationStep, publicID, nil
	}

	return false, "", "", nil
}

func (repository *EmployerRepository) UpdateEmployerPassword(publicID, password, newPassword string) (bool, error) {

	if publicID == "" || password == "" || newPassword == "" {
		return false, nil
	}
	var databasePassword, registrationStep string

	stmt, err := repository.Database.Prepare(`SELECT password, registrationstep FROM employers WHERE publicid=$1;`)

	if err != nil {
		log.Println(err)
		return false, err
	}

	err = stmt.QueryRow(publicID).Scan(&databasePassword, &registrationStep)

	if err != nil {

		if err.Error() == "sql: no rows in result set" {
			return false, nil
		} else {
			log.Println(err)
			return false, err
		}

	}

	if encryption.CompareHashes([]byte(databasePassword), []byte(password)) {

		if registrationStep == ChangePassword.String() {
			stmt, err = repository.Database.Prepare(`UPDATE employers SET registrationstep='personal-information' WHERE publicid=$1;`)

			if err != nil {
				log.Println(err)
				return false, err
			}

			_, err = stmt.Exec(publicID)

			if err != nil {
				log.Println(err)
				return false, err
			}

		}
		stmt, err = repository.Database.Prepare(`UPDATE employers SET password=$1 WHERE publicid=$2;`)

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

func (repository *EmployerRepository) UpdateEmployerAccount(publicID, firstName, lastName, email, phoneNumber, mobileNumber, role string) (*Employer, error) {

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

	if phoneNumber != "" {
		Employer.PhoneNumber = phoneNumber
	}

	if mobileNumber != "" {
		Employer.MobileNumber = mobileNumber
	}

	if role != "" {
		Employer.Role = role
	}

	Employer.PublicID = publicID
	stmt, err = repository.Database.Prepare(`UPDATE employers SET firstname=$1, lastname=$2, email=$3, phonenumber=$4, mobilenumber=$5, role=$6 WHERE publicid=$7;`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	_, err = stmt.Exec(Employer.FirstName, Employer.LastName, Employer.Email, Employer.PhoneNumber, Employer.MobileNumber, Employer.Role, Employer.PublicID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	emp, _ := repository.GetEmployer(publicID)

	if emp.RegistrationStep == PersonalInformation.String() {
		stmt, _ = repository.Database.Prepare(`UPDATE employers SET registrationstep='company-details' WHERE publicid=$1;`)

		stmt.Exec(publicID)

	}

	return Employer, nil
}

func (repository *EmployerRepository) UpdateEmployerCompany(employerPublicID, companyName, location, url, facebook, twitter, instagram, description, logo, extradetails string) (*companies.Company, error) {

	var company companies.Company

	stmt, err := repository.Database.Prepare(`
		SELECT companies.name, companies.location, companies.url, 
			companies.facebook, companies.twitter, companies.instagram, 
			companies.logo, companies.description, companies.extradetails,
			companies.domain, companies.publicid
		FROM companies 
		JOIN employers ON employers.companyid = companies.id 
		WHERE employers.publicid=$1;`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = stmt.QueryRow(employerPublicID).Scan(&company.Name, &company.Location, &company.URL, &company.Facebook, &company.Twitter, &company.Instagram, &company.Description, &company.Logo, &company.ExtraDetails, &company.Domain, &company.PublicID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if companyName != "" {
		company.Name = companyName
	}

	if location != "" {
		company.Location = location
	}

	if url != "" {
		company.URL = url
	}

	if facebook != "" {
		company.Facebook = facebook
	}

	if twitter != "" {
		company.Twitter = twitter
	}

	if instagram != "" {
		company.Instagram = instagram
	}

	if description != "" {
		company.Description = description
	}

	if logo != "" {
		company.Logo = logo
	}

	if extradetails != "" {
		company.ExtraDetails = extradetails
	}

	stmt, err = repository.Database.Prepare(`UPDATE companies SET name=$1, location=$2, url=$3, facebook=$4, twitter=$5, instagram=$6, description=$7, logo=$8, extradetails=$9 WHERE publicid=$10;`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	_, err = stmt.Exec(company.Name, company.Location, company.URL, company.Facebook, company.Twitter, company.Instagram, company.Description, company.Logo, company.ExtraDetails, company.PublicID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	emp, _ := repository.GetEmployer(employerPublicID)

	if emp.RegistrationStep == CompanyDetails.String() {
		stmt, _ = repository.Database.Prepare(`UPDATE employers SET registrationstep='payment-method' WHERE publicid=$1;`)

		stmt.Exec(employerPublicID)

	}

	return &company, nil

}

func (repository *EmployerRepository) UpdateEmployerPaymentMethod(employerPublicID, paymentMethod string) error {

	emp, _ := repository.GetEmployer(employerPublicID)

	if emp.RegistrationStep == PaymentMethod.String() {
		stmt, err := repository.Database.Prepare(`UPDATE employers SET registrationstep='payment-details' WHERE publicid=$1;`)

		if err != nil {
			log.Println(err)
			return err
		}

		stmt.Exec(employerPublicID)

	}

	return nil
}

func (repository *EmployerRepository) UpdateEmployerPaymentDetails(employerPublicID, paymentDetails string) error {

	emp, _ := repository.GetEmployer(employerPublicID)

	if emp.RegistrationStep == PaymentDetails.String() {
		stmt, err := repository.Database.Prepare(`UPDATE employers SET registrationstep='registration-complete' WHERE publicid=$1;`)

		if err != nil {
			log.Println(err)
			return err
		}

		stmt.Exec(employerPublicID)

	}

	return nil

}

func (repository *EmployerRepository) EmployerCreateJob(employerPublicID, jobTitle, jobType, category, jobDescription, postStartDatetime, postEndDatetime, payPeriod string, minSalary, maxSalary int) (*Job, error) {

	if jobTitle == "" {
		return nil, errors.New("data cannot be empty")
	}

	var job Job

	stmt, err := repository.Database.Prepare(`
		INSERT INTO 
		jobs(title, jobtype, category, description, minsalary, 
				maxsalary, payperiod, poststartdatetime, postenddatetime, employerid) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, (SELECT id FROM employers WHERE publicid=$10)) 
		RETURNING publicid;`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = stmt.QueryRow(jobTitle, jobType, category, jobDescription, minSalary, maxSalary, payPeriod, postStartDatetime, postEndDatetime, employerPublicID).Scan(&job.PublicID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return repository.GetJob(job.PublicID)

}

func (repository *EmployerRepository) SetEmployerCompany(employerPublicID, companyPublicID string) error {

	if employerPublicID == "" || companyPublicID == "" {
		return errors.New("missing required value")
	}

	stmt, err := repository.Database.Prepare(`UPDATE employers SET companyid=(SELECT id FROM companies WHERE publicid=$1) WHERE publicid=$2;`)

	if err != nil {
		return err
	}

	_, err = stmt.Exec(companyPublicID, employerPublicID)

	if err != nil {
		return err
	}

	return nil
}

func (repository *EmployerRepository) GetEmployerCompany(employerPublicID string) (*companies.Company, error) {

	if employerPublicID == "" {
		return nil, errors.New("missing required value")
	}
	var company companies.Company
	stmt, err := repository.Database.Prepare(`
				SELECT 
					name, domain, location, url, facebook, twitter, instagram,
					description, logo, extradetails, publicid
				FROM companies 
				WHERE id = (SELECT companyid FROM employers WHERE publicid=$1);`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = stmt.QueryRow(employerPublicID).Scan(&company.Name, &company.Domain, &company.Location, &company.URL, &company.Facebook, &company.Twitter, &company.Instagram, &company.Description, &company.Logo, &company.ExtraDetails, &company.PublicID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &company, nil
}

func (repository *EmployerRepository) GetJob(jobPublicID string) (*Job, error) {

	if jobPublicID == "" {
		return nil, errors.New("missing required value")
	}
	var job Job

	stmt, err := repository.Database.Prepare(`
		SELECT jobs.title, jobs.jobtype, jobs.category, jobs.description, 
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

	err = stmt.QueryRow(jobPublicID).Scan(&job.Title, &job.JobType, &job.Category, &job.Description, &job.MinSalary, &job.MaxSalary, &job.PayPeriod, &job.PostStartDatetime, &job.PostEndDatetime, &job.EmployerPublicID)

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
			SELECT jobs.title, jobs.jobtype, jobs.category, jobs.description, 
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

		err := rows.Scan(&job.Title, &job.JobType, &job.Category, &job.Description, &job.MinSalary, &job.MaxSalary, &job.PayPeriod, &job.PostStartDatetime, &job.PostEndDatetime, &job.PublicID, &employerTotalPostsBought)

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

	stmt, err := repository.Database.Prepare(`DELETE FROM jobs WHERE publicid=$1 AND employerid=(SELECT id FROM employers WHERE publicid=$2) RETURNING title, jobtype, category, description;`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = stmt.QueryRow(jobPublicID, employerPublicID).Scan(&job.Title, &job.JobType, &job.Category, &job.Description)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &job, nil
}

func (repository *EmployerRepository) EditJob(employerPublicID, jobPublicID, jobTitle, jobType, category, jobDescription, postStartDatetime, postEndDatetime, payPeriod string, minSalary, maxSalary int) (*Job, error) {

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

	if jobType != "" {
		job.JobType = jobType
	}

	if category != "" {
		job.Category = category
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
	stmt, err := repository.Database.Prepare(`UPDATE jobs SET title=$1, jobtype=$2, category=$3, description=$4, poststartdatetime=$5, postenddatetime=$6, minsalary=$7, maxsalary=$8, payperiod=$9 WHERE publicid=$10 AND employerid=(SELECT id FROM employers WHERE publicid=$11);`)

	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec(job.Title, job.JobType, job.Category, job.Description, job.PostStartDatetime, job.PostEndDatetime, job.MinSalary, job.MaxSalary, job.PayPeriod, job.PublicID, employerPublicID)

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
