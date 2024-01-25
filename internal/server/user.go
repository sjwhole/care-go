package server

import (
	"care-backend/internal/auth"
	"care-backend/internal/models"
	pb "care-backend/internal/pb"
	"care-backend/internal/utils"
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"strconv"
	"sync"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	db         *gorm.DB
	mu         sync.Mutex
	jwtManager *auth.JwtManager
	// Add any additional fields you need here (like a database connection)
}

func (s *UserService) GetUser(ctx context.Context, req *emptypb.Empty) (*pb.User, error) {

	md, _ := metadata.FromIncomingContext(ctx)
	fmt.Println(md["jwt"])

	userId, err := s.jwtManager.VerifyJwtToken(md["jwt"][0])
	if err != nil {
		return nil, err
	}
	fmt.Println(userId)

	s.mu.Lock()

	var dbUser models.User
	s.db.First(&dbUser, userId)

	s.mu.Unlock()

	if dbUser == (models.User{}) {
		return nil, status.Errorf(codes.NotFound, "dbUser not found")
	} else {
		return &pb.User{Id: strconv.Itoa(int(dbUser.ID)), Name: dbUser.Name, KakoId: dbUser.KakaoId}, nil
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
	return &pb.User{Id: strconv.Itoa(int(user.ID)), Name: user.Name, KakoId: user.KakaoId}, nil
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

	s.mu.Lock()

	var dbUser models.User
	s.db.First(&dbUser, "kakao_id=?", kakao.Id)

	s.mu.Unlock()

	if dbUser == (models.User{}) {
		return nil, status.Errorf(codes.NotFound, "dbUser not found")
	} else {
		return &pb.User{Id: strconv.Itoa(int(dbUser.ID)), Name: dbUser.Name, KakoId: dbUser.KakaoId}, nil
	}
}

func (s *UserService) GetJWTByAccessToken(ctx context.Context, req *pb.GetJWTByAccessTokenRequest) (*pb.JWT, error) {
	accessToken := req.KakaoAccessToken
	kakao, err := utils.GetUserInfoFromKakao(accessToken)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "Can't get user info from kakao")
	}

	s.mu.Lock()

	var dbUser models.User
	s.db.First(&dbUser, "kakao_id=?", kakao.Id)

	s.mu.Unlock()

	if dbUser == (models.User{}) {
		return nil, status.Errorf(codes.NotFound, "dbUser not found")
	} else {
		token, err := s.jwtManager.IssueJwtToken(int(dbUser.ID))
		if err != nil {
			return nil, status.Errorf(codes.Aborted, "error occurred while issuing jwt token")
		}
		return &pb.JWT{Jwt: token}, nil
	}
}

func NewServer(db *gorm.DB, jwtManager *auth.JwtManager) *UserService {
	s := &UserService{db: db, jwtManager: jwtManager}
	return s
}
