package model

type Product struct {
	ID       string
	Name     string
	Category Category
	Fields   map[string]string
}

func NewProduct(ID string, Name string, Category Category, Fields map[string]string) *Product {
	return &Product{
		ID:       ID,
		Name:     Name,
		Category: Category,
		Fields:   Fields,
	}
}
