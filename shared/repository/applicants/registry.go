package applicants

import (
	"autumnomous.com/bit-jobs-api/shared/database"
	"autumnomous.com/bit-jobs-api/shared/repository/applicants/accountmanagement"
)

type ApplicantRegistry struct {
}

func NewApplicantRegistry() *ApplicantRegistry {
	return &ApplicantRegistry{}
}

func (*ApplicantRegistry) GetApplicantRepository() *accountmanagement.ApplicantRepository {
	return accountmanagement.NewApplicantRepository(database.DB)
}
