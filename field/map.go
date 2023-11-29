package field

type MapField map[string]string

var UserMapField MapField = MapField{
	"Username":  "user.username",
	"Password":  "user.password",
	"Email":     "user.email",
	"FirstName": "user.first_name",
	"LastName":  "user.last_name",
}

func (m *MapField) GetField(k string) string {
	if v, ok := (*m)[k]; ok {
		return v
	}
	return k
}
