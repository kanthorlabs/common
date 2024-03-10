package strategies

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/cipher/password"
	"github.com/kanthorlabs/common/passport/config"
	"github.com/kanthorlabs/common/passport/entities"
	sqlxconfig "github.com/kanthorlabs/common/persistence/sqlx/config"
	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/common/testify"
	"github.com/stretchr/testify/require"
)

func TestDurability_New(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		conf := &config.Durability{Sqlx: sqlxconfig.Config{
			Uri: testdata.SqliteUri,
			Connection: sqlxconfig.Connection{
				MaxLifetime:  sqlxconfig.DefaultConnMaxLifetime,
				MaxIdletime:  sqlxconfig.DefaultConnMaxIdletime,
				MaxIdleCount: sqlxconfig.DefaultConnMaxIdleCount,
				MaxOpenCount: sqlxconfig.DefaultConnMaxOpenCount,
			},
		}}
		_, err := NewDurability(conf, testify.Logger())
		require.NoError(st, err)
	})

	t.Run("KO - configuration error", func(st *testing.T) {
		_, err := NewDurability(&config.Durability{}, testify.Logger())
		require.ErrorContains(st, err, "SQLX.CONFIG")
	})

	t.Run("KO - sqlx error", func(st *testing.T) {
		conf := &config.Durability{Sqlx: sqlxconfig.Config{}}
		_, err := NewDurability(conf, testify.Logger())
		require.ErrorContains(st, err, "SQLX.CONFIG")
	})
}

func TestDurability_ParseCredentials(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		c, err := NewDurability(durabilityconf, testify.Logger())
		require.NoError(st, err)

		_, err = c.ParseCredentials(context.Background(), "basic "+testdata.Fake.Internet().Password())
		require.ErrorIs(st, err, ErrParseCredentials)
	})

	t.Run("KO - unknown scheme error", func(st *testing.T) {
		c, err := NewDurability(durabilityconf, testify.Logger())
		require.NoError(st, err)

		_, err = c.ParseCredentials(context.Background(), "")
		require.ErrorIs(st, err, ErrCredentialsScheme)
	})
}

func TestDurability_Connect(t *testing.T) {

	t.Run("OK", func(st *testing.T) {
		c, err := NewDurability(durabilityconf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, c.Connect(context.Background()))
	})

	t.Run("KO - already connected error", func(st *testing.T) {
		c, err := NewDurability(durabilityconf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, c.Connect(context.Background()))
		require.ErrorIs(st, c.Connect(context.Background()), ErrAlreadyConnected)
	})
}

func TestDurability_Readiness(t *testing.T) {
	conf := &config.Durability{Sqlx: sqlxconfig.Config{
		Uri: testdata.SqliteUri,
		Connection: sqlxconfig.Connection{
			MaxLifetime:  sqlxconfig.DefaultConnMaxLifetime,
			MaxIdletime:  sqlxconfig.DefaultConnMaxIdletime,
			MaxIdleCount: sqlxconfig.DefaultConnMaxIdleCount,
			MaxOpenCount: sqlxconfig.DefaultConnMaxOpenCount,
		},
	}}

	t.Run("OK", func(st *testing.T) {
		c, err := NewDurability(conf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, c.Connect(context.Background()))
		require.NoError(st, c.Readiness())
	})

	t.Run("OK - disconnected", func(st *testing.T) {
		c, err := NewDurability(conf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, c.Connect(context.Background()))
		require.NoError(st, c.Disconnect(context.Background()))
		require.NoError(st, c.Readiness())
	})

	t.Run("KO - not connected error", func(st *testing.T) {
		c, err := NewDurability(conf, testify.Logger())
		require.NoError(st, err)

		require.ErrorIs(st, c.Readiness(), ErrNotConnected)
	})
}

