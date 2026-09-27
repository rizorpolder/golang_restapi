package mapper

import (
	accountpb "golang_restapi/contracts/account/go"
	"golang_restapi/internal/gateway/model"
)

func PbsToUsers(pbs []*accountpb.User) []model.User {
	users := make([]model.User, len(pbs))
	for i, pb := range pbs {
		users[i] = PbToUser(pb)
	}
	return users
}

func PbToUser(userpb *accountpb.User) model.User {
	return model.User{
		ID:         userpb.Id,
		Login:      userpb.Login,
		Email:      userpb.Email,
		Phone:      userpb.Phone,
		FirstName:  userpb.FirstName,
		LastName:   userpb.LastName,
		MiddleName: userpb.MiddleName,
		Age:        userpb.Age,
		CreatedAt:  userpb.CreatedAt.AsTime(),
		UpdatedAt:  userpb.UpdatedAt.AsTime(),
	}
}

func UserCreateToPb(user model.User) *accountpb.CreateUser {
	return &accountpb.CreateUser{
		Login:      user.Login,
		Email:      user.Email,
		Phone:      user.Phone,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		MiddleName: user.MiddleName,
		Age:        user.Age,
	}
}

func UserUpdateToPb(user model.User) *accountpb.UpdateUser {
	return &accountpb.UpdateUser{
		Email:      user.Email,
		Phone:      user.Phone,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		MiddleName: user.MiddleName,
		Age:        user.Age,
	}
}

func CreateUserToUser(createUser model.CreateUser) model.User {
	return model.User{
		Login:      createUser.Login,
		Email:      createUser.Email,
		Phone:      createUser.Phone,
		FirstName:  createUser.FirstName,
		LastName:   createUser.LastName,
		MiddleName: createUser.MiddleName,
		Age:        createUser.Age,
	}
}
