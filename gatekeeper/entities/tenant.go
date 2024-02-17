package entities

type Tenant struct {
	Tenant string   `json:"tenant" yaml:"tenant" gorm:"index"`
	Role   []string `json:"role" yaml:"role"`
}
