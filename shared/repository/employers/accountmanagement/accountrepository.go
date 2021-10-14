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
	Facebook         string `json:"facebook"`
	Twitter          string `json:"twitter"`
	Instagram        string `json:"instagram"`
	TotalPostsBought int    `json:"totalpostsbought"`
	RegistrationStep string `json:"registrationstep"`
	Password         string
	CompanyPublicID  string `json:"companypublicid"`
	PublicID         string `json:"publicid"`
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
		SELECT firstname, lastname, email, totalpostsbought, registrationstep, mobilenumber, phonenumber, role, facebook, twitter, instagram
		FROM employers
		WHERE publicid=$1;`,
	)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	var emp_mobile_number, emp_work_number, emp_role, emp_facebook, emp_twitter, emp_instagram sql.NullString

	err = stmt.QueryRow(userID).Scan(&employer.FirstName, &employer.LastName, &employer.Email, &employer.TotalPostsBought, &employer.RegistrationStep, &emp_mobile_number, &emp_work_number, &emp_role, &emp_facebook, &emp_twitter, &emp_instagram)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if emp_mobile_number.Valid {
		employer.MobileNumber = emp_mobile_number.String
	}

	if emp_work_number.Valid {
		employer.PhoneNumber = emp_work_number.String
	}

	if emp_role.Valid {
		employer.Role = emp_role.String
	}

	if emp_facebook.Valid {
		employer.Facebook = emp_facebook.String
	}

	if emp_twitter.Valid {
		employer.Twitter = emp_twitter.String
	}

	if emp_instagram.Valid {
		employer.Instagram = emp_instagram.String
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

func (repository *EmployerRepository) UpdateEmployerAccount(publicID, firstName, lastName, email, phoneNumber, mobileNumber, role, facebook, twitter, instagram string) (*Employer, error) {

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

	Employer.Facebook = facebook
	Employer.Twitter = twitter
	Employer.Instagram = instagram
	Employer.PublicID = publicID
	stmt, err = repository.Database.Prepare(`UPDATE employers SET firstname=$1, lastname=$2, email=$3, phonenumber=$4, mobilenumber=$5, role=$6, facebook=$7, twitter=$8, instagram=$9 WHERE publicid=$10;`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	_, err = stmt.Exec(Employer.FirstName, Employer.LastName, Employer.Email, Employer.PhoneNumber, Employer.MobileNumber, Employer.Role, facebook, twitter, instagram, Employer.PublicID)

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

func (repository *EmployerRepository) GetJobPackage(typeID string) (*JobPackage, error) {

	var pack JobPackage
	stmt, err := repository.Database.Prepare(`
			SELECT id, typeid, isactive, title, numberofjobs, description, price
			FROM jobpackages
			WHERE typeid=$1;`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = stmt.QueryRow(typeID).Scan(&pack.ID, &pack.TypeID, &pack.IsActive, &pack.Title, &pack.NumberOfJobs, &pack.Description, &pack.Price)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pack, nil

}
