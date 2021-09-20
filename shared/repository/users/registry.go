package users

import (
	"bit-jobs-api/shared/database"
	"bit-jobs-api/shared/repository/users/accountmanagement"
)

type UserRegistry struct {
}

func NewUserRegistry() *UserRegistry {
	return &UserRegistry{}
}

func (*UserRegistry) GetUserRepository() *accountmanagement.UserRepository {
	return accountmanagement.NewUserRepository(database.DB)
}
