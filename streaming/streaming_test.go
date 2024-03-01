package streaming

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/testdata"
)

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
