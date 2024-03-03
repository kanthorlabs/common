package safe

import (
	"testing"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/testdata"
	"github.com/sourcegraph/conc"
	"github.com/stretchr/testify/require"
)

func TestSlice(t *testing.T) {
	s := Slice[string]{}

	var wg conc.WaitGroup
	counter := testdata.Fake.IntBetween(100000, 999999)
	for i := 0; i < counter; i++ {
		wg.Go(func() {
			s.Append(uuid.NewString())
		})
	}
	wg.Wait()

	require.Equal(t, counter, s.Count())
	_, err := uuid.Parse(s.Data()[s.Count()-1])
	require.NoError(t, err)
}
