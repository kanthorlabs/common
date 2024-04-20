package circuitbreaker

import (
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/circuitbreaker/config"
	"github.com/kanthorlabs/common/clock"
	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/common/testify"
	"github.com/stretchr/testify/require"
)

func TestGoBreaker_New(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		_, err := NewGoBreaker(testconf, testify.Logger())
		require.NoError(st, err)
	})

	t.Run(testify.CaseKOValidationError, func(st *testing.T) {
		conf := &config.Config{}
		_, err := NewGoBreaker(conf, testify.Logger())
		require.ErrorContains(st, err, "CIRCUIT_BREAKER.CONFIG")
	})
}

func TestGoBreaker_Do(t *testing.T) {
	t.Run("KO - consecutive error", func(st *testing.T) {
		cb, err := NewGoBreaker(testconf, testify.Logger())
		require.NoError(st, err)

		cmd := uuid.NewString()
		for {
			_, err = cb.Do(cmd, func() (any, error) {
				return nil, errors.New(testdata.Fake.RandomStringWithLength(10))
			}, passerror)

			if err != nil && strings.Contains(err.Error(), "CIRCUIT_BREAKER.STAGE_OPEN.ERROR") {
				break
			}
		}
	})

	t.Run("KO - ratio error", func(st *testing.T) {
		cb, err := NewGoBreaker(testconf, testify.Logger())
		require.NoError(st, err)

		cmd := uuid.NewString()
		errorable := false
		for {
			_, err = cb.Do(cmd, func() (any, error) {
				errorable = !errorable

				if errorable {
					return nil, errors.New(testdata.Fake.RandomStringWithLength(10))
				}

				return testdata.NewUser(clock.New()), nil
			}, passerror)

			if err != nil && strings.Contains(err.Error(), "CIRCUIT_BREAKER.STAGE_OPEN.ERROR") {
				break
			}
		}
	})
}
