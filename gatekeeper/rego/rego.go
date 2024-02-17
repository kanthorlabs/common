package rego

import "github.com/kanthorlabs/common/gatekeeper/entities"

type Evaluate func(permission *entities.Permission, privileges []entities.Privilege) error
