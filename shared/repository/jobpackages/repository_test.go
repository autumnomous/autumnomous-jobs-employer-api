package jobpackages_test

import (
	"bit-jobs-api/shared/repository/jobpackages"
	"bit-jobs-api/shared/testhelper"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	testhelper.Init()
}

func Test_EmployerRepository_GetActiveJobPackages_Correct(t *testing.T) {
	assert := assert.New(t)

	repository := jobpackages.NewJobPackageRegistry().GetJobPackageRepository()

	testhelper.Helper_RandomJobPackage(t)

	result, err := repository.GetActiveJobPackages()

	assert.Nil(err)
	assert.NotNil(result)
	assert.GreaterOrEqual(len(result), 1)

}

func Test_EmployerRepository_GetJobPackage_Correct(t *testing.T) {
	assert := assert.New(t)

	repository := jobpackages.NewJobPackageRegistry().GetJobPackageRepository()

	jobpackage := testhelper.Helper_RandomJobPackage(t)

	result, err := repository.GetJobPackage(jobpackage.TypeID)

	assert.Nil(err)
	assert.NotNil(result)
	assert.Equal(result.TypeID, jobpackage.TypeID)

}
