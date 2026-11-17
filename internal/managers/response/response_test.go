package response

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRetryCleanup(t *testing.T) {
	cases := []struct {
		name     string
		statuses []int
	}{
		{"успех без повтора", []int{http.StatusNoContent}},
		{"временный отказ, затем - успех", []int{http.StatusServiceUnavailable, http.StatusNoContent}},
		{"ошибка запроса не повторяется", []int{http.StatusBadRequest}},
		{"число попыток ограничено", []int{http.StatusServiceUnavailable, http.StatusServiceUnavailable, http.StatusServiceUnavailable}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			raw := RetryCleanup(t, func() *http.Response {
				require.Less(t, calls, len(tc.statuses), "Избыточный повтор HTTP-запроса")
				status := tc.statuses[calls]
				calls++
				return &http.Response{StatusCode: status, Body: io.NopCloser(strings.NewReader(""))}
			})
			require.NoError(t, raw.Body.Close())
			require.Equal(t, len(tc.statuses), calls)
			require.Equal(t, tc.statuses[len(tc.statuses)-1], raw.StatusCode)
		})
	}
}
