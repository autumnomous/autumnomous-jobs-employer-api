package accountmanagement_test

import (
	"testing"

	"autumnomous.com/bit-jobs-api/shared/repository/users"
	"autumnomous.com/bit-jobs-api/shared/testhelper"
	"github.com/stretchr/testify/assert"
)

func init() {
	testhelper.Init()
}
func Test_UserRepository_GetUser(t *testing.T) {

	assert := assert.New(t)

	repository := users.NewUserRegistry().GetUserRepository()

	applicant := testhelper.Helper_RandomApplicant(t)

	result, err := repository.GetUser(applicant.PublicID)

	assert.Nil(err)
	assert.NotNil(result)
}
