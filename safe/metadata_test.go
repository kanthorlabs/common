package safe

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/testdata"
	"github.com/sourcegraph/conc"
	"github.com/stretchr/testify/require"
)

func TestMetadata(t *testing.T) {
	t.Run(".Set", func(st *testing.T) {
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

		require.Equal(st, counter, len(metadata.kv))
	})

	t.Run(".Set", func(st *testing.T) {
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

		wg.Go(func() {
			v, has := metadata.Get(fmt.Sprintf("index_%d", counter-1))
			require.True(st, has)
			require.Equal(st, counter-1, v)
		})
		wg.Wait()
	})

	t.Run(".Merge", func(st *testing.T) {
		dest := &Metadata{}

		dest.Merge(nil)
		require.Equal(st, 0, len(dest.kv))

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
		require.Equal(st, counter, len(dest.kv))

		dest.Merge(&Metadata{})
		require.Equal(st, counter, len(dest.kv))
	})

	t.Run(".String", func(st *testing.T) {
		st.Run("OK", func(sst *testing.T) {
			var metadata Metadata
			id := uuid.NewString()
			metadata.Set(id, true)

			str := metadata.String()
			require.Contains(sst, str, id)
		})

		st.Run("nullable", func(sst *testing.T) {
			var metadata Metadata
			require.Empty(sst, metadata.String())
		})
	})

	t.Run(".Value", func(st *testing.T) {
		st.Run("OK", func(sst *testing.T) {
			var metadata Metadata
			id := uuid.NewString()
			metadata.Set(id, true)

			value, err := metadata.Value()
			require.Nil(sst, err)
			require.Contains(sst, value, id)
		})

		st.Run("nullable", func(sst *testing.T) {
			var metadata Metadata
			value, err := metadata.Value()
			require.Nil(sst, err)
			require.Empty(sst, value)
		})
	})

	t.Run(".Scan", func(st *testing.T) {
		st.Run("OK", func(sst *testing.T) {
			var src Metadata
			src.Set(uuid.NewString(), true)

			var metadata Metadata
			require.Nil(sst, metadata.Scan(src.String()))

			require.Equal(sst, src.kv, metadata.kv)
		})

		st.Run("nullable", func(sst *testing.T) {
			var metadata Metadata
			require.Nil(sst, metadata.Scan(""))
		})
	})
}
