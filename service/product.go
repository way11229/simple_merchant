package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/way11229/simple_merchant/domain"
	mysql_sqlc "github.com/way11229/simple_merchant/repo/mysql/sqlc"
)

type ProductService struct {
	mysqlStore domain.MysqlStore

	redisClient domain.RedisClient

	recommendedProductCacheExpired time.Duration
}

func NewProductService(
	mysqlStore domain.MysqlStore,
	redisClient domain.RedisClient,
	recommendedProductCacheExpired time.Duration,
) domain.ProductService {
	return &ProductService{
		mysqlStore:                     mysqlStore,
		redisClient:                    redisClient,
		recommendedProductCacheExpired: recommendedProductCacheExpired,
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

	if input.IsRecommendation {
		go p.batchSetRecommendedProductsCache(context.Background(), &batchSetRecommendedProductsCacheParams{
			Products: []*domain.RecommendedProduct{
				{
					Id:      productId,
					Name:    input.Name,
					Price:   input.Price,
					OrderBy: input.OrderBy,
				},
			},
		})
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

	go p.removeRecommendedProductCacheByProductId(context.Background(), product.ID)

	return nil
}

func (p *ProductService) ListTheRecommendedProducts(ctx context.Context) (*domain.ListTheRecommendedProductsResult, error) {
	products, err := p.listRecommendedProductsWithCacheAndDB(ctx)
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

type batchSetRecommendedProductsCacheParams struct {
	Products []*domain.RecommendedProduct
}

func (p *ProductService) batchSetRecommendedProductsCache(ctx context.Context, input *batchSetRecommendedProductsCacheParams) error {
	redisZAddParams := &domain.ZAddParams{
		Key:     domain.RECOMMENDED_PRODUCT_CACHE_SORTED_SET_KEY_PREFIX,
		Members: []*domain.ZScoreMember{},
	}

	for _, product := range input.Products {
		productJson, err := json.Marshal(product)
		if err != nil {
			log.Printf("json.Marshal error = %v, params = %v", err, product)
			return domain.ErrUnknown
		}

		rpKey := p.getRecommendedProductStringCacheKey(product.Id)
		if err := p.redisClient.SetEx(ctx, &domain.SetExParams{
			Key:        rpKey,
			Value:      string(productJson),
			Expiration: p.recommendedProductCacheExpired,
		}); err != nil {
			return err
		}

		redisZAddParams.Members = append(redisZAddParams.Members, &domain.ZScoreMember{
			Score:  float64(product.OrderBy),
			Member: rpKey,
		})
	}

	if err := p.redisClient.ZAdd(ctx, redisZAddParams); err != nil {
		return err
	}

	if err := p.redisClient.SetExpiration(ctx, &domain.SetExpirationParams{
		Key:        redisZAddParams.Key,
		Expiration: p.recommendedProductCacheExpired,
	}); err != nil {
		p.redisClient.Del(ctx, &domain.DelParams{
			Keys: []string{redisZAddParams.Key},
		})

		return err
	}

	return nil
}

func (p *ProductService) getRecommendedProductStringCacheKey(productId uint32) string {
	return fmt.Sprintf("%s%d", domain.RECOMMENDED_PRODUCT_CACHE_STRING_KEY_PREFIX, productId)
}

func (p *ProductService) removeRecommendedProductCacheByProductId(ctx context.Context, productId uint32) error {
	rpKey := p.getRecommendedProductStringCacheKey(productId)

	if err := p.redisClient.ZRem(ctx, &domain.ZRemParams{
		Key: domain.RECOMMENDED_PRODUCT_CACHE_SORTED_SET_KEY_PREFIX,
		Members: []interface{}{
			rpKey,
		},
	}); err != nil {
		return err
	}

	if err := p.redisClient.Del(ctx, &domain.DelParams{
		Keys: []string{rpKey},
	}); err != nil {
		return err
	}

	return nil
}

func (p *ProductService) listRecommendedProductsFromCache(ctx context.Context) ([]*domain.RecommendedProduct, error) {
	cacheKeyList, err := p.redisClient.ZRevRange(ctx, &domain.ZRevRangeParams{
		Key:   domain.RECOMMENDED_PRODUCT_CACHE_SORTED_SET_KEY_PREFIX,
		Start: 0,
		Stop:  -1, // list all
	})
	if err != nil {
		return nil, err
	}

	rtn := []*domain.RecommendedProduct{}
	for _, cacheKey := range cacheKeyList {
		cacheProduct, err := p.redisClient.Get(ctx, &domain.GetParams{
			Key: cacheKey,
		})
		if err != nil {
			return nil, err
		}

		var recommendedProduct domain.RecommendedProduct
		if err := json.Unmarshal([]byte(cacheProduct), &recommendedProduct); err != nil {
			log.Printf("json.Unmarshal error = %v, params = %v", err, cacheProduct)
			return nil, domain.ErrUnknown
		}

		rtn = append(rtn, &recommendedProduct)
	}

	return rtn, nil
}

func (p *ProductService) listRecommendedProductsWithCacheAndDB(ctx context.Context) ([]*domain.RecommendedProduct, error) {
	cacheProducts, _ := p.listRecommendedProductsFromCache(ctx)
	if len(cacheProducts) > 0 {
		return cacheProducts, nil
	}

	products, err := p.batchListRecommendedProducts(ctx)
	if err != nil {
		return nil, err
	}

	if len(products) > 0 {
		go p.batchSetRecommendedProductsCache(context.Background(), &batchSetRecommendedProductsCacheParams{
			Products: products,
		})
	}

	return products, nil
}
