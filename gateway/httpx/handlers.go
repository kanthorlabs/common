package httpx

import (
	"net/http"
	"time"

	"github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/common/project"
)

func UseHealthz(check func() error) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := check(); err != nil {
			writer.ErrUnknown(w, err)
			return
		}

		writer.Ok(w, writer.M{"timestamp": time.Now().UTC().Format(time.RFC3339Nano)})
	})
}

func UseVersion() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer.Ok(w, writer.M{"version": project.GetVersion()})
	})
}
