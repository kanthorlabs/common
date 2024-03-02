package gatekeeper

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/gatekeeper/config"
	"github.com/kanthorlabs/common/gatekeeper/entities"
	sqlx "github.com/kanthorlabs/common/persistence/sqlx/config"
	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/common/testify"
	"github.com/stretchr/testify/require"
)

func TestOpa(t *testing.T) {
	privileges, count := setup(t)
	defs := definitions(t)

	t.Run("New", func(st *testing.T) {
		st.Run("KO - configuration error", func(sst *testing.T) {
			_, err := New(&config.Config{}, testify.Logger())
			require.ErrorContains(sst, err, "GATEKEEPER.CONFIG")
		})

		st.Run("KO - sqlx error", func(sst *testing.T) {
			conf := &config.Config{
				Engine: config.EngineRBAC,
				Privilege: config.Privilege{
					Sqlx: sqlx.Config{},
				},
			}
			_, err := New(conf, testify.Logger())
			require.ErrorContains(sst, err, "SQLX.CONFIG")
		})
	})

	t.Run(".Connect/.Readiness/.Liveness/.Disconnect", func(st *testing.T) {
		gk, err := New(testconf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, gk.Connect(context.Background()))
		require.NoError(st, gk.Readiness())
		require.NoError(st, gk.Liveness())
		require.NoError(st, gk.Disconnect(context.Background()))
	})

	t.Run(".Grant", func(st *testing.T) {
		gk, err := New(testconf, testify.Logger())
		require.NoError(st, err)

		gk.Connect(context.Background())
		defer gk.Disconnect(context.Background())

		orm := gk.(*opa).orm
		tx := orm.Create(privileges)
		require.NoError(st, tx.Error)

		st.Run("OK", func(sst *testing.T) {
			i := testdata.Fake.IntBetween(0, len(privileges)-1)
			evaluation := &entities.Evaluation{
				Tenant:   uuid.NewString(),
				Username: uuid.NewString(),
				Role:     privileges[i].Role,
			}

			require.NoError(sst, gk.Grant(context.Background(), evaluation))
		})

		st.Run("KO - role not exist error", func(sst *testing.T) {
			evaluation := &entities.Evaluation{
				Tenant:   uuid.NewString(),
				Username: uuid.NewString(),
				Role:     uuid.NewString(),
			}

			err := gk.Grant(context.Background(), evaluation)
			require.ErrorContains(sst, err, "GATEKEEPER.GRANT.ROLE_NOT_EXIST")
		})

		st.Run("KO - evaluation validate error", func(sst *testing.T) {
			i := testdata.Fake.IntBetween(0, len(privileges)-1)
			evaluation := &entities.Evaluation{
				Tenant:   privileges[i].Tenant,
				Username: privileges[i].Username,
			}

			// duplicated
			require.ErrorContains(sst, gk.Grant(context.Background(), evaluation), "GATEKEEPER.EVALUATION.")
		})

		st.Run("KO - duplicated error", func(sst *testing.T) {
			i := testdata.Fake.IntBetween(0, len(privileges)-1)
			evaluation := &entities.Evaluation{
				Tenant:   privileges[i].Tenant,
				Username: privileges[i].Username,
				Role:     privileges[i].Role,
			}

			// duplicated
			require.NotNil(sst, gk.Grant(context.Background(), evaluation))
		})
	})

	t.Run(".Enforce", func(st *testing.T) {
		gk, err := New(testconf, testify.Logger())
		require.NoError(st, err)

		gk.Connect(context.Background())
		defer gk.Disconnect(context.Background())

		orm := gk.(*opa).orm
		tx := orm.Create(privileges)
		require.NoError(st, tx.Error)

		st.Run("OK", func(sst *testing.T) {
			evaluation := &entities.Evaluation{
				Tenant:   uuid.NewString(),
				Username: uuid.NewString(),
				Role:     "own",
			}
			// grant the permissions first
			require.NoError(sst, gk.Grant(context.Background(), evaluation))

			permission := &entities.Permission{
				Action: "DELETE",
				Object: "/",
			}

			// then enforce it and expect error because the permission is non-sense
			require.NotNil(sst, gk.Enforce(context.Background(), evaluation, permission))
		})

		st.Run("KO - evaluation validate error", func(sst *testing.T) {
			i := testdata.Fake.IntBetween(0, len(privileges)-1)
			evaluation := &entities.Evaluation{
				Tenant: privileges[i].Tenant,
			}

			// duplicated
			require.ErrorContains(sst, gk.Grant(context.Background(), evaluation), "GATEKEEPER.EVALUATION.")
		})

		st.Run("KO - permission validate error", func(sst *testing.T) {
			i := testdata.Fake.IntBetween(0, len(privileges)-1)
			evaluation := &entities.Evaluation{
				Tenant:   privileges[i].Tenant,
				Username: privileges[i].Username,
			}
			permission := &entities.Permission{
				Action: "DELETE",
			}

			err := gk.Enforce(context.Background(), evaluation, permission)
			require.ErrorContains(sst, err, "GATEKEEPER.PERMISSION.")
		})

		st.Run("KO - no privileges error", func(sst *testing.T) {
			evaluation := &entities.Evaluation{
				Username: uuid.NewString(),
				Tenant:   uuid.NewString(),
			}
			permission := &entities.Permission{
				Action: "DELETE",
				Object: "/",
			}

			err := gk.Enforce(context.Background(), evaluation, permission)
			require.ErrorContains(sst, err, "GATEKEEPER.ENFORCE.PRIVILEGE_EMPTY")
		})
	})

	t.Run(".Revoke", func(st *testing.T) {
		gk, err := New(testconf, testify.Logger())
		require.NoError(st, err)

		gk.Connect(context.Background())
		defer gk.Disconnect(context.Background())

		orm := gk.(*opa).orm
		tx := orm.Create(privileges)
		require.NoError(st, tx.Error)

		st.Run("OK", func(sst *testing.T) {
			i := testdata.Fake.IntBetween(0, len(privileges)-1)
			evaluation := &entities.Evaluation{
				Tenant:   privileges[i].Tenant,
				Username: privileges[i].Username,
			}

			require.NoError(sst, gk.Revoke(context.Background(), evaluation))
		})

		st.Run("KO - evaluation error", func(sst *testing.T) {
			i := testdata.Fake.IntBetween(0, len(privileges)-1)
			evaluation := &entities.Evaluation{
				Username: privileges[i].Username,
			}

			err := gk.Revoke(context.Background(), evaluation)
			require.ErrorContains(sst, err, "GATEKEEPER.EVALUATION.")
		})

		st.Run("KO - revoke not exist privilege", func(sst *testing.T) {
			evaluation := &entities.Evaluation{
				Username: testdata.Fake.Internet().Email(),
				Tenant:   uuid.NewString(),
				Role:     testdata.Fake.Color().SafeColorName(),
			}

			err := gk.Revoke(context.Background(), evaluation)
			require.ErrorContains(sst, err, "GATEKEEPER.REVOKE.PRIVILEGE_NOT_EXIST")
		})
	})

	t.Run(".Users", func(st *testing.T) {
		gk, err := New(testconf, testify.Logger())
		require.NoError(st, err)

		gk.Connect(context.Background())
		defer gk.Disconnect(context.Background())

		orm := gk.(*opa).orm
		tx := orm.Create(privileges)
		require.NoError(st, tx.Error)

		st.Run("OK", func(sst *testing.T) {
			i := testdata.Fake.IntBetween(0, len(privileges)-1)
			tenant := privileges[i].Tenant

			users, err := gk.Users(context.Background(), tenant)
			require.NoError(sst, err)

			require.Equal(sst, count, len(users))

			for i := range users {
				require.Equal(sst, len(defs), len(users[i].Roles))
			}
		})
	})

	t.Run(".Tenants", func(st *testing.T) {
		gk, err := New(testconf, testify.Logger())
		require.NoError(st, err)

		gk.Connect(context.Background())
		defer gk.Disconnect(context.Background())

		orm := gk.(*opa).orm
		tx := orm.Create(privileges)
		require.NoError(st, tx.Error)

		st.Run("OK", func(sst *testing.T) {
			i := testdata.Fake.IntBetween(0, len(privileges)-1)
			username := privileges[i].Username

			users, err := gk.Tenants(context.Background(), username)
			require.NoError(sst, err)

			require.Equal(sst, 1, len(users))

			for i := range users {
				require.Equal(sst, len(defs), len(users[i].Roles))
			}
		})
	})
}
