package employers

import (
	"bit-jobs-api/shared/database"
	"bit-jobs-api/shared/repository/employers/accountmanagement"
)

type EmployerRegistry struct {
}

func NewEmployerRegistry() *EmployerRegistry {
	return &EmployerRegistry{}
}

func (*EmployerRegistry) GetEmployerRepository() *accountmanagement.EmployerRepository {
	return accountmanagement.NewEmployerRepository(database.DB)
}
