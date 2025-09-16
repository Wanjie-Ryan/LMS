package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/Wanjie-Ryan/LMS/cmd/api/handlers"
	"github.com/Wanjie-Ryan/LMS/cmd/api/services"
	"github.com/Wanjie-Ryan/LMS/common"
	auth "github.com/Wanjie-Ryan/LMS/genproto"
)

func StartGRPCServer() {

	db, err := common.ConnectionDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	redisClient := common.ConnectRedis()
	authService := services.NewAuthService(db, redisClient)
	grpcHandler := &handlers.GrpcAuthHandler{Service: authService}

	lis, err := net.Listen("tcp", ":50051")
	if err !=nil{
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	// pb.RegisterAuthServiceServer(grpcServer, grpcHandler)
	auth.RegisterAuthServiceServer(grpcServer, grpcHandler)
	log.Default().Println("gRPC server running on port 50051")
	grpcServer.Serve(lis)

}
