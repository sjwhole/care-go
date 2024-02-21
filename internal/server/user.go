package server

import (
	db "care-backend/db/gen"
	"care-backend/internal/auth"
	pb "care-backend/internal/pb"
	"care-backend/internal/utils"
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	queries    *db.Queries
	jwtManager *auth.JwtManager
	// Add any additional fields you need here (like a database connection)
}

func userToPbUser(user *db.User) *pb.User {
	return &pb.User{Id: strconv.FormatUint(user.ID, 10), CreatedAt: timestamppb.New(user.CreatedAt.Time), Name: user.Name.String, KakaoId: uint64(user.KakaoID.Int64)}
}

func (s *UserService) GetUser(ctx context.Context, _ *emptypb.Empty) (*pb.User, error) {
	userId := ctx.Value("userId").(uint)

	user, err := s.queries.GetUserById(ctx, uint64(userId))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "dbUser not found")
	} else {
		return userToPbUser(&user), nil
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	accessToken := req.KakaoAccessToken
	kakao, err := utils.GetUserInfoFromKakao(accessToken)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "Can't get user info from kakao")
	}

	result, err := s.queries.CreateUser(ctx, db.CreateUserParams{
		KakaoID: sql.NullInt64{Int64: int64(kakao.Id), Valid: true},
		Name:    sql.NullString{String: kakao.KakaoAccount.Name, Valid: true},
		PhoneNo: sql.NullString{String: kakao.KakaoAccount.PhoneNumber, Valid: true},
	})
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "Can't create user")
	} else {
		id, _ := result.LastInsertId()
		user, err := s.queries.GetUserById(ctx, uint64(id))
		if err != nil {
			return nil, status.Errorf(codes.NotFound, "dbUser not found")
		} else {
			return userToPbUser(&user), nil
		}
	}
}

//	func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
//		// Implement your logic here
//		return nil, status.Errorf(codes.Unimplemented, "method UpdateUser not implemented")
//	}
//
// //func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.Empty, error) {
// //	// Implement your logic here
// //	return nil, status.Errorf(codes.Unimplemented, "method DeleteUser not implemented")
// //}
//
//	func (s *UserService) GetUserByKakaoAccessToken(ctx context.Context, req *pb.GetUserByKakaoAccessTokenRequest) (*pb.User, error) {
//		accessToken := req.KakaoAccessToken
//		kakao, err := utils.GetUserInfoFromKakao(accessToken)
//		if err != nil {
//			return nil, status.Errorf(codes.Aborted, "Can't get user info from kakao")
//		}
//
//		user, err := s.queries.GetUserByKakaoId(ctx, kakao.Id)
//		if err != nil {
//			return nil, status.Errorf(codes.NotFound, "dbUser not found")
//		} else {
//			return userToPbUser(&user), nil
//		}
//	}
func (s *UserService) GetJWTByAccessToken(ctx context.Context, req *pb.GetJWTByAccessTokenRequest) (*pb.JWT, error) {
	accessToken := req.KakaoAccessToken
	kakao, err := utils.GetUserInfoFromKakao(accessToken)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "Can't get user info from kakao")
	}

	user, err := s.queries.GetUserByKakaoId(ctx, sql.NullInt64{Int64: int64(kakao.Id), Valid: true})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "dbUser not found")
	} else {
		token, err := s.jwtManager.IssueJwtToken(int(user.ID))
		if err != nil {
			return nil, status.Errorf(codes.Aborted, "error occurred while issuing jwt token")
		}
		return &pb.JWT{Jwt: token}, nil
	}
}

func NewUserServer(queries *db.Queries, jwtManager *auth.JwtManager) *UserService {
	s := &UserService{queries: queries, jwtManager: jwtManager}
	return s
}
