package addresses

type AddressQuery struct {
	ID      uint
	StoreID uint
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
