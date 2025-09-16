package handlers

import (
	"context"
	"fmt"

	"github.com/Wanjie-Ryan/LMS/cmd/api/requests"
	"github.com/Wanjie-Ryan/LMS/cmd/api/services"
	pb "github.com/Wanjie-Ryan/LMS/genproto"
)

type GrpcAuthHandler struct {
	// embed the uninimplemented version so we don't have to write all methods at once
	pb.UnimplementedAuthServiceServer
	Service services.AuthService
}

// gRPC Register method
// ctx is the request context for timeouts
func (h *GrpcAuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {

	// converting proto enum to string
	var role string
	switch req.Role {
	case pb.Role_ADMIN:
		role = "admin"
	case pb.Role_MEMBER:
		role = "member"
	default:
		role = "member"
	}

	// convert proto input to existing DTO

	payload := &requests.RegisterRequest{
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Email:     req.Email,
		Password:  req.Password,
		Role:      requests.Role(role),
	}

	// check if user already exists by mail
	existingUser, err := h.Service.GetUserByMail(payload.Email)

	if err != nil {
		fmt.Println("error getting user by mail", err)
		return nil, err
	}

	if existingUser != nil {
		fmt.Println("user already exists", err)
		return nil, err
	}

	// call the registerService
	user, err := h.Service.RegisterService(payload)

	if err != nil {
		fmt.Println("error registering user", err)
		return nil, err
	}

	// return safe user info
	userView := &pb.UserView{
		Id:        uint64(user.ID),
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		Role:      string(user.Role),
	}

	return &pb.RegisterResponse{
		Message: "Registration successful",
		User:    userView,
	}, nil

}
