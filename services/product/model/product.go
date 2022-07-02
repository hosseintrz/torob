package model

type Product struct {
	ID       string
	Name     string
	Category Category
	ImageUrl string
	MinPrice int32
	Fields   map[string]string
}

func NewProduct(ID string, Name string, Category Category, ImageUrl string, MinPrice int32,
	Fields map[string]string) *Product {
	return &Product{
		ID:       ID,
		Name:     Name,
		Category: Category,
		ImageUrl: ImageUrl,
		MinPrice: MinPrice,
		Fields:   Fields,
	}
}
