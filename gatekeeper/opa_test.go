package gatekeeper

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/gatekeeper/config"
	"github.com/kanthorlabs/common/gatekeeper/entities"
	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/common/testify"
	"github.com/stretchr/testify/require"
)

func TestOpa_New(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		_, err := NewOpa(testconf, testify.Logger())
		require.NoError(st, err)
	})

	t.Run("KO - configuration error", func(st *testing.T) {
		conf := &config.Config{}
		_, err := NewOpa(conf, testify.Logger())
		require.ErrorContains(st, err, "GATEKEEPER.CONFIG")
	})
}

func TestOpa_Connect(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		gk, err := NewOpa(testconf, testify.Logger())
		require.NoError(t, err)

		require.NoError(st, gk.Connect(context.Background()))
	})

	t.Run("KO - already connected error", func(st *testing.T) {
		gk, err := NewOpa(testconf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, gk.Connect(context.Background()))
		require.ErrorIs(st, gk.Connect(context.Background()), ErrAlreadyConnected)
	})

	t.Run("KO - sqlx error", func(st *testing.T) {
		conf := &config.Config{
			Engine:      testconf.Engine,
			Definitions: testconf.Definitions,
			Privilege:   testconf.Privilege,
		}
		conf.Privilege.Sqlx.Uri = testdata.PostgresUri

		gk, err := NewOpa(conf, testify.Logger())
		require.NoError(st, err)

		require.ErrorContains(st, gk.Connect(context.Background()), "SQLX.CONNECT")
	})

	t.Run("KO - parse definitions error", func(st *testing.T) {
		conf := &config.Config{
			Engine:      testconf.Engine,
			Definitions: testconf.Definitions,
			Privilege:   testconf.Privilege,
		}
		conf.Definitions.Uri = "base64://-"

		gk, err := NewOpa(conf, testify.Logger())
		require.NoError(st, err)

		require.ErrorContains(st, gk.Connect(context.Background()), "GATEKEEPER.CONFIG.DEFINITIONS.BASE64")
	})

	t.Run("KO - init rego RBAC error", func(st *testing.T) {
		conf := &config.Config{
			Engine:      testconf.Engine,
			Definitions: testconf.Definitions,
			Privilege:   testconf.Privilege,
		}
		conf.Definitions.Uri = "base64://eyJhZG1pbmlzdHJhdG9yIjogW3siYWN0aW9uIjogIioifV19Cg=="

		gk, err := NewOpa(conf, testify.Logger())
		require.NoError(st, err)

		require.ErrorContains(st, gk.Connect(context.Background()), "GATEKEEPER.REGO.RBAC")
	})
}

func TestOpa_Readiness(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		gk, err := NewOpa(testconf, testify.Logger())
		require.NoError(t, err)

		require.NoError(st, gk.Connect(context.Background()))
		require.NoError(st, gk.Readiness())
	})

	t.Run("OK - disconnected", func(st *testing.T) {
		gk, err := NewOpa(testconf, testify.Logger())
		require.NoError(t, err)

		require.NoError(st, gk.Connect(context.Background()))
		require.NoError(st, gk.Disconnect(context.Background()))
		require.NoError(st, gk.Readiness())
	})
}

func TestOpa_Liveness(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		gk, err := NewOpa(testconf, testify.Logger())
		require.NoError(t, err)

		require.NoError(st, gk.Connect(context.Background()))
		require.NoError(st, gk.Liveness())
	})

	t.Run("OK - disconnected", func(st *testing.T) {
		gk, err := NewOpa(testconf, testify.Logger())
		require.NoError(t, err)

		require.NoError(st, gk.Connect(context.Background()))
		require.NoError(st, gk.Disconnect(context.Background()))
		require.NoError(st, gk.Liveness())
	})
}

func TestOpa_Disconnect(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		gk, err := NewOpa(testconf, testify.Logger())
		require.NoError(t, err)

		require.NoError(st, gk.Connect(context.Background()))
		require.NoError(st, gk.Disconnect(context.Background()))
	})

	t.Run("KO - not connected error", func(st *testing.T) {
		gk, err := NewOpa(testconf, testify.Logger())
		require.NoError(t, err)

		require.ErrorIs(st, gk.Disconnect(context.Background()), ErrNotConnected)
	})
}

func TestOpa_Grant(t *testing.T) {
	privileges, _ := setup(t)

	gk, err := NewOpa(testconf, testify.Logger())
	require.NoError(t, err)

	gk.Connect(context.Background())
	defer gk.Disconnect(context.Background())

	orm := gk.(*opa).orm
	tx := orm.Create(privileges)
	require.NoError(t, tx.Error)

	t.Run("OK", func(st *testing.T) {
		i := testdata.Fake.IntBetween(0, len(privileges)-1)
		evaluation := &entities.Evaluation{
			Tenant:   uuid.NewString(),
			Username: uuid.NewString(),
			Role:     privileges[i].Role,
		}

		require.NoError(st, gk.Grant(context.Background(), evaluation))
	})

	t.Run("KO - role not exist error", func(st *testing.T) {
		evaluation := &entities.Evaluation{
			Tenant:   uuid.NewString(),
			Username: uuid.NewString(),
			Role:     uuid.NewString(),
		}

		err := gk.Grant(context.Background(), evaluation)
		require.ErrorContains(st, err, "GATEKEEPER.GRANT.ROLE_NOT_EXIST")
	})

	t.Run("KO - evaluation validate error", func(st *testing.T) {
		i := testdata.Fake.IntBetween(0, len(privileges)-1)
		evaluation := &entities.Evaluation{
			Tenant:   privileges[i].Tenant,
			Username: privileges[i].Username,
		}

		// duplicated
		require.ErrorContains(st, gk.Grant(context.Background(), evaluation), "GATEKEEPER.EVALUATION.")
	})

	t.Run("KO - duplicated error", func(st *testing.T) {
		i := testdata.Fake.IntBetween(0, len(privileges)-1)
		evaluation := &entities.Evaluation{
			Tenant:   privileges[i].Tenant,
			Username: privileges[i].Username,
			Role:     privileges[i].Role,
		}

		// duplicated
		require.NotNil(st, gk.Grant(context.Background(), evaluation))
	})
}

