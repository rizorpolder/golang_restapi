package mapper

import (
	accountpb "golang_restapi/contracts/account/go"
	"golang_restapi/internal/gateway/model"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func UserToPb(user model.User) *accountpb.User {
	return &accountpb.User{
		Id:         user.ID,
		Login:      user.Login,
		Email:      user.Email,
		Phone:      user.Phone,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		MiddleName: user.MiddleName,
		Age:        user.Age,
		CreatedAt:  timestamppb.New(user.CreatedAt),
		UpdatedAt:  timestamppb.New(user.UpdatedAt),
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
