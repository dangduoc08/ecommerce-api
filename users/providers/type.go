package providers

type Query struct {
	Username string
	Email    string
	StoreID  uint
	Statuses []string
	Sort     string
	Order    string
	Limit    int
	Offset   int
}

type Creation struct {
	StoreID   uint
	Password  string
	Username  string
	Email     string
	FirstName string
	LastName  string
	GroupIDs  []uint
}
