package main

import (
	"care-backend/internal/auth"
	"care-backend/internal/models"
	pb "care-backend/internal/pb"
	"care-backend/internal/server"
	"context"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net"
	"os"
	"strconv"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	// Load env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	secretKey := os.Getenv("SECRET_KEY")
	validMin, err := strconv.Atoi(os.Getenv("VALID_MIN"))
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Create a new JWT manager
	jwtManager := auth.NewJwtManager(secretKey, validMin)

	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "root@tcp(127.0.0.1:3306)/care?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info), // Log every SQL query
	})
	if err != nil {
		panic(err.Error())
	}

	// Migrate the schema
	err = db.AutoMigrate(&Product{}, &models.User{})
	if err != nil {
		panic(err.Error())
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
		grpc.UnaryInterceptor(authInterceptor),
	)

	// Register the service with the server
	pb.RegisterUserServiceServer(s, server.NewServer(db, jwtManager))

	// Serve the server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func authInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	// Extract JWT from metadata
	md, _ := metadata.FromIncomingContext(ctx)
	//jwt := md["authorization"]
	log.Println(md)

	// Validate JWT (example using a hypothetical validation function)
	//if userId,  err := auth.JwtManager.VerifyJwtToken(jwt); err != nil {
	//	return nil, status.Errorf(codes.Unauthenticated, "Invalid authentication token")
	//}

	// Proceed with the original RPC call
	return handler(ctx, req)
}
