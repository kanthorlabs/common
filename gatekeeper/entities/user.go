package entities

type User struct {
	Username string   `json:"username" yaml:"username"`
	Roles    []string `json:"role" yaml:"role"`
}
