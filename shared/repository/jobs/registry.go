package jobs

import "autumnomous-jobs-employer-api/shared/database"

type JobRegistry struct {
}

func NewJobRegistry() *JobRegistry {
	return &JobRegistry{}
}

func (*JobRegistry) GetJobRepository() *JobRepository {
	return NewJobRepository(database.DB)
}
