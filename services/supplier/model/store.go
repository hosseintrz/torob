package model

type Store struct {
	ID      string
	Name    string
	OwnerId string
	Url     string
	City    string
}

func NewStore(name, ownerId, url, city string) *Store {
	return &Store{
		Name:    name,
		OwnerId: ownerId,
		Url:     url,
		City:    city,
	}
}
