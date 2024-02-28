package testdata

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/clock"
	"github.com/segmentio/ksuid"
)

var (
	UserNs = "u"
)

type User struct {
	Id       string `gorm:"primaryKey"`
	Username string
	Created  int64
	Updated  int64
}

func NewUser(watch clock.Clock) User {
	now := watch.Now()
	// error could not be happen because we provide a valid payload
	id, _ := ksuid.FromParts(now, []byte("0000000000000000"))

	return User{
		Id:       fmt.Sprintf("%s_%s", UserNs, id.String()),
		Username: uuid.NewString() + "/" + Fake.Internet().Email(),
		Created:  now.UnixMilli(),
		Updated:  now.UnixMilli(),
	}
}
