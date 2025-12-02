package repositories

type ProductRepository struct{}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{}
}
