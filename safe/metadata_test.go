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
		metadata := &Metadata{}

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
		metadata := &Metadata{}

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
		metadata := &Metadata{}
		id := uuid.NewString()
		metadata.Set(id, true)

		str := metadata.String()
		require.Contains(st, str, id)
	})

	t.Run(".Value", func(st *testing.T) {
		metadata := &Metadata{}
		id := uuid.NewString()
		metadata.Set(id, true)

		value, err := metadata.Value()
		require.Nil(st, err)

		require.Contains(st, value, id)
	})

	t.Run(".Scan", func(st *testing.T) {
		src := &Metadata{}
		src.Set(uuid.NewString(), true)

		metadata := &Metadata{}
		metadata.Scan(src.String())

		require.Equal(st, src.kv, metadata.kv)
	})
}
