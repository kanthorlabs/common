package validator

import (
	"fmt"
	"time"
)

var MinDatetime = time.Date(2014, 5, 30, 17, 0, 0, 0, time.UTC)

func DatetimeBefore(prop string, value, target time.Time) Fn {
	return func() error {
		if value.Before(MinDatetime) {
			return fmt.Errorf("%s (%s) must after %s", prop, value.Format(time.RFC3339Nano), MinDatetime.Format(time.RFC3339Nano))
		}

		if target.Before(MinDatetime) {
			return fmt.Errorf("%s (%s) must after %s", prop, target.Format(time.RFC3339Nano), MinDatetime.Format(time.RFC3339Nano))
		}

		if value.After(target) {
			return fmt.Errorf("%s (%s) must before %s", prop, value.Format(time.RFC3339Nano), target.Format(time.RFC3339Nano))
		}

		return nil
	}
}
