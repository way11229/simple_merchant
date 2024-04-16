package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/way11229/simple_merchant/domain"
	pb "github.com/way11229/simple_merchant/pb"
)

func (g *GrpcHandler) DeleteUserById(ctx context.Context, req *pb.DeleteUserByIdRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUserById not implemented")
}

func (g *GrpcHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	resp, err := g.userService.CreateUser(ctx, &domain.CreateUserParams{
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return nil, g.getResponseError(err)
	}

	return &pb.CreateUserResponse{
		UserId: uint64(resp.UserId),
	}, nil
}

func (g *GrpcHandler) GetUserEmailVerificationCode(ctx context.Context, req *pb.GetUserEmailVerificationCodeRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserEmailVerificationCode not implemented")
}

func (g *GrpcHandler) VerifyUserEmail(ctx context.Context, req *pb.VerifyUserEmailRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyUserEmail not implemented")
}

func (g *GrpcHandler) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginUser not implemented")
}

func (g *GrpcHandler) LogoutUser(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LogoutUser not implemented")
}
