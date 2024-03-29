package utils

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/project"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func Test_IsBasicScheme(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		require.True(st, IsBasicScheme(SchemeBasic))
	})

	t.Run("OK - case insensitive", func(st *testing.T) {
		require.True(st, IsBasicScheme(strings.ToLower(SchemeBasic)))
	})

	t.Run("KO - empty error", func(st *testing.T) {
		require.False(st, IsBasicScheme(""))
	})

	t.Run("KO - not matching error", func(st *testing.T) {
		require.False(st, IsBasicScheme("Bearer "))
	})
}

func Test_ParseBasicCredentials_Standard(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		credentials, err := ParseBasicCredentials(SchemeBasic + basic)
		require.NoError(st, err)
		require.Equal(st, user, credentials.Username)
		require.Equal(st, pass, credentials.Password)
		require.Empty(st, credentials.Region)
	})

	t.Run("KO - not basic scheme error", func(st *testing.T) {
		_, err := ParseBasicCredentials("Bearer ")
		require.ErrorContains(st, err, "PASSPORT.UTILS.PARSE_BASIC_CREDENTIALS.SCHEME_UNKNOWN.ERROR")
	})

	t.Run("KO - base64 error", func(st *testing.T) {
		_, err := ParseBasicCredentials(SchemeBasic + "invalid")
		require.ErrorContains(st, err, "PASSPORT.UTILS.PARSE_BASIC_CREDENTIALS.DECODE.ERROR")
	})

	t.Run("KO - not matching user:pass pattern error", func(st *testing.T) {
		_, err := ParseBasicCredentials(SchemeBasic + base64.StdEncoding.EncodeToString([]byte(user)))
		require.ErrorContains(st, err, "PASSPORT.UTILS.PARSE_BASIC_CREDENTIALS.PARSE.ERROR")
	})
}

func Test_ParseBasicCredentials_Regional(t *testing.T) {
	basicregion := CreateRegionalBasicCredentials(user + ":" + pass)
	t.Run("OK", func(st *testing.T) {
		credentials, err := ParseBasicCredentials(SchemeBasic + basicregion)
		require.NoError(st, err)
		require.Equal(st, user, credentials.Username)
		require.Equal(st, pass, credentials.Password)
		require.Equal(st, project.Region(), credentials.Region)
	})
}

var (
	user  = testdata.Fake.Internet().Email()
	pass  = uuid.NewString()
	basic = base64.StdEncoding.EncodeToString([]byte(user + ":" + pass))
)
