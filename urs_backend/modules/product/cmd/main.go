package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"urs_backend/modules/product"
	pb "urs_backend/proto/protobuf"

	_ "github.com/lib/pq"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func startGRPCServer(db *sql.DB) {
	grpcServer := grpc.NewServer()

	repo := product.NewProductRepository(db)

	// Register the ProductService handler with the DB connection
	pb.RegisterProductServiceServer(grpcServer, product.NewProductService(repo))

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Println("gRPC server is running on port :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func startRESTServer() {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	// Register the gRPC service with the HTTP proxy (gRPC-Gateway)
	err := pb.RegisterProductServiceHandlerFromEndpoint(
		context.Background(),
		mux,
		"localhost:50051", // gRPC server address
		opts,
	)
	if err != nil {
		log.Fatalf("could not register gRPC gateway: %v", err)
	}

	// Start the HTTP server to serve REST API
	http.ListenAndServe(":8080", mux)
}

func main() {
	// Set up the DB connection
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"localhost", "5432", "srinukesari", "srinukesari", "urs_products")

	var err error
	DB, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	defer DB.Close()

	// Run both servers concurrently
	go startGRPCServer(DB) // Start gRPC server
	startRESTServer()      // Start REST API server using gRPC-Gateway
}
