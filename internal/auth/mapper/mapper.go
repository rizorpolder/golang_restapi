package mapper

import (
	authpb "golang_restapi/contracts/auth/go"
	"golang_restapi/internal/auth/model"
)

func PbRegisterToRegister(register *authpb.RegisterRequest) model.Register {
	return model.Register{
		Login:    register.Login,
		Email:    register.Email,
		Password: register.Password,
	}
}

func PbLoginToLogin(login *authpb.LoginRequest) model.Login {
	return model.Login{
		LoginOrEmail: login.LoginOrEmail,
		Password:     login.Password,
	}
}

func TokenPairToPbTokenPair(pair model.TokenPair) *authpb.TokenPair {
	return &authpb.TokenPair{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
	}
}
