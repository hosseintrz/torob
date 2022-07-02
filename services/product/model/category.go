package model

type Category struct {
	Name        string
	Path        string
	Description string
}

func NewCategory(name, path, desc string) *Category {
	return &Category{
		Name:        name,
		Path:        path,
		Description: desc,
	}
}
