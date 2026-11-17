package response

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type Error struct {
	Code    string          `json:"error"`
	Message string          `json:"message"`
	Details json.RawMessage `json:"details"`
}

const cleanupAttempts = 3

// RetryCleanup повторяет идемпотентную очистку при временном отказе сервера
func RetryCleanup(t *testing.T, request func() *http.Response) *http.Response {
	t.Helper()
	for attempt := 1; ; attempt++ {
		raw := request()
		transient := raw.StatusCode == http.StatusTooManyRequests ||
			raw.StatusCode == http.StatusBadGateway ||
			raw.StatusCode == http.StatusServiceUnavailable ||
			raw.StatusCode == http.StatusGatewayTimeout
		if !transient || attempt == cleanupAttempts {
			return raw
		}
		_, err := io.Copy(io.Discard, raw.Body)
		closeErr := raw.Body.Close()
		require.NoError(t, err, "Не удалось прочитать ответ при очистке")
		require.NoError(t, closeErr, "Не удалось закрыть ответ при очистке")
		t.Logf("[Очистка] HTTP %d, повтор %d/%d", raw.StatusCode, attempt+1, cleanupAttempts)
		time.Sleep(time.Second)
	}
}

// ReadJSON проверяет статус, закрывает тело ответа и разбирает JSON в заданную модель
func ReadJSON[T any](t *testing.T, raw *http.Response, expectedStatus int) T {
	t.Helper()
	defer raw.Body.Close()
	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err, "Не удалось прочитать HTTP-ответ")
	require.Equal(t, expectedStatus, raw.StatusCode, "Неожиданный HTTP-статус; ответ: %s", body)
	var result T
	require.NoError(t, json.Unmarshal(body, &result), "Некорректный JSON в HTTP-ответе: %s", body)
	return result
}

// RequireStatus проверяет статус ответа без JSON, включая повторную очистку удалённых данных
func RequireStatus(t *testing.T, raw *http.Response, expectedStatuses ...int) {
	t.Helper()
	defer raw.Body.Close()
	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err, "Не удалось прочитать HTTP-ответ")
	require.Contains(t, expectedStatuses, raw.StatusCode, "Неожиданный HTTP-статус; ответ: %s", body)
}
