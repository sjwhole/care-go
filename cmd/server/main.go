package main

import (
	db "care-backend/db/gen"
	"care-backend/internal/auth"
	"care-backend/internal/interceptors"
	pb "care-backend/internal/pb"
	"care-backend/internal/server"
	"database/sql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	// Load jwt manager
	jwtManager := auth.InitliazeJWTManager()

	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "root@tcp(127.0.0.1:3306)/care?charset=utf8mb4&parseTime=True&loc=Local"
	//dsn := "root:root@tcp(mysql:3306)/care?charset=utf8mb4&parseTime=True&loc=Local"
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to open a DB connection: %v", err)
	}
	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {
			log.Printf("Failed to close the DB connection: %v\n", err)
		}
	}(conn)

	queries := db.New(conn)

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
	pb.RegisterUserServiceServer(s, server.NewUserServer(queries, jwtManager))
	pb.RegisterSubscriptionServiceServer(s, server.NewSubscriptionServer(queries, jwtManager))

	reflection.Register(s)

	// Serve the server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
