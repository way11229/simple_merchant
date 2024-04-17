package grpc

import (
	"context"

	"github.com/way11229/simple_merchant/domain"
	pb "github.com/way11229/simple_merchant/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (g *GrpcHandler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	resp, err := g.productService.CreateProduct(ctx, &domain.CreateProductParams{
		Name:             req.GetName(),
		Description:      req.GetDescription(),
		Price:            req.GetPrice(),
		OrderBy:          req.GetOrderBy(),
		IsRecommendation: req.GetIsRecommendation(),
		TotalQuantity:    req.GetTotalQuantity(),
		SoldQuantity:     req.GetSoldQuantity(),
	})
	if err != nil {
		return nil, g.getResponseError(err)
	}

	return &pb.CreateProductResponse{
		ProductId: resp.ProductId,
	}, nil
}

func (g *GrpcHandler) DeleteProductById(ctx context.Context, req *pb.DeleteProductByIdRequest) (*emptypb.Empty, error) {
	if err := g.productService.DeleteProductById(ctx, &domain.DeleteProductByIdParams{
		ProductId: req.GetProductId(),
	}); err != nil {
		return nil, g.getResponseError(err)
	}

	return &emptypb.Empty{}, nil
}

func (g *GrpcHandler) ListTheRecommendedProducts(ctx context.Context, _ *emptypb.Empty) (*pb.ListTheRecommendedProductsResponse, error) {
	if _, err := g.authorizedUser(ctx); err != nil {
		return nil, g.getResponseError(err)
	}

	resp, err := g.productService.ListTheRecommendedProducts(ctx)
	if err != nil {
		return nil, g.getResponseError(err)
	}

	rtn := &pb.ListTheRecommendedProductsResponse{
		Products: []*pb.RecommendedProduct{},
	}

	for _, e := range resp.Products {
		rtn.Products = append(rtn.Products, &pb.RecommendedProduct{
			Id:    e.Id,
			Name:  e.Name,
			Price: e.Price,
		})
	}

	return rtn, nil
}
