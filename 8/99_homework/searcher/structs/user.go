package structs

type User struct {
	Browsers interface{} `json:"browsers"`
	Company  interface{} `json:"company"`
	Country  interface{} `json:"country"`
	Email    interface{} `json:"email"`
	Job      interface{} `json:"job"`
	Name     interface{} `json:"name"`
}

type User1 struct {
	Browsers []string `json:"browsers"`
	Company  string   `json:"company"`
	Country  string   `json:"country"`
	Email    string   `json:"email"`
	Job      string   `json:"job"`
	Name     string   `json:"name"`
}
