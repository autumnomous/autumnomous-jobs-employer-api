package companies

import "bit-jobs-api/shared/database"

type CompanyRegistry struct {
}

func NewCompanyRegistry() *CompanyRegistry {
	return &CompanyRegistry{}
}

func (*CompanyRegistry) GetCompanyRepository() *CompanyRepository {
	return NewCompanyRepository(database.DB)
}
