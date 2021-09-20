package accountmanagement

import (
	"database/sql"
	"log"
)

type UserRepository struct {
	Database *sql.DB
}

type User struct {
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	AccountType string `json:"accounttype"`
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{Database: db}
}

func (repository *UserRepository) GetUser(publicID string) (*User, error) {

	var user User

	stmt, err := repository.Database.Prepare(`SELECT firstname, lastname, accounttype FROM users WHERE publicid=$1;`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = stmt.QueryRow(publicID).Scan(&user.FirstName, &user.LastName, &user.AccountType)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &user, nil
}
