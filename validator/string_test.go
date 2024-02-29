package validator

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestStringRequired(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		require.Nil(st, StringRequired("name", testdata.Fake.App().Name())())
	})
	t.Run("KO", func(st *testing.T) {
		require.ErrorContains(st, StringRequired("name", "")(), "is required")
	})
}

func TestStringStartsWithIfNotEmpty(t *testing.T) {
	prefix := testdata.Fake.Company().Name()
	t.Run("OK - empty", func(st *testing.T) {
		require.Nil(st, StringStartsWithIfNotEmpty("name", "", prefix)())
	})
	t.Run("OK", func(st *testing.T) {
		require.Nil(st, StringStartsWithIfNotEmpty("name", prefix+" - "+testdata.Fake.Company().JobTitle(), prefix)())
	})
	t.Run("KO", func(st *testing.T) {
		require.ErrorContains(st, StringStartsWithIfNotEmpty("name", testdata.Fake.Company().CatchPhrase(), prefix)(), "must be started with")
	})
}

func TestStringStartsWith(t *testing.T) {
	prefix := testdata.Fake.Company().Name()
	t.Run("OK", func(st *testing.T) {
		require.Nil(st, StringStartsWith("name", prefix+" - "+testdata.Fake.Company().JobTitle(), prefix)())
	})
	t.Run("KO", func(st *testing.T) {
		require.ErrorContains(st, StringStartsWith("name", testdata.Fake.Company().CatchPhrase(), prefix)(), "must be started with")
	})
	t.Run("KO - empty", func(st *testing.T) {
		require.ErrorContains(st, StringStartsWith("name", "", prefix)(), "is required")
	})
}

func TestStringStartsWithOneOf(t *testing.T) {
	oneOf := []string{
		testdata.Fake.UserAgent().Chrome(),
		testdata.Fake.UserAgent().Firefox(),
		testdata.Fake.UserAgent().Safari(),
	}

	t.Run("OK", func(st *testing.T) {
		require.Nil(st, StringStartsWithOneOf("agent", oneOf[0], oneOf)())
	})
	t.Run("KO", func(st *testing.T) {
		require.ErrorContains(st, StringStartsWithOneOf("agent", testdata.Fake.UserAgent().Opera(), oneOf)(), "prefix must be started with one of")
	})
}

func TestStringUri(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		require.Nil(st, StringUri("uri", testdata.Fake.Internet().URL())())
	})
	t.Run("KO", func(st *testing.T) {
		require.ErrorContains(st, StringUri("uri", testdata.Fake.Internet().Email())(), "is not a valid uri")
	})
	t.Run("KO - empty", func(st *testing.T) {
		require.ErrorContains(st, StringUri("uri", "")(), "is required")
	})
}

func TestStringLen(t *testing.T) {
	id := testdata.Fake.UUID().V4()
	t.Run("OK", func(st *testing.T) {
		require.Nil(st, StringLen("id", id, len(id)/2, len(id)*2)())
	})
	t.Run("KO - greater than", func(st *testing.T) {
		require.ErrorContains(st, StringLen("id", id, len(id)*2, len(id)*4)(), "length must be greater than or equal")
	})
	t.Run("KO - less than", func(st *testing.T) {
		require.ErrorContains(st, StringLen("id", id, len(id)/4, len(id)/2)(), "length must be less than or equal")
	})
}

func TestStringLenIfNotEmpty(t *testing.T) {
	id := testdata.Fake.UUID().V4()
	t.Run("OK", func(st *testing.T) {
		require.Nil(st, StringLenIfNotEmpty("id", id, len(id)/2, len(id)*2)())
	})
	t.Run("Ok - empty", func(st *testing.T) {
		require.Nil(st, StringLenIfNotEmpty("id", "", len(id)/2, len(id)*2)())
	})
}

func TestStringOneOf(t *testing.T) {
	oneOf := []string{
		testdata.Fake.UserAgent().Chrome(),
		testdata.Fake.UserAgent().Firefox(),
		testdata.Fake.UserAgent().Safari(),
	}

	t.Run("OK", func(st *testing.T) {
		require.Nil(st, StringOneOf("agent", oneOf[0], oneOf)())
	})
	t.Run("KO", func(st *testing.T) {
		require.ErrorContains(st, StringOneOf("agent", testdata.Fake.UserAgent().Opera(), oneOf)(), "must be one of")
	})
	t.Run("KO - empty", func(st *testing.T) {
		require.ErrorContains(st, StringOneOf("agent", "", oneOf)(), "is required")
	})
}

func TestStringAlphaNumericUnderscore(t *testing.T) {
	name := strings.ReplaceAll(uuid.NewString(), "-", "")
	t.Run("OK", func(st *testing.T) {
		require.Nil(st, StringAlphaNumericUnderscore("name", name+"_"+name)())
	})
	t.Run("KO", func(st *testing.T) {
		require.ErrorContains(st, StringAlphaNumericUnderscore("name", name+"."+name)(), "is not matched")
	})
}

func TestStringAlphaNumericUnderscoreDot(t *testing.T) {
	name := strings.ReplaceAll(uuid.NewString(), "-", "")
	t.Run("OK", func(st *testing.T) {
		require.Nil(st, StringAlphaNumericUnderscoreDot("name", name+"_"+name+"."+name)())
	})
	t.Run("KO", func(st *testing.T) {
		require.ErrorContains(st, StringAlphaNumericUnderscoreDot("name", name+"_"+name+"."+name+"#"+name)(), "is not matched")
	})
}
