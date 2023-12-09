package users

type UserQuery struct {
	Username string
	Email    string
}

type UserCreation struct {
	StoreID   uint
	Password  string
	Username  string
	Email     string
	FirstName string
	LastName  string
}
