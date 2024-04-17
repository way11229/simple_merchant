package grpc

import (
	"github.com/way11229/simple_merchant/domain"
	pb "github.com/way11229/simple_merchant/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcHandler struct {
	pb.UnimplementedSimpleMerchantServer

	userService    domain.UserService
	authService    domain.AuthService
	productService domain.ProductService
}

func NewGrpcHandler(
	userService domain.UserService,
	authService domain.AuthService,
	productService domain.ProductService,
) *GrpcHandler {
	return &GrpcHandler{
		userService:    userService,
		authService:    authService,
		productService: productService,
	}
}

func (g *GrpcHandler) getResponseError(err error) error {
	var code codes.Code

	switch err {
	case domain.ErrPermissionDeny:
		code = codes.PermissionDenied
	case domain.ErrMissingRequiredParameter,
		domain.ErrInvalidUserName,
		domain.ErrInvalidEmail,
		domain.ErrInvalidUserPassword,
		domain.ErrEmailHasVerified,
		domain.ErrInvalidProductName:
		code = codes.InvalidArgument
	case domain.ErrLoginAborted,
		domain.ErrInvalidVerificationCode,
		domain.ErrVerificationCodeExpired,
		domain.ErrSendVerificationCodeInShortPeriod:
		code = codes.Aborted
	case domain.ErrUserEmailDuplicated:
		code = codes.AlreadyExists
	case domain.ErrRecordNotFound,
		domain.ErrEmailNotFound:
		code = codes.NotFound
	default:
		code = codes.Unknown
	}

	return status.Errorf(code, err.Error())
}
