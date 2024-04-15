package grpc

import (
	"context"

	pb "github.com/way11229/simple_merchant/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (g *GrpcHandler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateProduct not implemented")
}

func (g *GrpcHandler) DeleteProductById(ctx context.Context, req *pb.DeleteProductByIdRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteProductById not implemented")
}

func (g *GrpcHandler) ListTheRecommendedProducts(ctx context.Context, _ *emptypb.Empty) (*pb.ListTheRecommendedProductsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListTheRecommendedProducts not implemented")
}
