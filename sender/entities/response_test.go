package entities

import (
	"net/http"
	"testing"

	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestResponse(t *testing.T) {
	testcases := map[int]bool{
		http.StatusSwitchingProtocols:  false,
		http.StatusOK:                  true,
		http.StatusBadRequest:          false,
		http.StatusInternalServerError: false,
	}

	for status := range testcases {
		res := Response{Status: status}
		require.Equal(t, res.Ok(), testcases[status])
		require.Equal(t, res.StatusText(), http.StatusText(status))
	}

	exception := &Response{Status: -1, Body: testdata.Fake.Lorem().Bytes(100)}
	require.Equal(t, exception.StatusText(), string(exception.Body))

}
