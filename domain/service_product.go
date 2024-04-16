package domain

import "context"

const (
	RECOMMENDED_PRODUCT_CACHE_SORTED_SET_KEY_PREFIX = "zrp:"
	RECOMMENDED_PRODUCT_CACHE_STRING_KEY_PREFIX     = "rp:"
)

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
	Id    uint32 `json:"id"`
	Name  string `json:"name"`
	Price uint32 `json:"price"`

	OrderBy int32
}

type ProductService interface {
	CreateProduct(ctx context.Context, input *CreateProductParams) (*CreateProductResult, error)
	DeleteProductById(ctx context.Context, input *DeleteProductByIdParams) error
	ListTheRecommendedProducts(ctx context.Context) (*ListTheRecommendedProductsResult, error)
}
