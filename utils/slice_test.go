package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChunkNext(t *testing.T) {
	t.Run("begin", func(st *testing.T) {
		prev := 0
		end := 17
		step := 5

		next := ChunkNext(prev, end, step)
		require.Equal(st, next, 5)
	})

	t.Run("middle", func(st *testing.T) {
		prev := 5
		end := 17
		step := 5

		next := ChunkNext(prev, end, step)
		require.Equal(st, next, 10)
	})

	t.Run("end", func(st *testing.T) {
		prev := 15
		end := 17
		step := 5

		next := ChunkNext(prev, end, step)
		require.Equal(st, next, 17)
	})
}
