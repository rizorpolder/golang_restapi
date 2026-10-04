package mapper

import (
	authpb "golang_restapi/contracts/auth/go"
	"golang_restapi/internal/gateway/model"
)

func PbToTokenPair(tokenPair *authpb.TokenPair) model.TokenPair {
	return model.TokenPair{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}
}

func TokenPairToPb(tokenPair model.TokenPair) *authpb.TokenPair {
	return &authpb.TokenPair{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}
}
