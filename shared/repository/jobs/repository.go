package jobs

import (
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

func (repository *JobRepository) EmployerCreateJob(employerPublicID, jobTitle, jobType, category, jobDescription, postStartDatetime, postEndDatetime, payPeriod string, minSalary, maxSalary int) (*Job, error) {

	if jobTitle == "" {
		return nil, errors.New("data cannot be empty")
	}

	var job Job
	var slug string

	slug = strings.ReplaceAll(jobTitle, " ", "-")
	slug = strings.ToLower(slug)

	stmt, err := repository.Database.Prepare(`
		INSERT INTO 
		jobs(title, jobtype, category, description, minsalary, 
				maxsalary, payperiod, poststartdatetime, postenddatetime, employerid, slug) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, (SELECT id FROM employers WHERE publicid=$10), $11) 
		RETURNING publicid;`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = stmt.QueryRow(jobTitle, jobType, category, jobDescription, minSalary, maxSalary, payPeriod, postStartDatetime, postEndDatetime, employerPublicID, slug).Scan(&job.PublicID)

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

func (repository *JobRepository) GetEmployerJobs(employerPublicID string) ([]*Job, int, error) {

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

func (repository *JobRepository) EditJob(employerPublicID, jobPublicID, jobTitle, jobType, category, jobDescription, postStartDatetime, postEndDatetime, payPeriod string, minSalary, maxSalary int) (*Job, error) {

	if employerPublicID == "" || jobPublicID == "" {
		return nil, errors.New("missing required value")
	}

	var slug string

	job, err := repository.GetJob(jobPublicID)

	if err != nil {
		return nil, err
	}

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
	stmt, err := repository.Database.Prepare(`UPDATE jobs SET title=$1, jobtype=$2, category=$3, description=$4, poststartdatetime=$5, postenddatetime=$6, minsalary=$7, maxsalary=$8, payperiod=$9, slug=$10 WHERE publicid=$11 AND employerid=(SELECT id FROM employers WHERE publicid=$12);`)

	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec(job.Title, job.JobType, job.Category, job.Description, job.PostStartDatetime, job.PostEndDatetime, job.MinSalary, job.MaxSalary, job.PayPeriod, slug, job.PublicID, employerPublicID)

	if err != nil {
		return nil, err
	}

	return job, nil
}
