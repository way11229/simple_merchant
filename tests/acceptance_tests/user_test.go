package acceptance_tests

import (
	"context"
	"testing"

	pb "github.com/way11229/simple_merchant/pb"
	grpcMetadata "google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

func Test_createUser_getUserEmailVerificationCode_verifyUserEmail_loginUser_logoutUser(t *testing.T) {
	ctx := context.Background()

	client, closer := server(ctx)
	defer closer()

	userData := &pb.CreateUserRequest{
		Name:     "tester",
		Email:    "test@gmail.com",
		Password: "aaaAAA!",
	}
	createUserResp, err := client.CreateUser(ctx, userData)
	if err != nil {
		t.Fatalf("CreateUser error = %v", err)
	}

	defer func() {
		if _, err := client.DeleteUserById(ctx, &pb.DeleteUserByIdRequest{
			UserId: createUserResp.GetUserId(),
		}); err != nil {
			t.Fatalf("DeleteUserById error = %v", err)
		}
	}()

	if _, err := client.GetUserEmailVerificationCode(ctx, &pb.GetUserEmailVerificationCodeRequest{
		Email: userData.Email,
	}); err != nil {
		t.Fatalf("GetUserEmailVerificationCode error = %v", err)
	}

	if mailerClient == nil {
		t.Fatal("mailerClient is nil")
	}

	if _, err := client.VerifyUserEmail(ctx, &pb.VerifyUserEmailRequest{
		Email:            userData.Email,
		VerificationCode: mailerClient.GetVerificationCode(),
	}); err != nil {
		t.Fatalf("GetUserEmailVerificationCode error = %v", err)
	}

	loginResp, err := client.LoginUser(ctx, &pb.LoginUserRequest{
		Email:    userData.Email,
		Password: userData.Password,
	})
	if err != nil {
		t.Fatalf("LoginUser error = %v", err)
	}

	if !loginResp.EmailHasVerified {
		t.Fatal("the email has not verified")
	}

	ctxWithAuth := grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+loginResp.Token)
	_, err = client.LogoutUser(ctxWithAuth, &emptypb.Empty{})
	if err != nil {
		t.Fatalf("LogoutUser error = %v", err)
	}
}
