package server

import (
	"care-backend/internal/auth"
	"care-backend/internal/models"
	pb "care-backend/internal/pb"
	"care-backend/internal/utils"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"strconv"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	db         *gorm.DB
	jwtManager *auth.JwtManager
	// Add any additional fields you need here (like a database connection)
}

func (s *UserService) GetUser(ctx context.Context, _ *emptypb.Empty) (*pb.User, error) {
	userId := ctx.Value("userId").(uint)

	var dbUser models.User
	//result := s.db.First(&dbUser, userId)
	result := s.db.Model(&models.User{}).Preload("Subscriptions", func(tx *gorm.DB) *gorm.DB {
		return tx.Limit(1).Order("expires_at desc")
	}).First(&dbUser, userId)

	if result.Error != nil {
		return nil, status.Errorf(codes.NotFound, "dbUser not found")
	} else {
		return &pb.User{Id: strconv.Itoa(int(dbUser.ID)), Name: dbUser.Name, KakaoId: dbUser.KakaoId}, nil
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	accessToken := req.KakaoAccessToken
	kakao, err := utils.GetUserInfoFromKakao(accessToken)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "Can't get user info from kakao")
	}

	user := models.User{Name: kakao.KakaoAccount.Name, KakaoId: kakao.Id, PhoneNo: kakao.KakaoAccount.PhoneNumber}
	result := s.db.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &pb.User{Id: strconv.Itoa(int(user.ID)), Name: user.Name, KakaoId: user.KakaoId}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	// Implement your logic here
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUser not implemented")
}

//func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.Empty, error) {
//	// Implement your logic here
//	return nil, status.Errorf(codes.Unimplemented, "method DeleteUser not implemented")
//}

func (s *UserService) GetUserByKakaoAccessToken(ctx context.Context, req *pb.GetUserByKakaoAccessTokenRequest) (*pb.User, error) {
	accessToken := req.KakaoAccessToken
	kakao, err := utils.GetUserInfoFromKakao(accessToken)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "Can't get user info from kakao")
	}

	var dbUser models.User
	result := s.db.First(&dbUser, "kakao_id=?", kakao.Id)

	if result.Error != nil {
		return nil, status.Errorf(codes.NotFound, "dbUser not found")
	} else {
		return &pb.User{Id: strconv.Itoa(int(dbUser.ID)), Name: dbUser.Name, KakaoId: dbUser.KakaoId, CreatedAt: timestamppb.New(dbUser.CreatedAt)}, nil
	}
}

func (s *UserService) GetJWTByAccessToken(ctx context.Context, req *pb.GetJWTByAccessTokenRequest) (*pb.JWT, error) {
	accessToken := req.KakaoAccessToken
	kakao, err := utils.GetUserInfoFromKakao(accessToken)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "Can't get user info from kakao")
	}

	var dbUser models.User
	result := s.db.First(&dbUser, "kakao_id=?", kakao.Id)

	if result.Error != nil {
		return nil, status.Errorf(codes.NotFound, "dbUser not found")
	} else {
		token, err := s.jwtManager.IssueJwtToken(int(dbUser.ID))
		if err != nil {
			return nil, status.Errorf(codes.Aborted, "error occurred while issuing jwt token")
		}
		return &pb.JWT{Jwt: token}, nil
	}
}

func NewUserServer(db *gorm.DB, jwtManager *auth.JwtManager) *UserService {
	s := &UserService{db: db, jwtManager: jwtManager}
	return s
}
