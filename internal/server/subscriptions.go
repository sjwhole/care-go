package server

import (
	db "care-backend/db/gen"
	"care-backend/internal/auth"
	pb "care-backend/internal/pb"
	"context"
	"database/sql"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
)

type SubscriptionService struct {
	pb.UnimplementedSubscriptionServiceServer
	quries     *db.Queries
	jwtManager *auth.JwtManager
	// Add any additional fields you need here (like a database connection)
}

func subscriptionToPbSubscription(subscription *db.Subscription) *pb.Subscription {
	return &pb.Subscription{Id: strconv.FormatUint(subscription.ID, 10),
		CreatedAt: timestamppb.New(subscription.CreatedAt.Time),
		UpdatedAt: timestamppb.New(subscription.UpdatedAt.Time),
		DeletedAt: timestamppb.New(subscription.DeletedAt.Time),
		ExpiresAt: timestamppb.New(subscription.ExpiresAt.Time)}
}

func (s *SubscriptionService) GetSubscriptions(ctx context.Context, _ *emptypb.Empty) (*pb.SubscriptionList, error) {
	userId := ctx.Value("userId").(uint)

	subscriptions, err := s.quries.GetSubscriptionsByUserId(ctx, sql.NullInt64{Int64: int64(userId), Valid: true})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "subscriptions not found")
	} else {
		return &pb.SubscriptionList{Subscriptions: lo.Map(subscriptions, func(subscription db.Subscription, _ int) *pb.Subscription {
			return subscriptionToPbSubscription(&subscription)
		})}, nil
	}
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, req *pb.SubscriptionCreateRequest) (*pb.Subscription, error) {
	userId := ctx.Value("userId").(uint)

	result, err := s.quries.CreateSubscription(ctx, db.CreateSubscriptionParams{
		UserID:    sql.NullInt64{Int64: int64(uint64(userId)), Valid: true},
		ExpiresAt: sql.NullTime{Time: req.ExpiresAt.AsTime(), Valid: true},
	})
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "Can't create subscription")
	} else {
		id, _ := result.LastInsertId()
		subscription, err := s.quries.GetSubscriptionById(ctx, uint64(id))
		if err != nil {
			return nil, status.Errorf(codes.NotFound, "subscription not found")
		} else {
			return subscriptionToPbSubscription(&subscription), nil
		}
	}
}

func NewSubscriptionServer(quries *db.Queries, jwtManager *auth.JwtManager) *SubscriptionService {
	s := &SubscriptionService{quries: quries, jwtManager: jwtManager}
	return s
}
