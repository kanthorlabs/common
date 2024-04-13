package configuration

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/project"
	"github.com/kanthorlabs/common/safe"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestNewFile(t *testing.T) {
	t.Run("OK - $KANTHOR_HOME", func(st *testing.T) {
		home := "/tmp/" + uuid.NewString()
		require.NoError(st, os.Mkdir(home, 0755))

		st.Setenv("KANTHOR_HOME", home)
		_, data := setupdata(t)
		require.NoError(st, os.WriteFile(home+"/"+FileName+"."+FileExt, data, 0644))

		provider, err := NewFile(project.Namespace(), FileLookingDirs)
		require.NoError(st, err)
		require.NotNil(st, provider)
	})

	t.Run("KO - empty dirs error", func(st *testing.T) {
		_, err := NewFile(project.Namespace(), []string{})
		require.ErrorContains(st, err, "CONFIGURATION.FILE.NO_DIRECTOY.ERROR")
	})

	t.Run("KO - read config error", func(st *testing.T) {
		home := "/tmp/" + uuid.NewString()
		require.NoError(st, os.Mkdir(home, 0755))

		st.Setenv("KANTHOR_HOME", home)
		require.NoError(st, os.WriteFile(home+"/"+FileName+"."+FileExt, []byte("test"), 0644))

		_, err := NewFile(project.Namespace(), FileLookingDirs)
		require.ErrorContains(st, err, "CONFIGURATION.FILE.ERROR")
	})
}

func TestFile_Unmarshal(t *testing.T) {
	t.Run("OK - $KANTHOR_HOME", func(st *testing.T) {
		orignal, data := setupdata(t)

		home := "/tmp/kanthor-" + uuid.NewString()
		require.NoError(st, os.Mkdir(home, 0755))

		st.Setenv("KANTHOR_HOME", home)
		require.NoError(st, os.WriteFile(home+"/"+FileName+"."+FileExt, data, 0644))

		provider, err := NewFile(project.Namespace(), FileLookingDirs)
		require.NoError(st, err)

		var conf configs
		require.NoError(st, provider.Unmarshal(&conf))
		require.Equal(st, orignal.Counter, conf.Counter)
		require.Equal(st, orignal.Blood, conf.Blood)

		ob, _ := orignal.Metadata.Get("bool")
		b, _ := conf.Metadata.Get("bool")
		require.Equal(st, ob, b)

		on, _ := orignal.Metadata.Get("number")
		n, _ := conf.Metadata.Get("number")
		require.Equal(st, on, int64(n.(int)))
	})

	t.Run("OK - $HOME", func(st *testing.T) {
		orignal, data := setupdata(t)

		home := "/tmp/" + uuid.NewString()
		require.NoError(st, os.Mkdir(home, 0755))
		require.NoError(st, os.Mkdir(home+"/.kanthor/", 0755))

		st.Setenv("HOME", home)
		require.NoError(st, os.WriteFile(home+"/.kanthor/"+FileName+"."+FileExt, data, 0644))

		provider, err := NewFile(project.Namespace(), FileLookingDirs)
		require.NoError(st, err)

		var conf configs
		require.NoError(st, provider.Unmarshal(&conf))
		require.Equal(st, orignal.Counter, conf.Counter)
		require.Equal(st, orignal.Blood, conf.Blood)
	})

	t.Run("OK - current directory", func(st *testing.T) {
		orignal, data := setupdata(t)
		require.NoError(st, os.WriteFile("./"+FileName+"."+FileExt, data, 0644))

		provider, err := NewFile(project.Namespace(), FileLookingDirs)
		require.NoError(st, err)

		var conf configs
		require.NoError(st, provider.Unmarshal(&conf))
		require.Equal(st, orignal.Counter, conf.Counter)
		require.Equal(st, orignal.Blood, conf.Blood)
	})
}

func TestFile_SetDefault(t *testing.T) {
	setupfile(t)

	provider, err := NewFile(project.Namespace(), FileLookingDirs)
	require.NoError(t, err)

	id := uuid.NewString()
	provider.SetDefault("unset_by_default", id)

	var conf configs
	require.NoError(t, provider.Unmarshal(&conf))

	require.Equal(t, id, conf.UnsetByDefault)
}

func TestFile_Sources(t *testing.T) {
	setupfile(t)

	_, data := setupdata(t)
	require.NoError(t, os.WriteFile("./"+FileName+"."+FileExt, data, 0644))

	provider, err := NewFile(project.Namespace(), FileLookingDirs)
	require.NoError(t, err)

	sources := provider.Sources()
	require.Len(t, sources, 3)
	require.Equal(t, sources[0].Looking, "$KANTHOR_HOME/configs.yaml")
	require.Equal(t, sources[1].Looking, "$HOME/.kanthor/configs.yaml")
	require.Equal(t, sources[2].Looking, "configs.yaml")
	require.False(t, sources[0].Used)
	require.True(t, sources[1].Used)
	require.False(t, sources[2].Used)
}

func TestFile_Set(t *testing.T) {
	setupfile(t)

	provider, err := NewFile(project.Namespace(), FileLookingDirs)
	require.NoError(t, err)

	provider.SetDefault("unset_by_default", uuid.NewString())
	id := uuid.NewString()
	provider.Set("unset_by_default", id)

	var conf configs
	require.NoError(t, provider.Unmarshal(&conf))

	require.Equal(t, id, conf.UnsetByDefault)
}

func setupfile(t *testing.T) {
	_, data := setupdata(t)

	home := "/tmp/" + uuid.NewString()
	require.NoError(t, os.Mkdir(home, 0755))
	require.NoError(t, os.Mkdir(home+"/.kanthor/", 0755))

	t.Setenv("HOME", home)
	require.NoError(t, os.WriteFile(home+"/.kanthor/"+FileName+"."+FileExt, data, 0644))
}

func setupdata(t *testing.T) (*configs, []byte) {
	conf := &configs{
		Counter:  testdata.Fake.IntBetween(1, 100),
		Blood:    testdata.Fake.Blood().Name(),
		Metadata: &safe.Metadata{},
	}
	conf.Metadata.Set("id", uuid.NewString())
	conf.Metadata.Set("bool", true)
	conf.Metadata.Set("number", testdata.Fake.Int64Between(1, 100))

	data, err := yaml.Marshal(conf)
	require.NoError(t, err)
	return conf, data
}
