package acceptance_tests

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	pb "github.com/way11229/simple_merchant/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

func Test_createProduct_listTheRecommendedProducts(t *testing.T) {
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
	wantProduct := &pb.RecommendedProduct{
		Id:    createProductResp.ProductId,
		Name:  productData.Name,
		Price: productData.Price,
	}
	if !reflect.DeepEqual(wantProduct, product) {
		t.Fatalf("recommended product want %v, got %v", wantProduct, product)
	}
}

func Test_recommendedProductsSorting(t *testing.T) {
	ctx := context.Background()

	client, closer := server(ctx)
	defer closer()

	deleteProductIds := []uint32{}

	type productTmp struct {
		Id   uint32
		Data *pb.CreateProductRequest
	}

	productDataList := []*productTmp{}
	for i := 10; i > 0; i -= 1 {
		productData := &pb.CreateProductRequest{
			Name:             fmt.Sprintf("test_%d", i),
			Description:      "test",
			Price:            100,
			OrderBy:          0,
			IsRecommendation: true,
			TotalQuantity:    10,
			SoldQuantity:     0,
		}

		if i == 10 {
			productData.OrderBy = 99999
		}

		createProductResp, err := client.CreateProduct(ctx, productData)
		if err != nil {
			t.Fatalf("CreateProduct error = %v", err)
		}

		deleteProductIds = append(deleteProductIds, createProductResp.ProductId)
		productDataList = append(productDataList, &productTmp{
			Id:   createProductResp.ProductId,
			Data: productData,
		})
	}

	defer func() {
		for _, productId := range deleteProductIds {
			if _, err := client.DeleteProductById(ctx, &pb.DeleteProductByIdRequest{
				ProductId: productId,
			}); err != nil {
				t.Fatalf("DeleteProductById error = %v", err)
			}
		}
	}()

	// waiting for cache
	time.Sleep(time.Second)

	ctxWithAuth := getCtxWithAuth(ctx)
	recommendedProductsResp, err := client.ListTheRecommendedProducts(ctxWithAuth, &emptypb.Empty{})
	if err != nil {
		t.Fatalf("ListTheRecommendedProducts error = %v", err)
	}

	if len(recommendedProductsResp.GetProducts()) == 0 {
		t.Fatal("ListTheRecommendedProducts response is empty")
	}

	product := recommendedProductsResp.GetProducts()[0]
	wantProduct := &pb.RecommendedProduct{
		Id:    productDataList[0].Id,
		Name:  productDataList[0].Data.Name,
		Price: productDataList[0].Data.Price,
	}
	if !reflect.DeepEqual(wantProduct, product) {
		t.Fatalf("recommended product want %v, got %v", wantProduct, product)
	}
}
