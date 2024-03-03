package safe

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/testdata"
	"github.com/sourcegraph/conc"
	"github.com/stretchr/testify/require"
)

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
