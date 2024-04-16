package acceptance_tests

import (
	"context"
	"reflect"
	"testing"

	pb "github.com/way11229/simple_merchant/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

func Test_createProduct_ListTheRecommendedProducts(t *testing.T) {
	ctx := context.Background()

	client, closer := server(ctx)
	defer closer()

	productData := &pb.CreateProductRequest{
		Name:             "test",
		Description:      "test",
		Price:            100,
		OrderBy:          99999,
		IsRecommendation: true,
		TotalQuantity:    10,
		SoldQuantity:     0,
	}
	createProductResp, err := client.CreateProduct(ctx, productData)
	if err != nil {
		t.Fatalf("CreateProduct error = %v", err)
	}

	defer func() {
		if _, err := client.DeleteProductById(ctx, &pb.DeleteProductByIdRequest{
			ProductId: createProductResp.ProductId,
		}); err != nil {
			t.Fatalf("DeleteProductById error = %v", err)
		}
	}()

	ctxWithAuth := getCtxWithAuth(ctx)
	recommendedProductsResp, err := client.ListTheRecommendedProducts(ctxWithAuth, &emptypb.Empty{})
	if err != nil {
		t.Fatalf("ListTheRecommendedProducts error = %v", err)
	}

	if len(recommendedProductsResp.GetProducts()) == 0 {
		t.Fatal("ListTheRecommendedProducts response is empty")
	}

	product := recommendedProductsResp.GetProducts()[0]
	wantProduct := &pb.Product{
		Id:    createProductResp.ProductId,
		Name:  productData.Name,
		Price: product.Price,
	}
	if !reflect.DeepEqual(wantProduct, product) {
		t.Fatalf("recommended product want %v, got %v", wantProduct, product)
	}
}
