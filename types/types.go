package types

type User struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Age int `json:"age"`
	Salary int `json:"salary"`
	Occupation string `json:"occupation"`
}
