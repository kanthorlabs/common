package gatekeeper

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/gatekeeper/entities"
	"github.com/kanthorlabs/common/safe"
	"github.com/kanthorlabs/common/testdata"
)

func setup(t *testing.T) ([]entities.Privilege, int) {
	privileges := make([]entities.Privilege, 0)

	count := testdata.Fake.IntBetween(5, 10)

	for i := 0; i < count; i++ {
		tenant := uuid.NewString()
		for j := 0; j < count; j++ {
			username := testdata.Fake.Internet().Email()
			for k := 0; k < count; k++ {
				privilege := entities.Privilege{
					Username:  username,
					Tenant:    tenant,
					Role:      uuid.NewString(),
					Metadata:  &safe.Metadata{},
					CreatedAt: time.Now().UnixMilli(),
					UpdatedAt: time.Now().UnixMilli(),
				}
				privileges = append(privileges, privilege)
			}
		}
	}

	return privileges, count
}