func TestDurability_Liveness(t *testing.T) {
	conf := &config.Durability{Sqlx: sqlxconfig.Config{
		Uri: testdata.SqliteUri,
		Connection: sqlxconfig.Connection{
			MaxLifetime:  sqlxconfig.DefaultConnMaxLifetime,
			MaxIdletime:  sqlxconfig.DefaultConnMaxIdletime,
			MaxIdleCount: sqlxconfig.DefaultConnMaxIdleCount,
			MaxOpenCount: sqlxconfig.DefaultConnMaxOpenCount,
		},
	}}

	t.Run("OK", func(st *testing.T) {
		c, err := NewDurability(conf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, c.Connect(context.Background()))
		require.NoError(st, c.Liveness())
	})

	t.Run("OK - disconnected", func(st *testing.T) {
		c, err := NewDurability(conf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, c.Connect(context.Background()))
		require.NoError(st, c.Disconnect(context.Background()))
		require.NoError(st, c.Liveness())
	})

	t.Run("KO - not connected error", func(st *testing.T) {
		c, err := NewDurability(conf, testify.Logger())
		require.NoError(st, err)

		require.ErrorIs(st, c.Liveness(), ErrNotConnected)
	})
}

func TestDurability_Disconnect(t *testing.T) {
	conf := &config.Durability{Sqlx: sqlxconfig.Config{
		Uri: testdata.SqliteUri,
		Connection: sqlxconfig.Connection{
			MaxLifetime:  sqlxconfig.DefaultConnMaxLifetime,
			MaxIdletime:  sqlxconfig.DefaultConnMaxIdletime,
			MaxIdleCount: sqlxconfig.DefaultConnMaxIdleCount,
			MaxOpenCount: sqlxconfig.DefaultConnMaxOpenCount,
		},
	}}

	t.Run("OK", func(st *testing.T) {
		c, err := NewDurability(conf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, c.Connect(context.Background()))
		require.NoError(st, c.Disconnect(context.Background()))
	})

	t.Run("KO - not connected error", func(st *testing.T) {
		c, err := NewDurability(conf, testify.Logger())
		require.NoError(st, err)

		require.ErrorIs(st, c.Disconnect(context.Background()), ErrNotConnected)
	})
}

func TestDurability_Login(t *testing.T) {
	accounts, passwords := setup(t)

	conf := &config.Durability{Sqlx: sqlxconfig.Config{
		Uri: testdata.SqliteUri,
		Connection: sqlxconfig.Connection{
			MaxLifetime:  sqlxconfig.DefaultConnMaxLifetime,
			MaxIdletime:  sqlxconfig.DefaultConnMaxIdletime,
			MaxIdleCount: sqlxconfig.DefaultConnMaxIdleCount,
			MaxOpenCount: sqlxconfig.DefaultConnMaxOpenCount,
		},
	}}
	strategy, err := NewDurability(conf, testify.Logger())
	require.NoError(t, err)

	strategy.Connect(context.Background())
	defer strategy.Disconnect(context.Background())

	orm := strategy.(*durability).orm
	require.NoError(t, orm.Create(accounts).Error)

	t.Run("OK", func(st *testing.T) {
		i := testdata.Fake.IntBetween(0, len(passwords)-1)
		credentials := &entities.Credentials{
			Username: accounts[i].Username,
			Password: passwords[i],
		}
		acc, err := strategy.Login(context.Background(), credentials)
		require.NoError(st, err)
		require.Equal(st, credentials.Username, acc.Username)
		require.Empty(st, acc.PasswordHash)
	})

	t.Run("KO - credentials error", func(st *testing.T) {
		_, err := strategy.Login(context.Background(), nil)
		require.ErrorContains(st, err, "PASSPORT.CREDENTIALS")

		_, err = strategy.Login(context.Background(), &entities.Credentials{})
		require.ErrorContains(st, err, "PASSPORT.CREDENTIALS")
	})

	t.Run("KO - user not found", func(st *testing.T) {
		credentials := &entities.Credentials{
			Username: uuid.NewString(),
			Password: testdata.Fake.Internet().Password(),
		}
		_, err := strategy.Login(context.Background(), credentials)
		require.ErrorIs(st, err, ErrLogin)
	})

	t.Run("KO - password not match", func(st *testing.T) {
		i := testdata.Fake.IntBetween(0, len(passwords)-1)
		credentials := &entities.Credentials{
			Username: accounts[i].Username,
			Password: testdata.Fake.Internet().Password(),
		}
		_, err := strategy.Login(context.Background(), credentials)
		require.ErrorIs(st, err, ErrLogin)
	})
}

