package companies_test

import (
	"bit-jobs-api/shared/repository/companies"
	"bit-jobs-api/shared/services/security/encryption"
	"bit-jobs-api/shared/testhelper"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	testhelper.Init()
}

func Test_CompanyRepository_GetOrCreateCompany_Create_Correct(t *testing.T) {
	assert := assert.New(t)

	data := map[string]string{
		"domain":       fmt.Sprintf("example-%s.com", string(encryption.GeneratePassword(5))),
		"name":         string(encryption.GeneratePassword(5)),
		"location":     "",
		"url":          "",
		"facebook":     "",
		"twitter":      "",
		"instagram":    "",
		"description":  "",
		"logo":         "",
		"extradetails": "",
	}

	repository := companies.NewCompanyRegistry().GetCompanyRepository()

	result, err := repository.GetOrCreateCompany(data["domain"], data["name"], data["location"], data["url"], data["facebook"], data["twitter"], data["instagram"], data["description"], data["logo"], data["extradetails"])

	assert.Nil(err)
	assert.NotNil(result)
	assert.Equal(result.Name, data["name"])

}

func Test_CompanyRepository_GetOrCreateCompany_Get_Correct(t *testing.T) {
	assert := assert.New(t)

	company := testhelper.Helper_RandomCompany(t)

	repository := companies.NewCompanyRegistry().GetCompanyRepository()

	result, err := repository.GetOrCreateCompany(company.Domain, "", "", "", "", "", "", "", "", "")

	assert.Nil(err)
	assert.NotNil(result)
	assert.Equal(result.Name, company.Name)

}
