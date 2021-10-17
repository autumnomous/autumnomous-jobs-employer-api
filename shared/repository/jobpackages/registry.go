package jobpackages

import (
	"bit-jobs-api/shared/database"
)

type JobPackageRegistry struct {
}

func NewJobPackageRegistry() *JobPackageRegistry {
	return &JobPackageRegistry{}
}

func (registry *JobPackageRegistry) GetJobPackageRepository() *JobPackageRepository {
	return NewJobPackageRepository(database.DB)
}
