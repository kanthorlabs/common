package safe

import (
	"fmt"
	"testing"

	"github.com/kanthorlabs/common/testdata"
	"github.com/sourcegraph/conc"
	"github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {
	m := Map[error]{}

	var wg conc.WaitGroup
	counter := testdata.Fake.IntBetween(100000, 999999)
	for i := 0; i < counter; i++ {
		index := i
		wg.Go(func() {
			m.Set(fmt.Sprintf("index_%d", index), testdata.ErrGeneric)
		})
	}
	wg.Wait()

	v, has := m.Get(fmt.Sprintf("index_%d", counter-1))
	require.True(t, has)
	require.ErrorIs(t, testdata.ErrGeneric, v)

	require.ErrorIs(t, testdata.ErrGeneric, m.Sample())
	require.Equal(t, counter, m.Count())
	require.Equal(t, counter, len(m.Data()))
	require.Equal(t, counter, len(m.Keys()))

	m.Merge(map[string]error{"new": testdata.ErrGeneric})
	require.Equal(t, counter+1, m.Count())
}
