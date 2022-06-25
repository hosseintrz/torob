package model

type Category struct {
	Name string
	Path string
}

func NewCategory(name, path string) *Category {
	return &Category{
		Name: name,
		Path: path,
	}
}
