package client

const (
	UserPath    = "users"
	UserPathGet = UserPath + "/%s"
)

type User struct {
	EmailId   string `json:"emailId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password,omitempty"`
}
