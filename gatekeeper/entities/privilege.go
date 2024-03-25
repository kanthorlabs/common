package entities

import (
	"github.com/kanthorlabs/common/project"
	"github.com/kanthorlabs/common/safe"
	"github.com/kanthorlabs/common/validator"
)

type Privilege struct {
	Tenant   string         `json:"tenant" yaml:"tenant" gorm:"primaryKey"`
	Username string         `json:"username" yaml:"username" gorm:"primaryKey;index:idx_username"`
	Role     string         `json:"role" yaml:"role" gorm:"primaryKey"`
	Metadata *safe.Metadata `json:"metadata" yaml:"metadata"`

	CreatedAt int64 `json:"created_at" yaml:"created_at" gorm:"autoCreateTime:milli"`
	UpdatedAt int64 `json:"updated_at" yaml:"updated_at" gorm:"autoUpdateTime:milli"`
}

func (privilege *Privilege) TableName() string {
	return project.Name("gatekeeper_privilege")
}

func (privilege *Privilege) Validate() error {
	return validator.Validate(
		validator.StringRequired("GATEKEEPER.PREVILEGE.TENANT", privilege.Tenant),
		validator.StringRequired("GATEKEEPER.PREVILEGE.USERNAME", privilege.Username),
		validator.StringRequired("GATEKEEPER.PREVILEGE.ROLE", privilege.Role),
	)
}
