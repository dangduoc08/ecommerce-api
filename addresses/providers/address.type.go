package providers

type AddressQuery struct {
	ID      uint
	StoreID uint
	Sort    string
	Order   string
	Limit   int
	Offset  int
}

type AddressCreation struct {
	StoreID    uint
	StreetName string
	LocationID *uint
}

type AddressUpdate struct {
	StoreID    uint
	StreetName string
	LocationID *uint
}
