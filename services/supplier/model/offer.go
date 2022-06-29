package model

type Offer struct {
	ID          string
	StoreId     string
	ProductId   string
	Price       int32
	Url         string
	Description string
}

func NewOffer(storeId, productId, url, description string, price int32) *Offer {
	return &Offer{
		StoreId:     storeId,
		ProductId:   productId,
		Price:       price,
		Url:         url,
		Description: description,
	}
}
