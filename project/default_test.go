package project

import (
	"testing"

	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestRegion(t *testing.T) {
	t.Run("default", func(st *testing.T) {
		got := Region()
		require.Equal(st, DefaultRegion, got)
	})

	t.Run("from env", func(st *testing.T) {
		region := testdata.Fake.Address().CountryCode()
		st.Setenv("PROJECT_REGION", region)
		got := Region()
		require.Equal(st, region, got)
	})
}

func TestNamespace(t *testing.T) {
	t.Run("default", func(st *testing.T) {
		got := Namespace()
		require.Equal(st, DefaultNamespace, got)
	})

	t.Run("from env", func(st *testing.T) {
		ns := testdata.Fake.App().Name()
		st.Setenv("PROJECT_NAMESPACE", ns)
		got := Namespace()
		require.Equal(st, ns, got)
	})
}

func TestTier(t *testing.T) {
	t.Run("default", func(st *testing.T) {
		got := Tier()
		require.Equal(st, DefaultTier, got)
	})

	t.Run("from env", func(st *testing.T) {
		tier := testdata.Fake.Blood().Name()
		st.Setenv("PROJECT_TIER", tier)
		got := Tier()
		require.Equal(st, tier, got)
	})
}
