package strategies

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/cipher/password"
	"github.com/kanthorlabs/common/passport/config"
	"github.com/kanthorlabs/common/passport/entities"
	sqlx "github.com/kanthorlabs/common/persistence/sqlx/config"
	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/common/testify"
	"github.com/stretchr/testify/require"
)

func TestDurability(t *testing.T) {
	accounts, passwords := setup(t)

	t.Run("New", func(st *testing.T) {
		st.Run("KO - configuration error", func(sst *testing.T) {
			_, err := NewDurability(&config.Durability{}, testify.Logger())
			require.ErrorContains(sst, err, "SQLX.CONFIG")
		})

		st.Run("KO - sqlx error", func(sst *testing.T) {
			conf := &config.Durability{Sqlx: sqlx.Config{}}
			_, err := NewDurability(conf, testify.Logger())
			require.ErrorContains(sst, err, "SQLX.CONFIG")
		})
	})

	t.Run(".Connect/.Readiness/.Liveness/.Logout/.Disconnect", func(st *testing.T) {
		conf := &config.Durability{Sqlx: sqlx.Config{
			Uri: testdata.SqliteUri,
			Connection: sqlx.Connection{
				MaxLifetime:  sqlx.DefaultConnMaxLifetime,
				MaxIdletime:  sqlx.DefaultConnMaxIdletime,
				MaxIdleCount: sqlx.DefaultConnMaxIdleCount,
				MaxOpenCount: sqlx.DefaultConnMaxOpenCount,
			},
		}}
		strategy, err := NewDurability(conf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, strategy.Connect(context.Background()))
		require.NoError(st, strategy.Readiness())
		require.NoError(st, strategy.Liveness())
		require.NoError(st, strategy.Logout(context.Background(), nil))
		require.NoError(st, strategy.Disconnect(context.Background()))
	})

	t.Run(".Login", func(st *testing.T) {
		conf := &config.Durability{Sqlx: sqlx.Config{
			Uri: testdata.SqliteUri,
			Connection: sqlx.Connection{
				MaxLifetime:  sqlx.DefaultConnMaxLifetime,
				MaxIdletime:  sqlx.DefaultConnMaxIdletime,
				MaxIdleCount: sqlx.DefaultConnMaxIdleCount,
				MaxOpenCount: sqlx.DefaultConnMaxOpenCount,
			},
		}}
		strategy, err := NewDurability(conf, testify.Logger())
		require.NoError(st, err)

		strategy.Connect(context.Background())
		defer strategy.Disconnect(context.Background())

		orm := strategy.(*durability).orm
		tx := orm.Create(accounts)
		require.NoError(st, tx.Error)

		st.Run("OK", func(sst *testing.T) {
			i := testdata.Fake.IntBetween(0, len(passwords)-1)
			credentials := &entities.Credentials{
				Username: accounts[i].Username,
				Password: passwords[i],
			}
			acc, err := strategy.Login(context.Background(), credentials)
			require.NoError(sst, err)
			require.Equal(sst, credentials.Username, acc.Username)
			require.Empty(sst, acc.PasswordHash)
		})

		st.Run("KO - credentials error", func(sst *testing.T) {
			_, err := strategy.Login(context.Background(), nil)
			require.ErrorContains(sst, err, "PASSPORT.CREDENTIALS")

			_, err = strategy.Login(context.Background(), &entities.Credentials{})
			require.ErrorContains(sst, err, "PASSPORT.CREDENTIALS")
		})

		st.Run("KO - user not found", func(sst *testing.T) {
			credentials := &entities.Credentials{
				Username: uuid.NewString(),
				Password: testdata.Fake.Internet().Password(),
			}
			_, err := strategy.Login(context.Background(), credentials)
			require.ErrorIs(sst, err, ErrLogin)
		})

		st.Run("KO - password not match", func(sst *testing.T) {
			i := testdata.Fake.IntBetween(0, len(passwords)-1)
			credentials := &entities.Credentials{
				Username: accounts[i].Username,
				Password: testdata.Fake.Internet().Password(),
			}
			_, err := strategy.Login(context.Background(), credentials)
			require.ErrorIs(sst, err, ErrLogin)
		})
	})

	t.Run(".Verify", func(st *testing.T) {
		conf := &config.Durability{Sqlx: sqlx.Config{
			Uri: testdata.SqliteUri,
			Connection: sqlx.Connection{
				MaxLifetime:  sqlx.DefaultConnMaxLifetime,
				MaxIdletime:  sqlx.DefaultConnMaxIdletime,
				MaxIdleCount: sqlx.DefaultConnMaxIdleCount,
				MaxOpenCount: sqlx.DefaultConnMaxOpenCount,
			},
		}}
		strategy, err := NewDurability(conf, testify.Logger())
		require.NoError(st, err)

		strategy.Connect(context.Background())
		defer strategy.Disconnect(context.Background())

		orm := strategy.(*durability).orm
		tx := orm.Create(accounts)
		require.NoError(st, tx.Error)

		st.Run("OK", func(sst *testing.T) {
			i := testdata.Fake.IntBetween(0, len(passwords)-1)
			credentials := &entities.Credentials{
				Username: accounts[i].Username,
				Password: passwords[i],
			}
			acc, err := strategy.Verify(context.Background(), credentials)
			require.NoError(sst, err)
			require.Equal(sst, credentials.Username, acc.Username)
			require.Empty(sst, acc.PasswordHash)
		})
	})

	t.Run(".Register", func(st *testing.T) {
		conf := &config.Durability{Sqlx: sqlx.Config{
			Uri: testdata.SqliteUri,
			Connection: sqlx.Connection{
				MaxLifetime:  sqlx.DefaultConnMaxLifetime,
				MaxIdletime:  sqlx.DefaultConnMaxIdletime,
				MaxIdleCount: sqlx.DefaultConnMaxIdleCount,
				MaxOpenCount: sqlx.DefaultConnMaxOpenCount,
			},
		}}
		strategy, err := NewDurability(conf, testify.Logger())
		require.NoError(st, err)

		strategy.Connect(context.Background())
		defer strategy.Disconnect(context.Background())

		orm := strategy.(*durability).orm
		tx := orm.Create(accounts)
		require.NoError(st, tx.Error)

		t.Run("KO", func(sst *testing.T) {
			pass := uuid.NewString()
			hash, err := password.HashString(pass)
			require.NoError(st, err)

			acc := &entities.Account{
				Username:     uuid.NewString(),
				PasswordHash: hash,
				Name:         testdata.Fake.Internet().User(),
				CreatedAt:    time.Now().UnixMilli(),
				UpdatedAt:    time.Now().UnixMilli(),
			}

			require.NoError(st, strategy.Register(context.Background(), acc))
		})

		t.Run("KO", func(sst *testing.T) {
			i := testdata.Fake.IntBetween(0, len(accounts)-1)
			err := strategy.Register(context.Background(), &accounts[i])
			require.ErrorIs(sst, err, ErrRegister)
		})
	})

	t.Run(".Deactivate", func(st *testing.T) {
		conf := &config.Durability{Sqlx: sqlx.Config{
			Uri: testdata.SqliteUri,
			Connection: sqlx.Connection{
				MaxLifetime:  sqlx.DefaultConnMaxLifetime,
				MaxIdletime:  sqlx.DefaultConnMaxIdletime,
				MaxIdleCount: sqlx.DefaultConnMaxIdleCount,
				MaxOpenCount: sqlx.DefaultConnMaxOpenCount,
			},
		}}
		strategy, err := NewDurability(conf, testify.Logger())
		require.NoError(st, err)

		strategy.Connect(context.Background())
		defer strategy.Disconnect(context.Background())

		orm := strategy.(*durability).orm
		tx := orm.Create(accounts)
		require.NoError(st, tx.Error)

		st.Run("OK", func(sst *testing.T) {
			i := testdata.Fake.IntBetween(0, len(accounts)-1)
			username := accounts[i].Username
			ts := time.Now().UnixMilli()

			err := strategy.Deactivate(context.Background(), username, ts)
			require.NoError(sst, err)
		})

		st.Run("KO - user not found", func(sst *testing.T) {
			username := uuid.NewString()
			ts := time.Now().UnixMilli()

			err := strategy.Deactivate(context.Background(), username, ts)
			require.ErrorIs(sst, err, ErrAccountNotFound)
		})
	})
}
