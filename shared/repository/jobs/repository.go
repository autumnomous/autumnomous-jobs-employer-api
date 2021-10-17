package jobs

import (
	"bit-jobs-api/shared/services/utils"
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"
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
	Remote            bool   `json:"remote"`
	PostStartDatetime string `json:"poststartdatetime"`
	PostEndDatetime   string `json:"postenddatetime"`
}

func (repository *JobRepository) EmployerCreateJob(employerPublicID, jobTitle, jobType, category, jobDescription, postStartDatetime string, remote bool) (*Job, error) {

	if jobTitle == "" {
		return nil, errors.New("data cannot be empty")
	}

	var job Job
	var slug string

	slug = strings.ReplaceAll(jobTitle, " ", "-")
	slug = strings.ToLower(slug)

	stmt, err := repository.Database.Prepare(`
		INSERT INTO 
		jobs(title, jobtype, category, description,poststartdatetime, remote, employerid, slug) 
		VALUES ($1, $2, $3, $4, $5, $6, (SELECT id FROM employers WHERE publicid=$7), $8) 
		RETURNING publicid;`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = stmt.QueryRow(jobTitle, jobType, category, jobDescription, utils.NewNullString(postStartDatetime), remote, employerPublicID, slug).Scan(&job.PublicID)

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
	var postStartDatetime, postEndDatetime sql.NullString
	stmt, err := repository.Database.Prepare(`
		SELECT jobs.title, jobs.jobtype, jobs.category, jobs.description, jobs.poststartdatetime,
			jobs.postenddatetime, jobs.remote, employers.publicid
		FROM jobs
		JOIN employers ON employers.id=jobs.employerid
		WHERE jobs.publicid=$1;`,
	)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = stmt.QueryRow(jobPublicID).Scan(&job.Title, &job.JobType, &job.Category, &job.Description, &postStartDatetime, &postEndDatetime, &job.Remote, &job.EmployerPublicID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	job.PublicID = jobPublicID

	if postStartDatetime.Valid {
		job.PostStartDatetime = postStartDatetime.String
	}

	if postEndDatetime.Valid {
		job.PostEndDatetime = postEndDatetime.String
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
				jobs.poststartdatetime, jobs.postenddatetime, jobs.remote, jobs.publicid
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
		var postStartDatetime, postEndDatetime sql.NullString

		err := rows.Scan(&job.Title, &job.JobType, &job.Category, &job.Description, &postStartDatetime, &postEndDatetime, &job.Remote, &job.PublicID)

		if err != nil {
			log.Println(err)
			return nil, err
		}
		job.EmployerPublicID = employerPublicID

		if postStartDatetime.Valid {
			job.PostStartDatetime = postStartDatetime.String
		}

		if postEndDatetime.Valid {
			job.PostEndDatetime = postEndDatetime.String
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

func (repository *JobRepository) EditJob(employerPublicID, jobPublicID, jobTitle, jobType, category, jobDescription, postStartDatetime string, remote bool) (*Job, error) {

	if employerPublicID == "" || jobPublicID == "" {
		return nil, errors.New("missing required value")
	}

	var slug string
	var postEndDatetime time.Time

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

		startDate, err := time.Parse("2006-01-02", postStartDatetime)
		if err != nil {
			return nil, err
		}

		postEndDatetime = startDate.Add(time.Hour * 24 * 30)
		job.PostEndDatetime = postEndDatetime.Local().Format("2006-01-02")
	}

	stmt, err := repository.Database.Prepare(`UPDATE jobs SET title=$1, jobtype=$2, category=$3, description=$4, poststartdatetime=$5, postenddatetime=$6, slug=$7, remote=$8 WHERE publicid=$9 AND employerid=(SELECT id FROM employers WHERE publicid=$10);`)

	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec(job.Title, job.JobType, job.Category, job.Description, job.PostStartDatetime, job.PostEndDatetime, slug, job.Remote, job.PublicID, employerPublicID)

	if err != nil {
		return nil, err
	}

	return job, nil
}