func TestDurability_Logout(t *testing.T) {
	accounts, _ := setup(t)

	strategy, err := NewDurability(durabilityconf, testify.Logger())
	require.NoError(t, err)

	strategy.Connect(context.Background())
	defer strategy.Disconnect(context.Background())

	orm := strategy.(*durability).orm
	tx := orm.Create(accounts)
	require.NoError(t, tx.Error)

	t.Run("OK", func(st *testing.T) {
		i := testdata.Fake.IntBetween(0, len(accounts)-1)
		username := accounts[i].Username
		credentials := &entities.Credentials{
			Username: username,
		}

		err := strategy.Logout(context.Background(), credentials)
		require.NoError(st, err)
	})
}

func TestDurability_Verify(t *testing.T) {
	accounts, passwords := setup(t)

	conf := &config.Durability{Sqlx: sqlxconfig.Config{
		Uri: testdata.SqliteUri,
		Connection: sqlxconfig.Connection{
			MaxLifetime:  sqlxconfig.DefaultConnMaxLifetime,
			MaxIdletime:  sqlxconfig.DefaultConnMaxIdletime,
			MaxIdleCount: sqlxconfig.DefaultConnMaxIdleCount,
			MaxOpenCount: sqlxconfig.DefaultConnMaxOpenCount,
		},
	}}
	strategy, err := NewDurability(conf, testify.Logger())
	require.NoError(t, err)

	strategy.Connect(context.Background())
	defer strategy.Disconnect(context.Background())

	orm := strategy.(*durability).orm
	tx := orm.Create(accounts)
	require.NoError(t, tx.Error)

	t.Run("OK", func(st *testing.T) {
		i := testdata.Fake.IntBetween(0, len(passwords)-1)
		credentials := &entities.Credentials{
			Username: accounts[i].Username,
			Password: passwords[i],
		}
		acc, err := strategy.Verify(context.Background(), credentials)
		require.NoError(st, err)
		require.Equal(st, credentials.Username, acc.Username)
		require.Empty(st, acc.PasswordHash)
	})

	t.Run("KO - credentials error", func(st *testing.T) {
		_, err := strategy.Verify(context.Background(), nil)
		require.ErrorContains(st, err, "PASSPORT.CREDENTIALS")

		_, err = strategy.Verify(context.Background(), &entities.Credentials{})
		require.ErrorContains(st, err, "PASSPORT.CREDENTIALS")
	})

	t.Run("KO - user not found", func(st *testing.T) {
		credentials := &entities.Credentials{
			Username: uuid.NewString(),
			Password: testdata.Fake.Internet().Password(),
		}
		_, err := strategy.Verify(context.Background(), credentials)
		require.ErrorIs(st, err, ErrLogin)
	})

	t.Run("KO - deactivated error", func(st *testing.T) {
		// setup another batch to test deactivated account
		acc, pass := setup(t)
		require.NoError(st, orm.Create(acc).Error)

		i := testdata.Fake.IntBetween(0, len(pass)-1)
		credentials := &entities.Credentials{
			Username: acc[i].Username,
			Password: pass[i],
		}
		err := orm.
			Model(&entities.Account{}).
			Where("username = ?", credentials.Username).
			Update("deactivated_at", time.Now().Add(-time.Hour).UnixMilli()).
			Error
		require.NoError(st, err)

		_, err = strategy.Verify(context.Background(), credentials)
		require.ErrorIs(st, err, ErrAccountDeactivated)
	})

	t.Run("KO - password not match", func(st *testing.T) {
		i := testdata.Fake.IntBetween(0, len(passwords)-1)
		credentials := &entities.Credentials{
			Username: accounts[i].Username,
			Password: testdata.Fake.Internet().Password(),
		}
		_, err := strategy.Verify(context.Background(), credentials)
		require.ErrorIs(st, err, ErrLogin)
	})
}

