package entities

import (
	"github.com/kanthorlabs/common/project"
	"github.com/kanthorlabs/common/safe"
	"github.com/kanthorlabs/common/validator"
)

var AnyTenant = "*"

type Privilege struct {
	Id       string         `json:"id" yaml:"id" gorm:"primaryKey"`
	Tenant   string         `json:"tenant" yaml:"tenant" gorm:"index"`
	Username string         `json:"username" yaml:"username" gorm:"index"`
	Role     string         `json:"role" yaml:"role"`
	Metadata *safe.Metadata `json:"metadata" yaml:"metadata"`

	CreatedAt int64 `json:"created_at" yaml:"created_at"`
	UpdatedAt int64 `json:"updated_at" yaml:"updated_at"`
}

func (privilege *Privilege) TableName() string {
	return project.Name("opm_privilege")
}

func (privilege *Privilege) Validate() error {
	return validator.Validate(
		validator.StringRequired("GATEKEEPER.PREVILEGE.ID", privilege.Id),
		validator.StringRequired("GATEKEEPER.PREVILEGE.TENANT", privilege.Tenant),
		validator.StringRequired("GATEKEEPER.PREVILEGE.USERNAME", privilege.Username),
		validator.StringRequired("GATEKEEPER.PREVILEGE.ROLE", privilege.Role),
	)
}
