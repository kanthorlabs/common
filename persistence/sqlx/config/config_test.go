package config

import (
	"testing"

	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestDefault(t *testing.T) {
	conf := &Config{}
	require.ErrorContains(t, conf.Validate(), "SQLX.CONFIG.URI")

	conf.Uri = testdata.SqliteUri
	require.ErrorContains(t, conf.Validate(), "SQLX.CONFIG.CONNECTION")

	conf.Connection = Connection{
		MaxLifetime:  DefaultConnMaxLifetime,
		MaxIdletime:  DefaultConnMaxIdletime,
		MaxIdleCount: DefaultConnMaxIdleCount,
		MaxOpenCount: DefaultConnMaxOpenCount,
	}
	require.NoError(t, conf.Validate())
}
