package rego

import (
	"context"

	"github.com/kanthorlabs/common/gatekeeper/entities"
)

type Evaluate func(ectx context.Context, permission *entities.Permission, privileges []entities.Privilege) error
