package companies

import (
	"database/sql"
	"log"
)

type CompanyRepository struct {
	Database *sql.DB
}

type Company struct {
	Name         string  `json:"name"`
	Domain       string  `json:"domain"`
	Location     string  `json:"location"`
	Longitude    float64 `json:"longitude"`
	Latitude     float64 `json:"latitude"`
	URL          string  `json:"url"`
	Facebook     string  `json:"facebook"`
	Twitter      string  `json:"twitter"`
	Instagram    string  `json:"instagram"`
	Description  string  `json:"description"`
	Logo         string  `json:"logo"`
	ExtraDetails string  `json:"extradetails"`
	PublicID     string  `json:"publicid"`
}

func NewCompanyRepository(db *sql.DB) *CompanyRepository {
	return &CompanyRepository{Database: db}
}

func (repository *CompanyRepository) GetOrCreateCompany(domain, name, location, url, facebook, twitter, instagram, description, logo, extradetails string) (*Company, error) {
	var company Company

	stmt, err := repository.Database.Prepare(`SELECT name, location, url, facebook, twitter, instagram, description, logo, extradetails, publicid FROM companies WHERE domain=$1;`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = stmt.QueryRow(domain).Scan(&company.Name, &company.Location, &company.URL, &company.Facebook, &company.Twitter, &company.Instagram, &company.Description, &company.Logo, &company.ExtraDetails, &company.PublicID)

	if err != nil {

		if err.Error() == "sql: no rows in result set" {
			stmt, err := repository.Database.Prepare(`
				INSERT INTO companies(name, location, url, facebook, twitter, instagram, description, logo, extradetails, domain) VALUES
				($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
				RETURNING publicid;`)

			if err != nil {
				log.Println(err)
				return nil, err
			}

			err = stmt.QueryRow(name, location, url, facebook, twitter, instagram, description, logo, extradetails, domain).Scan(&company.PublicID)

			if err != nil {
				log.Println(err)
				return nil, err
			}

			company.Name = name
			company.Location = location
			company.URL = url
			company.Facebook = facebook
			company.Twitter = twitter
			company.Instagram = instagram
			company.Description = description
			company.Logo = logo
			company.ExtraDetails = extradetails
			company.Domain = domain

		} else {
			log.Println(err)
			return nil, err
		}

	}

	return &company, nil
}
