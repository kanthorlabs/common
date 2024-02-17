package entities

type Tenant struct {
	Tenant string   `json:"tenant" yaml:"tenant" gorm:"index"`
	Roles  []string `json:"role" yaml:"role"`
}
