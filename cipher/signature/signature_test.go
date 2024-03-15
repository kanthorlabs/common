package signature

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/testdata"
)

var key = uuid.NewString()
var data = fmt.Sprintf(
	"%s.%s.%d",
	testdata.Fake.Lorem().Sentence(1),
	testdata.Fake.Lorem().Sentence(1),
	time.Now().UTC().UnixMilli(),
)
