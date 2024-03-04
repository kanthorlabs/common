package streaming

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/common/testify"
	"github.com/stretchr/testify/require"
)

func TestStreaming_New(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		t.Run("OK", func(st *testing.T) {
			_, err := New(testconf("nats://127.0.0.1:42222"), testify.Logger())
			require.NoError(st, err)
		})
	})
}

func fakeitems(count int) map[string]*testdata.User {
	items := map[string]*testdata.User{}

	for i := 0; i < count; i++ {
		id := uuid.NewString()
		items[id] = &testdata.User{
			Id:       id,
			Username: testdata.Fake.Internet().Email(),
			Created:  time.Now().UnixMilli(),
			Updated:  time.Now().UnixMilli(),
			Metadata: map[string]string{
				"session_id": uuid.NewString(),
			},
		}
	}

	return items
}

func streamname() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "") + "_" + strings.ReplaceAll(uuid.NewString(), "-", "")
}

func subjectname() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "") + "." + strings.ReplaceAll(uuid.NewString(), "-", "")
}

func topicname() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "") + "." + strings.ReplaceAll(uuid.NewString(), "-", "")
}

func pubsubname() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "")
}
