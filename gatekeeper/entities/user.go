package entities

type User struct {
	Username string   `json:"username" yaml:"username"`
	Role     []string `json:"role" yaml:"role"`
}
