package rego

import (
	"context"

	"github.com/kanthorlabs/common/gatekeeper/entities"
)

// Evaluate is a function that evaluates a permission against a set of privileges
type Evaluate func(ectx context.Context, permission *entities.Permission, privileges []entities.Privilege) error
