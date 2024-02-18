package interceptors

import (
	"care-backend/internal/auth"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
)

var exemptMethods = map[string]bool{
	"/user.UserService/GetJWTByAccessToken": true,
	"/user.UserService/CreateUser":          true,
}

func AuthInterceptor(jwtManager *auth.JwtManager) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if exemptMethods[info.FullMethod] {
			// Allow exempt requests without authentication
			return handler(ctx, req)
		}
		// Extract JWT from metadata
		md, _ := metadata.FromIncomingContext(ctx)
		//jwt := md["authorization"]
		log.Println(md)

		log.Println(info.FullMethod)

		// Validate JWT (example using a hypothetical validation function)
		//if userId,  err := auth.JwtManager.VerifyJwtToken(jwt); err != nil {
		//	return nil, status.Errorf(codes.Unauthenticated, "Invalid authentication token")
		//}
		if len(md["jwt"]) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "Invalid authentication token")
		}
		userId, err := jwtManager.VerifyJwtToken(md["jwt"][0])
		if err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, "userId", userId)

		// Proceed with the original RPC call

		return handler(ctx, req)
	}
}
