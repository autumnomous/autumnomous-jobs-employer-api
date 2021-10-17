package jobpackages

import (
	"database/sql"
	"log"
)

type JobPackageRepository struct {
	Database *sql.DB
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

func NewJobPackageRepository(db *sql.DB) *JobPackageRepository {
	return &JobPackageRepository{Database: db}
}

func (repository *JobPackageRepository) GetActiveJobPackages() ([]*JobPackage, error) {

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

func (repository *JobPackageRepository) GetJobPackage(typeID string) (*JobPackage, error) {

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
