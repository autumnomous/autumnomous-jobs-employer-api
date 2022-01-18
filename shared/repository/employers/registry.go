package employers

import (
	"autumnomous-jobs-employer-api/shared/database"
	"autumnomous-jobs-employer-api/shared/repository/employers/accountmanagement"
)

type EmployerRegistry struct {
}

func NewEmployerRegistry() *EmployerRegistry {
	return &EmployerRegistry{}
}

func (*EmployerRegistry) GetEmployerRepository() *accountmanagement.EmployerRepository {
	return accountmanagement.NewEmployerRepository(database.DB)
}
