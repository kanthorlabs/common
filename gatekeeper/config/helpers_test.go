package config

import (
	"encoding/base64"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/gatekeeper/entities"
	"github.com/kanthorlabs/common/utils"
	"github.com/stretchr/testify/require"
)

func TestParseDefinitionsToPermissions(t *testing.T) {
	t.Run("KO - unsupported uri", func(st *testing.T) {
		_, err := ParseDefinitionsToPermissions("unsupported://")
		require.ErrorIs(t, err, ErrorDefinitionsUriUnsupported)
	})
}

func TestParseDefinitionsToPermissions_File(t *testing.T) {
	definitions := map[string][]entities.Permission{
		"administrator": {{Action: "*", Object: "*"}},
	}

	permission := utils.Stringify(definitions)

	dpath := os.TempDir() + "/" + uuid.NewString()
	require.NoError(t, os.WriteFile(dpath, []byte(permission), os.ModePerm))

	t.Run("OK", func(st *testing.T) {
		perms, err := ParseDefinitionsToPermissions("file://" + dpath)
		require.NoError(st, err)
		require.Equal(st, definitions, perms)
	})

	t.Run("KO - file not found", func(st *testing.T) {
		_, err := ParseDefinitionsToPermissions("file://./not-found/file")
		require.ErrorContains(st, err, ErrorDefinitionsFile.Error())
	})

	t.Run("KO - unmarshal error", func(st *testing.T) {
		p := os.TempDir() + "/" + uuid.NewString()
		require.NoError(st, os.WriteFile(p, []byte(""), os.ModePerm))
		_, err := ParseDefinitionsToPermissions("file://" + p)
		require.ErrorContains(st, err, ErrorDefinitionsFile.Error())
	})
}

func TestParseDefinitionsToPermissions_Base64(t *testing.T) {
	definitions := map[string][]entities.Permission{
		"administrator": {{Action: "*", Object: "*"}},
	}

	permission := utils.Stringify(definitions)
	dbase64 := base64.StdEncoding.EncodeToString([]byte(permission))

	t.Run("OK", func(st *testing.T) {
		perms, err := ParseDefinitionsToPermissions("base64://" + dbase64)
		require.NoError(st, err)
		require.Equal(st, definitions, perms)
	})

	t.Run("KO - decode error", func(st *testing.T) {
		_, err := ParseDefinitionsToPermissions("base64://---")
		require.ErrorContains(st, err, ErrorDefinitionsBase64.Error())
	})

	t.Run("KO - unmarshal error", func(st *testing.T) {
		_, err := ParseDefinitionsToPermissions("base64://ey19")
		require.ErrorContains(t, err, ErrorDefinitionsBase64.Error())
	})
}
