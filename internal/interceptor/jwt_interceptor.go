package interceptor

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

type JWTClaims struct {
	UserId uint64 `json:"user_id"`
	jwt.RegisteredClaims
}

type JWTInterceptor struct {
	jwtSecret []byte
}

const (
	AuthorizationHeader = "authorization"
	BearerPrefix        = "Bearer "
	UserIDKey           = "user_id"
)

var publicMethods = map[string]bool{
	"/gateway.Gateway/Register": true,
	"/gateway.Gateway/Login":    true,
	"/gateway.Gateway/Refresh":  true,
}

func NewJWTInterceptor(jwtSecret string) *JWTInterceptor {
	return &JWTInterceptor{jwtSecret: []byte(jwtSecret)}
}

func (i *JWTInterceptor) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}
		userId, err := i.extractUserIdFromJWT(ctx)
		if err != nil {
			return nil, err
		}

		ctxWithUserID := context.WithValue(ctx, UserIDKey, userId)
		return handler(ctxWithUserID, req)
	}
}

func (i *JWTInterceptor) StreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if publicMethods[info.FullMethod] {
			return handler(srv, stream)
		}

		userId, err := i.extractUserIdFromJWT(stream.Context())
		if err != nil {
			return err
		}

		ctxWithUserId := context.WithValue(stream.Context(), UserIDKey, userId)
		wrappedStream := &wrappedServerStream{
			ServerStream: stream,
			ctx:          ctxWithUserId,
		}
		return handler(srv, wrappedStream)
	}
}

func (i *JWTInterceptor) extractUserIdFromJWT(ctx context.Context) (uint64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	authHeaders := md.Get(AuthorizationHeader)
	if len(authHeaders) == 0 {
		return 0, status.Errorf(codes.Unauthenticated, "authorization header is not provided")
	}

	token := authHeaders[0]
	if !strings.HasPrefix(token, BearerPrefix) {
		return 0, status.Errorf(codes.Unauthenticated, "invalid token format")
	}
	accessToken := strings.TrimPrefix(token, BearerPrefix)
	if accessToken == "" {
		return 0, status.Errorf(codes.Unauthenticated, "access token is empty")
	}

	userId, err := i.parseAndValidateJWT(accessToken)
	if err != nil {
		return 0, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	return userId, nil
}

func (i *JWTInterceptor) parseAndValidateJWT(tokenString string) (uint64, error) {
	claimsNew := JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claimsNew, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return i.jwtSecret, nil
	})
	if err != nil {
		return 0, fmt.Errorf("failed to parse token: %w", err)
	}
	if !token.Valid {
		return 0, fmt.Errorf("token is not valid")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return 0, fmt.Errorf("fialed to extract claims")
	}
	if claims.UserId == 0 {
		return 0, fmt.Errorf("user_id not found")
	}
	return claims.UserId, nil
}

func (i *JWTInterceptor) ValidateTokenWithoutAuthService(tokenString string) (uint64, error) {
	return i.parseAndValidateJWT(tokenString)
}

func GetUserIDFromContext(ctx context.Context) (uint64, error) {
	userId, ok := ctx.Value(UserIDKey).(uint64)
	if !ok {
		return 0, fmt.Errorf("user_id not found in context")
	}
	return userId, nil
}

type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}
