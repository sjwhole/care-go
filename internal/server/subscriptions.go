package server

import (
	"care-backend/internal/auth"
	"care-backend/internal/models"
	pb "care-backend/internal/pb"
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/datatypes"
	"time"

	"github.com/samber/lo"

	"gorm.io/gorm"
	"strconv"
	"sync"
)

type SubscriptionService struct {
	pb.UnimplementedSubscriptionServiceServer
	db         *gorm.DB
	mu         sync.Mutex
	jwtManager *auth.JwtManager
	// Add any additional fields you need here (like a database connection)
}

func (s *SubscriptionService) GetSubscriptions(ctx context.Context, _ *emptypb.Empty) (*pb.SubscriptionList, error) {
	userId := ctx.Value("userId").(uint)

	fmt.Println(userId)

	s.mu.Lock()
	var dbSubscriptions []models.Subscription
	//result := s.db.First(&dbUser, userId)
	err := s.db.Model(&models.Subscription{}).Order("expires_at DESC").Find(&dbSubscriptions, "user_id = ?", userId).Error
	s.mu.Unlock()

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "dbUser not found")
	} else {
		return &pb.SubscriptionList{Subscriptions: lo.Map(dbSubscriptions, func(subscription models.Subscription, _ int) *pb.Subscription {
			return &pb.Subscription{Id: strconv.Itoa(int(subscription.ID)), ExpiresAt: timestamppb.New(time.Time(subscription.ExpiresAt))}
		})}, nil
	}
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, req *pb.SubscriptionCreateRequest) (*pb.Subscription, error) {
	userId := ctx.Value("userId").(uint)

	var dbSubscription models.Subscription
	dbSubscription.UserID = userId
	dbSubscription.ExpiresAt = datatypes.Date(req.ExpiresAt.AsTime())
	s.mu.Lock()
	result := s.db.Create(&dbSubscription)
	s.mu.Unlock()
	if result.Error != nil {
		return nil, status.Errorf(codes.Aborted, "Can't create subscription")
	} else {
		return &pb.Subscription{Id: strconv.Itoa(int(dbSubscription.ID)), ExpiresAt: timestamppb.New(time.Time(dbSubscription.ExpiresAt))}, nil
	}
}

func NewSubscriptionServer(db *gorm.DB, jwtManager *auth.JwtManager) *SubscriptionService {
	s := &SubscriptionService{db: db, jwtManager: jwtManager}
	return s
}
