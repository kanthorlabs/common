package safe

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/testdata"
	"github.com/sourcegraph/conc"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

type Value struct {
	Metadata *Metadata `json:"metadata" yaml:"metadata" mapstructure:"metadata"`
}

func TestMetadata_Parsing(t *testing.T) {
	value := Value{
		Metadata: &Metadata{},
	}
	value.Metadata.Set("bool", testdata.Fake.Bool())
	value.Metadata.Set("number", testdata.Fake.Int64Between(1, 100))

	t.Run("json", func(st *testing.T) {
		jsonv, err := json.Marshal(value)
		require.NoError(st, err)
		require.Contains(st, string(jsonv), "bool")
		require.Contains(st, string(jsonv), "number")

		var jsonp Value
		require.NoError(st, json.Unmarshal(jsonv, &jsonp))
		require.Equal(st, value.Metadata.kv["bool"], jsonp.Metadata.kv["bool"])
		require.Equal(st, float64(value.Metadata.kv["number"].(int64)), jsonp.Metadata.kv["number"])
	})

	t.Run("yaml", func(st *testing.T) {
		yamlv, err := yaml.Marshal(value)
		require.NoError(st, err)
		require.Contains(st, string(yamlv), "bool")
		require.Contains(st, string(yamlv), "number")

		var yamlp Value
		require.NoError(st, yaml.Unmarshal(yamlv, &yamlp))
		require.Equal(st, value.Metadata.kv["bool"], yamlp.Metadata.kv["bool"])
		require.Equal(st, int(value.Metadata.kv["number"].(int64)), yamlp.Metadata.kv["number"])
	})
}

func TestMetadata_Set(t *testing.T) {
	var metadata Metadata

	var wg conc.WaitGroup
	counter := testdata.Fake.IntBetween(100000, 999999)
	for i := 0; i < counter; i++ {
		index := i
		wg.Go(func() {
			metadata.Set(fmt.Sprintf("index_%d", index), index)
		})
	}
	wg.Wait()

	require.Equal(t, counter, len(metadata.kv))
}

func TestMetadata_Get(t *testing.T) {
	var metadata Metadata
	counter := testdata.Fake.IntBetween(100000, 999999)

	v, has := metadata.Get(fmt.Sprintf("index_%d", counter+1))
	require.False(t, has)
	require.Nil(t, v)

	var wg conc.WaitGroup
	for i := 0; i < counter; i++ {
		index := i
		wg.Go(func() {
			metadata.Set(fmt.Sprintf("index_%d", index), index)
		})
	}
	wg.Wait()

	wg.Go(func() {
		v, has := metadata.Get(fmt.Sprintf("index_%d", counter-1))
		require.True(t, has)
		require.Equal(t, counter-1, v)
	})
	wg.Wait()
}

func TestMetadata_Merge(t *testing.T) {
	dest := &Metadata{}

	dest.Merge(nil)
	require.Equal(t, 0, len(dest.kv))

	src := &Metadata{}

	var wg conc.WaitGroup
	counter := testdata.Fake.IntBetween(100000, 999999)
	for i := 0; i < counter; i++ {
		index := i
		wg.Go(func() {
			src.Set(fmt.Sprintf("index_%d", index), index)
		})
	}
	wg.Wait()

	dest.Merge(src)
	require.Equal(t, counter, len(dest.kv))

	dest.Merge(&Metadata{})
	require.Equal(t, counter, len(dest.kv))
}

func TestMetadata_String(t *testing.T) {
	var metadata Metadata
	require.Empty(t, metadata.String())

	metadata.Set(uuid.NewString(), true)
	require.NotEmpty(t, metadata.String())
}

func TestMetadata_Value(t *testing.T) {
	var metadata Metadata
	emptv, emptyerr := metadata.Value()
	require.NoError(t, emptyerr)
	require.Empty(t, emptv)

	id := uuid.NewString()
	metadata.Set(id, true)

	value, err := metadata.Value()
	require.NoError(t, err)
	require.Contains(t, value, id)
}

func TestMetadata_Scan(t *testing.T) {
	var metadata Metadata
	require.NoError(t, metadata.Scan(""))

	var src Metadata
	src.Set(uuid.NewString(), true)
	require.NoError(t, metadata.Scan(src.String()))
	require.Equal(t, src.kv, metadata.kv)
}

func TestMetadata_MetadataMapstructureHook(t *testing.T) {
	hook := MetadataMapstructureHook()
	t.Run("OK", func(st *testing.T) {
		data := map[string]interface{}{"bool": true, "number": 1}
		from := reflect.ValueOf(data)
		to := reflect.ValueOf(&Metadata{})

		metdata, err := hook(from.Type(), to.Type(), data)
		require.NoError(st, err)

		for k, v := range data {
			value, exist := metdata.(*Metadata).Get(k)
			require.True(st, exist)
			require.Equal(st, v, value)
		}
	})

	t.Run("OK - not metadata pointer", func(st *testing.T) {
		data := map[string]interface{}{"bool": true, "number": 1}
		from := reflect.ValueOf(data)
		to := reflect.ValueOf(Metadata{})

		metdata, err := hook(from.Type(), to.Type(), data)
		require.NoError(st, err)

		require.Equal(st, data, metdata)
	})
}
