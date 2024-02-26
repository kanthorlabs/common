package cache

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	mcache "github.com/kanthorlabs/common/mocks/cache"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetOrSet(t *testing.T) {
	c := mcache.NewCache(t)
	key := uuid.NewString()
	value := uuid.NewString()
	entry, _ := json.Marshal(value)
	ttl := time.Hour

	t.Run("OK - insert to cache", func(st *testing.T) {
		c.EXPECT().Get(mock.Anything, key).Return(nil, ErrEntryNotFound).Once()
		c.EXPECT().Set(mock.Anything, key, entry, ttl).Return(nil).Once()

		v, err := GetOrSet[string](c, context.Background(), key, ttl, func() (*string, error) {
			return &value, nil
		})

		require.Nil(st, err)
		require.Equal(st, value, *v)
	})

	t.Run("OK - get from cache", func(st *testing.T) {
		c.EXPECT().Get(mock.Anything, key).Return(entry, nil).Once()

		v, err := GetOrSet[string](c, context.Background(), key, ttl, func() (*string, error) {
			x := uuid.NewString()
			return &x, nil
		})

		require.Nil(st, err)
		require.Equal(st, value, *v)
	})

	t.Run("KO - get error", func(st *testing.T) {
		expected := errors.New(testdata.Fake.Lorem().Sentence(1))
		c.EXPECT().Get(mock.Anything, key).Return(nil, expected).Once()

		_, err := GetOrSet[string](c, context.Background(), key, ttl, func() (*string, error) {
			return &value, nil
		})

		require.ErrorIs(st, err, expected)
	})

	t.Run("KO - unmarshal error", func(st *testing.T) {
		c.EXPECT().Get(mock.Anything, key).Return([]byte(uuid.NewString()), nil).Once()

		_, err := GetOrSet[string](c, context.Background(), key, ttl, func() (*string, error) {
			x := uuid.NewString()
			return &x, nil
		})

		require.ErrorContains(st, err, "invalid character")
	})

	t.Run("KO - handler error", func(st *testing.T) {
		expected := errors.New(testdata.Fake.Lorem().Sentence(1))

		c.EXPECT().Get(mock.Anything, key).Return(nil, ErrEntryNotFound).Once()

		_, err := GetOrSet[string](c, context.Background(), key, ttl, func() (*string, error) {
			return nil, expected
		})

		require.ErrorIs(st, err, expected)
	})

	t.Run("KO - marshal error", func(st *testing.T) {
		c.EXPECT().Get(mock.Anything, key).Return(nil, ErrEntryNotFound).Once()

		_, err := GetOrSet[person](c, context.Background(), key, ttl, func() (*person, error) {
			parent := &person{
				Name: testdata.Fake.Person().FirstName(),
				Children: []*person{
					{Name: testdata.Fake.Person().FirstName()},
					{Name: testdata.Fake.Person().FirstName()},
				},
			}
			parent.Children[0].Parent = parent
			parent.Children[1].Parent = parent
			return parent, nil
		})

		require.ErrorContains(st, err, "json: unsupported value")
	})

	t.Run("KO - set error", func(st *testing.T) {
		expected := errors.New(testdata.Fake.Lorem().Sentence(1))

		c.EXPECT().Get(mock.Anything, key).Return(nil, ErrEntryNotFound).Once()
		c.EXPECT().Set(mock.Anything, key, entry, ttl).Return(expected).Once()

		_, err := GetOrSet[string](c, context.Background(), key, ttl, func() (*string, error) {
			return &value, nil
		})

		require.ErrorIs(st, err, expected)
	})
}

type person struct {
	Name     string
	Parent   *person
	Children []*person
}
