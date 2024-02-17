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
		conf := &config.Config{
			Engine: config.EngineRBAC,
			Privilege: config.Privilege{
				Sqlx: sqlx.Config{
					Uri: testdata.SqliteUri,
					Connection: sqlx.Connection{
						MaxLifetime:  sqlx.DefaultConnMaxLifetime,
						MaxIdletime:  sqlx.DefaultConnMaxIdletime,
						MaxIdleCount: sqlx.DefaultConnMaxIdleCount,
						MaxOpenCount: sqlx.DefaultConnMaxOpenCount,
					},
				},
			},
		}
		gk, err := New(conf, testify.Logger())
		require.Nil(st, err)

		require.Nil(st, gk.Connect(context.Background()))
		require.Nil(st, gk.Readiness())
		require.Nil(st, gk.Liveness())
		require.Nil(st, gk.Disconnect(context.Background()))
	})

	t.Run(".Grant", func(st *testing.T) {
		conf := &config.Config{
			Engine: config.EngineRBAC,
			Privilege: config.Privilege{
				Sqlx: sqlx.Config{
					Uri: testdata.SqliteUri,
					Connection: sqlx.Connection{
						MaxLifetime:  sqlx.DefaultConnMaxLifetime,
						MaxIdletime:  sqlx.DefaultConnMaxIdletime,
						MaxIdleCount: sqlx.DefaultConnMaxIdleCount,
						MaxOpenCount: sqlx.DefaultConnMaxOpenCount,
					},
				},
			},
		}
		gk, err := New(conf, testify.Logger())
		require.Nil(st, err)

		gk.Connect(context.Background())
		defer gk.Disconnect(context.Background())

		orm := gk.(*opa).orm
		tx := orm.Create(privileges)
		require.Nil(st, tx.Error)

		st.Run("OK", func(sst *testing.T) {
			evaluation := &entities.Evaluation{
				Username: testdata.Fake.Internet().Email(),
				Tenant:   uuid.NewString(),
				Role:     testdata.Fake.Color().SafeColorName(),
			}

			require.Nil(sst, gk.Grant(context.Background(), evaluation))
		})

		st.Run("KO - duplicated", func(sst *testing.T) {
			i := testdata.Fake.IntBetween(0, len(privileges)-1)
			evaluation := &entities.Evaluation{
				Username: privileges[i].Username,
				Tenant:   privileges[i].Tenant,
				Role:     privileges[i].Role,
			}

			// duplicated
			require.NotNil(sst, gk.Grant(context.Background(), evaluation))
		})
	})

	t.Run(".Revoke", func(st *testing.T) {
		conf := &config.Config{
			Engine: config.EngineRBAC,
			Privilege: config.Privilege{
				Sqlx: sqlx.Config{
					Uri: testdata.SqliteUri,
					Connection: sqlx.Connection{
						MaxLifetime:  sqlx.DefaultConnMaxLifetime,
						MaxIdletime:  sqlx.DefaultConnMaxIdletime,
						MaxIdleCount: sqlx.DefaultConnMaxIdleCount,
						MaxOpenCount: sqlx.DefaultConnMaxOpenCount,
					},
				},
			},
		}
		gk, err := New(conf, testify.Logger())
		require.Nil(st, err)

		gk.Connect(context.Background())
		defer gk.Disconnect(context.Background())

		orm := gk.(*opa).orm
		tx := orm.Create(privileges)
		require.Nil(st, tx.Error)

		st.Run("OK", func(sst *testing.T) {
			i := testdata.Fake.IntBetween(0, len(privileges)-1)
			evaluation := &entities.Evaluation{
				Username: privileges[i].Username,
				Tenant:   privileges[i].Tenant,
				Role:     privileges[i].Role,
			}

			require.Nil(sst, gk.Revoke(context.Background(), evaluation))
		})

		st.Run("OK - revoke not exist privilege", func(sst *testing.T) {
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
		conf := &config.Config{
			Engine: config.EngineRBAC,
			Privilege: config.Privilege{
				Sqlx: sqlx.Config{
					Uri: testdata.SqliteUri,
					Connection: sqlx.Connection{
						MaxLifetime:  sqlx.DefaultConnMaxLifetime,
						MaxIdletime:  sqlx.DefaultConnMaxIdletime,
						MaxIdleCount: sqlx.DefaultConnMaxIdleCount,
						MaxOpenCount: sqlx.DefaultConnMaxOpenCount,
					},
				},
			},
		}
		gk, err := New(conf, testify.Logger())
		require.Nil(st, err)

		gk.Connect(context.Background())
		defer gk.Disconnect(context.Background())

		orm := gk.(*opa).orm
		tx := orm.Create(privileges)
		require.Nil(st, tx.Error)

		st.Run("OK", func(sst *testing.T) {
			i := testdata.Fake.IntBetween(0, len(privileges)-1)
			tenant := privileges[i].Tenant

			users, err := gk.Users(context.Background(), tenant)
			require.Nil(sst, err)

			require.Equal(sst, count, len(users))

			for i := range users {
				require.Equal(sst, count, len(users[i].Roles))
			}
		})
	})

	t.Run(".Tenants", func(st *testing.T) {
		conf := &config.Config{
			Engine: config.EngineRBAC,
			Privilege: config.Privilege{
				Sqlx: sqlx.Config{
					Uri: testdata.SqliteUri,
					Connection: sqlx.Connection{
						MaxLifetime:  sqlx.DefaultConnMaxLifetime,
						MaxIdletime:  sqlx.DefaultConnMaxIdletime,
						MaxIdleCount: sqlx.DefaultConnMaxIdleCount,
						MaxOpenCount: sqlx.DefaultConnMaxOpenCount,
					},
				},
			},
		}
		gk, err := New(conf, testify.Logger())
		require.Nil(st, err)

		gk.Connect(context.Background())
		defer gk.Disconnect(context.Background())

		orm := gk.(*opa).orm
		tx := orm.Create(privileges)
		require.Nil(st, tx.Error)

		st.Run("OK", func(sst *testing.T) {
			i := testdata.Fake.IntBetween(0, len(privileges)-1)
			username := privileges[i].Username

			users, err := gk.Tenants(context.Background(), username)
			require.Nil(sst, err)

			require.Equal(sst, 1, len(users))

			for i := range users {
				require.Equal(sst, count, len(users[i].Roles))
			}
		})
	})

}
