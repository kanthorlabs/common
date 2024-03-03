package project

import (
	"testing"

	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestRegion(t *testing.T) {
	t.Run("default", func(st *testing.T) {
		require.Equal(st, DefaultRegion, Region())
	})

	t.Run("from env", func(st *testing.T) {
		region := testdata.Fake.Address().CountryCode()
		st.Setenv("PROJECT_REGION", region)
		require.Equal(st, region, Region())
	})
}

func TestNamespace(t *testing.T) {
	t.Run("default", func(st *testing.T) {
		require.Equal(st, DefaultNamespace, Namespace())
	})

	t.Run("from env", func(st *testing.T) {
		ns := testdata.Fake.App().Name()
		st.Setenv("PROJECT_NAMESPACE", ns)
		require.Equal(st, ns, Namespace())
	})
}

func TestTier(t *testing.T) {
	t.Run("default", func(st *testing.T) {
		require.Equal(st, DefaultTier, Tier())
	})

	t.Run("from env", func(st *testing.T) {
		tier := testdata.Fake.Blood().Name()
		st.Setenv("PROJECT_TIER", tier)
		require.Equal(st, tier, Tier())
	})
}
