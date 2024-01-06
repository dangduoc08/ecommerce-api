package providers

type Query struct {
	ID      uint
	StoreID uint
	Sort    string
	Order   string
	Limit   int
	Offset  int
}

type Creation struct {
	StoreID    uint
	StreetName string
	LocationID *uint
}

type Update struct {
	StoreID    uint
	StreetName string
	LocationID *uint
}