func TestOpa_Enforce(t *testing.T) {
	privileges, _ := setup(t)

	gk, err := NewOpa(testconf, testify.Logger())
	require.NoError(t, err)

	gk.Connect(context.Background())
	defer gk.Disconnect(context.Background())

	orm := gk.(*opa).orm
	tx := orm.Create(privileges)
	require.NoError(t, tx.Error)

	t.Run("OK", func(st *testing.T) {
		evaluation := &entities.Evaluation{
			Tenant:   uuid.NewString(),
			Username: uuid.NewString(),
			Role:     "own",
		}
		// grant the permissions first
		require.NoError(st, gk.Grant(context.Background(), evaluation))

		permission := &entities.Permission{
			Scope:  entities.AnyScope,
			Action: "DELETE",
			Object: "/",
		}

		// then enforce it and expect error because the permission is non-sense
		require.NotNil(st, gk.Enforce(context.Background(), evaluation, permission))
	})

	t.Run("KO - permission validate error", func(st *testing.T) {
		i := testdata.Fake.IntBetween(0, len(privileges)-1)
		evaluation := &entities.Evaluation{
			Tenant:   privileges[i].Tenant,
			Username: privileges[i].Username,
		}
		permission := &entities.Permission{
			Action: "DELETE",
			Object: "/",
		}

		err := gk.Enforce(context.Background(), evaluation, permission)
		require.ErrorContains(st, err, "GATEKEEPER.PERMISSION.")
	})

	t.Run("KO - no privileges error", func(st *testing.T) {
		evaluation := &entities.Evaluation{
			Username: uuid.NewString(),
			Tenant:   uuid.NewString(),
		}
		permission := &entities.Permission{
			Scope:  entities.AnyScope,
			Action: "DELETE",
			Object: "/",
		}

		err := gk.Enforce(context.Background(), evaluation, permission)
		require.ErrorContains(st, err, "GATEKEEPER.ENFORCE.PRIVILEGE_EMPTY")
	})
}

func TestOpa_Revoke(t *testing.T) {
	privileges, _ := setup(t)

	gk, err := NewOpa(testconf, testify.Logger())
	require.NoError(t, err)

	gk.Connect(context.Background())
	defer gk.Disconnect(context.Background())

	orm := gk.(*opa).orm
	tx := orm.Create(privileges)
	require.NoError(t, tx.Error)

	t.Run("OK", func(st *testing.T) {
		i := testdata.Fake.IntBetween(0, len(privileges)-1)
		evaluation := &entities.Evaluation{
			Tenant:   privileges[i].Tenant,
			Username: privileges[i].Username,
		}

		require.NoError(st, gk.Revoke(context.Background(), evaluation))
	})

	t.Run("KO - evaluation error", func(st *testing.T) {
		i := testdata.Fake.IntBetween(0, len(privileges)-1)
		evaluation := &entities.Evaluation{
			Username: privileges[i].Username,
		}

		err := gk.Revoke(context.Background(), evaluation)
		require.ErrorContains(st, err, "GATEKEEPER.EVALUATION.")
	})

	t.Run("KO - revoke not exist privilege", func(st *testing.T) {
		evaluation := &entities.Evaluation{
			Username: testdata.Fake.Internet().Email(),
			Tenant:   uuid.NewString(),
			Role:     testdata.Fake.Color().SafeColorName(),
		}

		err := gk.Revoke(context.Background(), evaluation)
		require.ErrorContains(st, err, "GATEKEEPER.REVOKE.PRIVILEGE_NOT_EXIST")
	})
}

func TestOpa_Users(t *testing.T) {
	privileges, count := setup(t)
	defs := definitions(t)

	gk, err := NewOpa(testconf, testify.Logger())
	require.NoError(t, err)

	gk.Connect(context.Background())
	defer gk.Disconnect(context.Background())

	orm := gk.(*opa).orm
	tx := orm.Create(privileges)
	require.NoError(t, tx.Error)

	t.Run("OK", func(st *testing.T) {
		i := testdata.Fake.IntBetween(0, len(privileges)-1)
		tenant := privileges[i].Tenant

		users, err := gk.Users(context.Background(), tenant)
		require.NoError(st, err)

		require.Equal(st, count, len(users))

		for i := range users {
			require.Equal(st, len(defs), len(users[i].Roles))
		}
	})
}

func TestOpa_Tenants(t *testing.T) {
	privileges, _ := setup(t)
	defs := definitions(t)
	gk, err := NewOpa(testconf, testify.Logger())
	require.NoError(t, err)

	gk.Connect(context.Background())
	defer gk.Disconnect(context.Background())

	orm := gk.(*opa).orm
	tx := orm.Create(privileges)
	require.NoError(t, tx.Error)

	t.Run("OK", func(st *testing.T) {
		i := testdata.Fake.IntBetween(0, len(privileges)-1)
		username := privileges[i].Username

		users, err := gk.Tenants(context.Background(), username)
		require.NoError(st, err)

		require.Equal(st, 1, len(users))

		for i := range users {
			require.Equal(st, len(defs), len(users[i].Roles))
		}
	})
}
