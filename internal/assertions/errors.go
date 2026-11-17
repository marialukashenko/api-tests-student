package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/response"
)

// сравнение кода ошибки и сообщения от сервера
func APIError(t *testing.T, actual response.Error, code, message string) {
	t.Helper()
	require.Equal(t, code, actual.Code, "Неверный код ошибки API")
	require.Equal(t, message, actual.Message, "Неверное сообщение об ошибке API")
	t.Logf("[Результат] Получена ожидаемая ошибка %s: %s", code, message)
}
