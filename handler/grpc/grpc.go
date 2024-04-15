package grpc

import (
	pb "github.com/way11229/simple_merchant/pb"
)

type GrpcHandler struct {
	pb.UnimplementedSimpleMerchantServer
}

func NewGrpcHandler() *GrpcHandler {
	return &GrpcHandler{}
}
