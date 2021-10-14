package jobs

import "bit-jobs-api/shared/database"

type JobRegistry struct {
}

func NewJobRegistry() *JobRegistry {
	return &JobRegistry{}
}

func (*JobRegistry) GetJobRepository() *JobRepository {
	return NewJobRepository(database.DB)
}
