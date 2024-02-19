package main

import (
	"care-backend/internal/auth"
	"care-backend/internal/interceptors"
	"care-backend/internal/models"
	pb "care-backend/internal/pb"
	"care-backend/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"net"
)

func main() {
	// Load jwt manager
	jwtManager := auth.InitliazeJWTManager()

	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	//dsn := "root@tcp(127.0.0.1:3306)/care?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := "root:root@tcp(mysql:3306)/care?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Log every SQL query
	})
	if err != nil {
		panic(err.Error())
	}

	// Migrate the schema
	err = db.AutoMigrate(&models.User{}, &models.Subscription{})
	if err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}

	// Listen on a specific host and port
	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	//Print listen address
	log.Println("Server is on " + lis.Addr().String())

	// Create a new gRPC server
	s := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.AuthInterceptor(jwtManager)),
	)

	// Register the service with the server
	pb.RegisterUserServiceServer(s, server.NewUserServer(db, jwtManager))
	pb.RegisterSubscriptionServiceServer(s, server.NewSubscriptionServer(db, jwtManager))

	reflection.Register(s)

	// Serve the server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
