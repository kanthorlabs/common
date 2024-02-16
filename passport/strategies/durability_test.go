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
		require.Nil(st, err)

		require.Nil(st, strategy.Connect(context.Background()))
		require.Nil(st, strategy.Readiness())
		require.Nil(st, strategy.Liveness())
		require.Nil(st, strategy.Logout(context.Background(), nil))
		require.Nil(st, strategy.Disconnect(context.Background()))
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
		require.Nil(st, err)

		strategy.Connect(context.Background())
		defer strategy.Disconnect(context.Background())

		orm := strategy.(*durability).orm
		tx := orm.Create(accounts)
		require.Nil(st, tx.Error)

		st.Run("OK", func(sst *testing.T) {
			i := testdata.Fake.IntBetween(0, len(passwords)-1)
			credentials := &entities.Credentials{
				Username: accounts[i].Username,
				Password: passwords[i],
			}
			acc, err := strategy.Login(context.Background(), credentials)
			require.Nil(sst, err)
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
		require.Nil(st, err)

		strategy.Connect(context.Background())
		defer strategy.Disconnect(context.Background())

		orm := strategy.(*durability).orm
		tx := orm.Create(accounts)
		require.Nil(st, tx.Error)

		st.Run("OK", func(sst *testing.T) {
			i := testdata.Fake.IntBetween(0, len(passwords)-1)
			credentials := &entities.Credentials{
				Username: accounts[i].Username,
				Password: passwords[i],
			}
			acc, err := strategy.Verify(context.Background(), credentials)
			require.Nil(sst, err)
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
		require.Nil(st, err)

		strategy.Connect(context.Background())
		defer strategy.Disconnect(context.Background())

		orm := strategy.(*durability).orm
		tx := orm.Create(accounts)
		require.Nil(st, tx.Error)

		t.Run("KO", func(sst *testing.T) {
			pass := uuid.NewString()
			hash, err := password.HashString(pass)
			require.Nil(st, err)

			acc := &entities.Account{
				Username:     uuid.NewString(),
				PasswordHash: hash,
				Name:         testdata.Fake.Internet().User(),
				CreatedAt:    time.Now().UnixMilli(),
				UpdatedAt:    time.Now().UnixMilli(),
			}

			require.Nil(st, strategy.Register(context.Background(), acc))
		})

		t.Run("KO", func(sst *testing.T) {
			i := testdata.Fake.IntBetween(0, len(accounts)-1)
			err := strategy.Register(context.Background(), &accounts[i])
			require.ErrorIs(sst, err, ErrRegister)
		})
	})
}
