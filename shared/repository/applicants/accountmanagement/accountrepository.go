package accountmanagement

import (
	"database/sql"
	"errors"
	"log"

	"autumnomous.com/bit-jobs-api/shared/services/security/encryption"
	_ "github.com/lib/pq"
)

type ApplicantRepository struct {
	Database *sql.DB
}

type Applicant struct {
	FirstName              string `json:"firstname"`
	LastName               string `json:"lastname"`
	Email                  string `json:"email"`
	Password               string
	Biography              string `json:"bio"`
	PublicID               string `json:"publicid"`
	InitialPasswordChanged bool   `json:"initialpasswordchanged"`
}

func NewApplicantRepository(db *sql.DB) *ApplicantRepository {
	return &ApplicantRepository{Database: db}
}

func (repository *ApplicantRepository) CreateApplicant(firstName, lastName, email, password string) (*Applicant, error) {

	if firstName == "" || lastName == "" || email == "" || password == "" {
		return nil, errors.New("data cannot be empty")
	}

	applicant := &Applicant{FirstName: firstName, LastName: lastName, Email: email, Password: password}

	stmt, err := repository.Database.Prepare(`INSERT INTO applicants(email, firstname, lastname, password) VALUES ($1, $2, $3, $4) RETURNING publicid;`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = stmt.QueryRow(email, firstName, lastName, password).Scan(&applicant.PublicID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return applicant, nil
}

func (repository *ApplicantRepository) AuthenticateApplicantPassword(email, password string) (bool, bool, string, error) {

	if email == "" || password == "" {
		return false, false, "", nil
	}

	var databasePassword, publicID string
	var initialPasswordChanged bool
	stmt, err := repository.Database.Prepare(`SELECT password, initialpasswordchanged, publicid FROM applicants WHERE email=$1;`)

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

func (repository *ApplicantRepository) UpdateApplicantPassword(publicID, password, newPassword string) (bool, error) {

	if publicID == "" || password == "" || newPassword == "" {
		return false, nil
	}
	var databasePassword string

	stmt, err := repository.Database.Prepare(`SELECT password FROM applicants WHERE publicid=$1;`)

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
		stmt, err = repository.Database.Prepare(`UPDATE applicants SET password=$1, initialpasswordchanged=true WHERE publicid=$2;`)

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

func (repository *ApplicantRepository) UpdateApplicantAccount(publicID, firstName, lastName, email, biography string) (*Applicant, error) {

	applicant := &Applicant{}

	stmt, err := repository.Database.Prepare(`SELECT firstname, lastname, email, biography FROM users WHERE publicid=$1;`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = stmt.QueryRow(publicID).Scan(&applicant.FirstName, &applicant.LastName, &applicant.Email, &applicant.Biography)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if firstName != "" {
		applicant.FirstName = firstName
	}

	if lastName != "" {
		applicant.LastName = lastName
	}

	if email != "" {
		applicant.Email = email
	}

	if biography != "" {
		applicant.Biography = biography
	}

	applicant.PublicID = publicID
	stmt, err = repository.Database.Prepare(`UPDATE applicants SET firstname=$1, lastname=$2, email=$3, biography=$4 WHERE publicid=$5;`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	_, err = stmt.Exec(applicant.FirstName, applicant.LastName, applicant.Email, applicant.Biography, applicant.PublicID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return applicant, nil
}
