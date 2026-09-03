package mapper

import (
	"golang_restapi/internal/account/model"
	repomodel "golang_restapi/internal/repository/account/model"
)

func UserToRepoUser(user model.User) repomodel.User {
	return repomodel.User{
		ID:         user.ID,
		Email:      user.Email,
		Phone:      user.Phone,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		MiddleName: user.MiddleName,
		Age:        user.Age,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
}

func RepoUserToUser(user repomodel.User) model.User {
	return model.User{
		ID:         user.ID,
		Email:      user.Email,
		Phone:      user.Phone,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		MiddleName: user.MiddleName,
		Age:        user.Age,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
}

func RepoUsersToUsers(users []repomodel.User) []model.User {
	res := make([]model.User, len(users))
	for i, user := range users {
		res[i] = RepoUserToUser(user)
	}
	return res
}

func UsersToRepoUsers(users []model.User) []repomodel.User {
	res := make([]repomodel.User, len(users))
	for i, user := range users {
		res[i] = UserToRepoUser(user)
	}
	return res
}
