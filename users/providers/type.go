package providers

type UserQuery struct {
	Username string
	Email    string
	StoreID  uint
	Statuses []string
	Sort     string
	Order    string
	Limit    int
	Offset   int
}

type UserCreation struct {
	StoreID   uint
	Password  string
	Username  string
	Email     string
	FirstName string
	LastName  string
	GroupIDs  []uint
}
