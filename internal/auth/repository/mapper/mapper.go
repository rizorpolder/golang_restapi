package mapper

import (
	"golang_restapi/internal/auth/model"
	rpmodel "golang_restapi/internal/auth/repository/model"
)

func UserToRepoUser(userModel model.User) (rpUser rpmodel.User) {
	return rpmodel.User{
		ID:        userModel.ID,
		Login:     userModel.Login,
		Password:  userModel.PasswordHash,
		Email:     userModel.Email,
		CreatedAt: userModel.CreatedAt,
		UpdatedAt: userModel.UpdatedAt,
	}
}

func RepoUserToUser(userModel rpmodel.User) (user model.User) {
	return model.User{
		ID:           userModel.ID,
		Login:        userModel.Login,
		PasswordHash: userModel.Password,
		Email:        userModel.Email,
		CreatedAt:    userModel.CreatedAt,
		UpdatedAt:    userModel.UpdatedAt,
	}
}

func RepoTokenToToken(repoToken rpmodel.RefreshToken) model.RefreshToken {
	return model.RefreshToken{
		ID:        repoToken.ID,
		UserID:    repoToken.UserID,
		Token:     repoToken.Token,
		ExpiresAt: repoToken.ExpiresAt,
		RevokedAt: repoToken.RevokedAt,
		CreatedAt: repoToken.CreatedAt,
	}
}
func TokenToRepoToken(token model.RefreshToken) rpmodel.RefreshToken {
	return rpmodel.RefreshToken{
		ID:        token.ID,
		UserID:    token.UserID,
		Token:     token.Token,
		ExpiresAt: token.ExpiresAt,
		RevokedAt: token.RevokedAt,
		CreatedAt: token.CreatedAt,
	}
}
