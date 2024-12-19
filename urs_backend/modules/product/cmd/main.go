package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"your-app/handler" // Replace with your actual Go module path
	"your-app/product" // Replace with your actual Go module path

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jmoiron/sqlx" // Import sqlx for DB interaction
	_ "github.com/lib/pq"     // Import PostgreSQL driver (or other DB driver)
	"google.golang.org/grpc"
)

func startGRPCServer(db *sqlx.DB) {
	grpcServer := grpc.NewServer()

	// Register the ProductService handler with the DB connection
	product.RegisterProductServiceServer(grpcServer, handler.NewProductServiceHandler(db))

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
	err := product.RegisterProductServiceHandlerFromEndpoint(
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
	// Set up the DB connection (PostgreSQL example using sqlx)
	db, err := sqlx.Connect("postgres", "user=youruser dbname=yourdb sslmode=disable")
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run both servers concurrently
	go startGRPCServer(db) // Start gRPC server
	startRESTServer()      // Start REST API server using gRPC-Gateway
}