func TestDurability_Register(t *testing.T) {
	accounts, _ := setup(t)

	conf := &config.Durability{Sqlx: sqlxconfig.Config{
		Uri: testdata.SqliteUri,
		Connection: sqlxconfig.Connection{
			MaxLifetime:  sqlxconfig.DefaultConnMaxLifetime,
			MaxIdletime:  sqlxconfig.DefaultConnMaxIdletime,
			MaxIdleCount: sqlxconfig.DefaultConnMaxIdleCount,
			MaxOpenCount: sqlxconfig.DefaultConnMaxOpenCount,
		},
	}}
	strategy, err := NewDurability(conf, testify.Logger())
	require.NoError(t, err)

	strategy.Connect(context.Background())
	defer strategy.Disconnect(context.Background())

	orm := strategy.(*durability).orm
	tx := orm.Create(accounts)
	require.NoError(t, tx.Error)

	t.Run("OK", func(st *testing.T) {
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

	t.Run("KO", func(st *testing.T) {
		i := testdata.Fake.IntBetween(0, len(accounts)-1)
		err := strategy.Register(context.Background(), &accounts[i])
		require.ErrorIs(st, err, ErrRegister)
	})
}

func TestDurability_Deactivate(t *testing.T) {
	accounts, _ := setup(t)

	conf := &config.Durability{Sqlx: sqlxconfig.Config{
		Uri: testdata.SqliteUri,
		Connection: sqlxconfig.Connection{
			MaxLifetime:  sqlxconfig.DefaultConnMaxLifetime,
			MaxIdletime:  sqlxconfig.DefaultConnMaxIdletime,
			MaxIdleCount: sqlxconfig.DefaultConnMaxIdleCount,
			MaxOpenCount: sqlxconfig.DefaultConnMaxOpenCount,
		},
	}}
	strategy, err := NewDurability(conf, testify.Logger())
	require.NoError(t, err)

	strategy.Connect(context.Background())
	defer strategy.Disconnect(context.Background())

	orm := strategy.(*durability).orm
	tx := orm.Create(accounts)
	require.NoError(t, tx.Error)

	t.Run("OK", func(st *testing.T) {
		i := testdata.Fake.IntBetween(0, len(accounts)-1)
		username := accounts[i].Username
		ts := time.Now().UnixMilli()

		err := strategy.Deactivate(context.Background(), username, ts)
		require.NoError(st, err)
	})

	t.Run("KO - user not found", func(st *testing.T) {
		username := uuid.NewString()
		ts := time.Now().UnixMilli()

		err := strategy.Deactivate(context.Background(), username, ts)
		require.ErrorIs(st, err, ErrAccountNotFound)
	})

	t.Run("KO - deactivate timestamp error", func(st *testing.T) {
		i := testdata.Fake.IntBetween(0, len(accounts)-1)
		username := accounts[i].Username
		ts := time.Now().UnixMilli()

		err := strategy.Deactivate(context.Background(), username, ts)
		require.NoError(st, err)

		err = strategy.Deactivate(context.Background(), username, ts-1)
		require.ErrorIs(st, err, ErrDeactivate)
	})
}

var durabilityconf = &config.Durability{Sqlx: sqlxconfig.Config{
	Uri: testdata.SqliteUri,
	Connection: sqlxconfig.Connection{
		MaxLifetime:  sqlxconfig.DefaultConnMaxLifetime,
		MaxIdletime:  sqlxconfig.DefaultConnMaxIdletime,
		MaxIdleCount: sqlxconfig.DefaultConnMaxIdleCount,
		MaxOpenCount: sqlxconfig.DefaultConnMaxOpenCount,
	},
}}
