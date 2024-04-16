package grpc

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/way11229/simple_merchant/domain"
	pb "github.com/way11229/simple_merchant/pb"
)

func (g *GrpcHandler) DeleteUserById(ctx context.Context, req *pb.DeleteUserByIdRequest) (*emptypb.Empty, error) {
	if err := g.userService.DeleteUserById(ctx, req.GetUserId()); err != nil {
		return nil, g.getResponseError(err)
	}

	return &emptypb.Empty{}, nil
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
		UserId: resp.UserId,
	}, nil
}

func (g *GrpcHandler) GetUserEmailVerificationCode(ctx context.Context, req *pb.GetUserEmailVerificationCodeRequest) (*emptypb.Empty, error) {
	if err := g.userService.GetUserEmailVerificationCode(ctx, &domain.GetUserEmailVerificationCodeParams{
		Email: req.GetEmail(),
	}); err != nil {
		return nil, g.getResponseError(err)
	}

	return &emptypb.Empty{}, nil
}

func (g *GrpcHandler) VerifyUserEmail(ctx context.Context, req *pb.VerifyUserEmailRequest) (*emptypb.Empty, error) {
	if err := g.userService.VerifyUserEmail(ctx, &domain.VerifyUserEmailParams{
		Email:                 req.GetEmail(),
		EmailVerificationCode: req.GetVerificationCode(),
		GetNowTimeFunc:        time.Now,
	}); err != nil {
		return nil, g.getResponseError(err)
	}

	return &emptypb.Empty{}, nil
}

func (g *GrpcHandler) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	resp, err := g.authService.LoginUser(ctx, &domain.LoginUserParams{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return nil, g.getResponseError(err)
	}

	return &pb.LoginUserResponse{
		Token: resp.Token,
	}, nil
}

func (g *GrpcHandler) LogoutUser(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	userId, err := g.authorizedUser(ctx)
	if err != nil {
		return nil, g.getResponseError(err)
	}

	if err := g.authService.LogoutUser(ctx, &domain.LogoutUserParams{
		UserId: userId,
	}); err != nil {
		return nil, g.getResponseError(err)
	}

	return &emptypb.Empty{}, nil
}
