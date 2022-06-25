package persistence

import "github.com/hosseintrz/torob/product/model"

type ProductRepository interface {
	CreateProduct(name string, category string, fields map[string]string) (string, error)
	GetProduct(id string) (*model.Product, error)
	CreateCategory(name, parent string) (string, error)
	GetCategory(name string) (*model.Category, error)
}
