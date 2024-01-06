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
	StoreID uint
	Name    string
	Slug    string
	Logo    *string
}

type Update struct {
	Name string
	Slug string
	Logo *string
}
