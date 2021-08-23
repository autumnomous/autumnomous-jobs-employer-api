package employers

import (
	"autumnomous.com/bit-jobs-api/shared/database"
	"autumnomous.com/bit-jobs-api/shared/repository/employers/accountmanagement"
)

type EmployerRegistry struct {
}

func NewEmployerRegistry() *EmployerRegistry {
	return &EmployerRegistry{}
}

func (*EmployerRegistry) GetEmployerRepository() *accountmanagement.EmployerRepository {
	return accountmanagement.NewEmployerRepository(database.DB)
}
