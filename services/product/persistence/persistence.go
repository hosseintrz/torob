package persistence

import "github.com/hosseintrz/torob/product/model"

type ProductRepository interface {
	CreateProduct(name string, category, imageUrl string, minPrice int32, fields map[string]string) (string, error)
	GetProduct(id string) (*model.Product, error)
	CreateCategory(name, parent, desc string) (string, error)
	GetCategory(name string) (*model.Category, error)
	GetProductsByType(category string) ([]*model.Product, error)
	GetCategories() ([]*model.Category, error)
}
