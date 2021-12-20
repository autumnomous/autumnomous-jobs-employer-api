package jobs

import (
	"bit-jobs-api/shared/services/utils"
	"database/sql"
	"errors"
	"log"
	"strings"
)

type JobRepository struct {
	Database *sql.DB
}

func NewJobRepository(db *sql.DB) *JobRepository {
	return &JobRepository{Database: db}
}

type Job struct {
	PublicID         string `json:"publicid"`
	Title            string `json:"title"`
	JobType          string `json:"jobtype"`
	Category         string `json:"category"`
	Description      string `json:"description"` // make required?
	EmployerPublicID string `json:"employerpublicid"`
	Remote           bool   `json:"remote"`
	VisibleDate      string `json:"visibledate"`
	MinSalary        int64  `json:"minsalary"`
	MaxSalary        int64  `json:"maxsalary"`
	PayPeriod        string `json:"payperiod"`
}

func (repository *JobRepository) EmployerCreateJob(employerPublicID, jobTitle, jobType, category, jobDescription, visibleDate, payPeriod string, remote bool, minSalary, maxSalary int64) (*Job, error) {

	if jobTitle == "" {
		return nil, errors.New("data cannot be empty")
	}

	var job Job
	var slug string

	slug = strings.ReplaceAll(jobTitle, " ", "-")
	slug = strings.ToLower(slug)

	stmt, err := repository.Database.Prepare(`
		INSERT INTO 
		jobs(title, jobtype, category, description,visibledate, remote, employerid, slug, minsalary, maxsalary, payperiod) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9 (SELECT id FROM employers WHERE publicid=$10), $11) 
		RETURNING publicid;`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = stmt.QueryRow(jobTitle, jobType, category, jobDescription, utils.NewNullString(visibleDate), remote, employerPublicID, slug, minSalary, maxSalary, payPeriod).Scan(&job.PublicID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return repository.GetJob(job.PublicID)

}

func (repository *JobRepository) GetJob(jobPublicID string) (*Job, error) {

	if jobPublicID == "" {
		return nil, errors.New("missing required value")
	}
	var job Job
	var visibleDate sql.NullString
	stmt, err := repository.Database.Prepare(`
		SELECT jobs.title, jobs.jobtype, jobs.category, jobs.description, jobs.visibledate, jobs.remote, jobs.minsalary, jobs.maxsalary, jobs.payperiod, employers.publicid
		FROM jobs
		JOIN employers ON employers.id=jobs.employerid
		WHERE jobs.publicid=$1;`,
	)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = stmt.QueryRow(jobPublicID).Scan(&job.Title, &job.JobType, &job.Category, &job.Description, &visibleDate, &job.Remote, &job.MinSalary, &job.MaxSalary, &job.PayPeriod, &job.EmployerPublicID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	job.PublicID = jobPublicID

	if visibleDate.Valid {
		job.VisibleDate = visibleDate.String
	}

	return &job, nil
}

func (repository *JobRepository) GetEmployerJobs(employerPublicID string) ([]*Job, error) {

	if employerPublicID == "" {
		return nil, errors.New("missing required value")
	}

	var jobs []*Job

	stmt, err := repository.Database.Prepare(`
			SELECT jobs.title, jobs.jobtype, jobs.category, jobs.description, 
				jobs.visibledate, jobs.remote, jobs.minsalary, jobs.maxsalary, jobs.payperiod, jobs.publicid
			FROM jobs
			JOIN employers ON employers.id=jobs.employerid
			WHERE jobs.employerid=(SELECT id FROM employers WHERE publicid=$1);`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	rows, err := stmt.Query(employerPublicID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		job := &Job{}
		var visibleDate sql.NullString

		err := rows.Scan(&job.Title, &job.JobType, &job.Category, &job.Description, &visibleDate, &job.Remote, &job.MinSalary, &job.MaxSalary, &job.PayPeriod, &job.PublicID)

		if err != nil {
			log.Println(err)
			return nil, err
		}
		job.EmployerPublicID = employerPublicID

		if visibleDate.Valid {
			job.VisibleDate = visibleDate.String
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

func (repository *JobRepository) DeleteJob(employerPublicID, jobPublicID string) (*Job, error) {

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

func (repository *JobRepository) EditJob(employerPublicID, jobPublicID, jobTitle, jobType, category, jobDescription, visibleDate, payPeriod string, remote bool, minSalary, maxSalary int64) (*Job, error) {

	if employerPublicID == "" || jobPublicID == "" {
		return nil, errors.New("missing required value")
	}

	var slug string
	// var postEndDatetime time.Time

	job, err := repository.GetJob(jobPublicID)

	if err != nil {
		return nil, err
	}

	job.Remote = remote

	if jobTitle != "" {
		job.Title = jobTitle
		slug = strings.ReplaceAll(jobTitle, " ", "-")
		slug = strings.ToLower(slug)
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

	if visibleDate != "" {
		job.VisibleDate = visibleDate
	}

	if minSalary != 0 {
		job.MinSalary = minSalary
	}

	if maxSalary != 0 {
		job.MaxSalary = maxSalary
	}

	stmt, err := repository.Database.Prepare(`UPDATE jobs SET title=$1, jobtype=$2, category=$3, description=$4, visibledate=$5, slug=$6, remote=$7 , minsalary=$8, maxsalary=$9, payperiod=$10 WHERE publicid=$11 AND employerid=(SELECT id FROM employers WHERE publicid=$12);`)

	if err != nil {
		return nil, err
	}
	_, err = stmt.Exec(job.Title, job.JobType, job.Category, job.Description, job.VisibleDate, slug, job.Remote, job.MinSalary, job.MaxSalary, job.PayPeriod, job.PublicID, employerPublicID)

	if err != nil {
		return nil, err
	}

	return job, nil
}
