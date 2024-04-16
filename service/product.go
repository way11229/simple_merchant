package service

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/way11229/simple_merchant/domain"
	mysql_sqlc "github.com/way11229/simple_merchant/repo/mysql/sqlc"
)

type ProductService struct {
	mysqlStore domain.MysqlStore
}

func NewProductService(
	mysqlStore domain.MysqlStore,
) domain.ProductService {
	return &ProductService{
		mysqlStore: mysqlStore,
	}
}

func (p *ProductService) CreateProduct(ctx context.Context, input *domain.CreateProductParams) (*domain.CreateProductResult, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	productId, err := p.mysqlStore.TxCreateProduct(ctx, &domain.MysqlTxCreateProductParams{
		CreateProductParams: mysql_sqlc.CreateProductParams{
			Name:             input.Name,
			Description:      input.Description,
			Price:            input.Price,
			OrderBy:          input.OrderBy,
			IsRecommendation: input.IsRecommendation,
			TotalQuantity:    input.TotalQuantity,
			SoldQuantity:     input.SoldQuantity,
			Status:           mysql_sqlc.ProductsStatusOn,
		},
	})
	if err != nil {
		return nil, err
	}

	return &domain.CreateProductResult{
		ProductId: productId,
	}, nil
}

func (p *ProductService) DeleteProductById(ctx context.Context, input *domain.DeleteProductByIdParams) error {
	product, err := p.mysqlStore.GetProductById(ctx, input.ProductId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrRecordNotFound
		}

		log.Printf("GetProductById error = %v, params = %d", err, input.ProductId)
		return domain.ErrUnknown
	}

	if err := p.mysqlStore.DeleteProductById(ctx, product.ID); err != nil {
		log.Printf("DeleteProductById error = %v, params = %d", err, product.ID)
		return domain.ErrUnknown
	}

	return nil
}

func (p *ProductService) ListTheRecommendedProducts(ctx context.Context) (*domain.ListTheRecommendedProductsResult, error) {
	products, err := p.batchListRecommendedProducts(ctx)
	if err != nil {
		return nil, err
	}

	return &domain.ListTheRecommendedProductsResult{
		Products: products,
	}, nil
}

/********************
 ********************
 ** private method **
 ********************
 ********************/

func (p *ProductService) batchListRecommendedProducts(ctx context.Context) ([]*domain.RecommendedProduct, error) {
	rtn := []*domain.RecommendedProduct{}
	offset := 0
	limit := 1000

	for {
		products, err := p.mysqlStore.ListTheRecommendedProducts(ctx, mysql_sqlc.ListTheRecommendedProductsParams{
			Offset: int32(offset),
			Limit:  int32(limit),
		})
		if err != nil {
			log.Printf("ListTheRecommendedProducts error = %v", err)
			return nil, domain.ErrUnknown
		}

		for _, e := range products {
			rtn = append(rtn, &domain.RecommendedProduct{
				Id:    e.ID,
				Name:  e.Name,
				Price: e.Price,
			})
		}

		numOfProducts := len(products)
		if numOfProducts < int(limit) || numOfProducts == 0 {
			break
		}

		offset += limit
	}

	return rtn, nil
}
