package domain

import "context"

type CreateProductParams struct {
	Name             string
	Description      string
	Price            uint32
	OrderBy          int32
	IsRecommendation bool
	TotalQuantity    uint32
	SoldQuantity     uint32
}

func (c *CreateProductParams) Validate() error {
	if len([]rune(c.Name)) > PRODUCT_NAME_MAX_LEN {
		return ErrInvalidUserName
	}

	return nil
}

type CreateProductResult struct {
	ProductId uint32
}

type DeleteProductByIdParams struct {
	ProductId uint32
}

type ListTheRecommendedProductsResult struct {
	Products []*RecommendedProduct
}

type RecommendedProduct struct {
	Id    uint32
	Name  string
	Price uint32
}

type ProductService interface {
	CreateProduct(ctx context.Context, input *CreateProductParams) (*CreateProductResult, error)
	DeleteProductById(ctx context.Context, input *DeleteProductByIdParams) error
	ListTheRecommendedProducts(ctx context.Context) (*ListTheRecommendedProductsResult, error)
}
