package structs

type User struct {
	Browsers interface{} `json:"browsers"`
	Company  interface{} `json:"company"`
	Country  interface{} `json:"country"`
	Email    interface{} `json:"email"`
	Job      interface{} `json:"job"`
	Name     interface{} `json:"name"`
}
