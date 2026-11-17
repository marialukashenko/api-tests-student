package favorites

import (
	"net/http"
	"testing"

	client "api-tests-template/internal/client/http/favorites"
	"api-tests-template/internal/managers/favorites/models"
	"api-tests-template/internal/managers/response"
)

// AddFavorite добавляет объявление в избранное и возвращает созданную запись
func AddFavorite(t *testing.T, token, id string) models.Favorite {
	t.Helper()
	t.Logf("[Шаг] Добавляем объявление %s в избранное", id)
	return response.ReadJSON[models.Favorite](t, client.HttpPostFavorite(t, token, id), http.StatusCreated)
}

// RejectAddFavorite проверяет статус неуспешного добавления и возвращает ошибку API
func RejectAddFavorite(t *testing.T, token, id string, status int) response.Error {
	t.Helper()
	t.Logf("[Шаг] Проверяем отказ в добавлении %s в избранное", id)
	return response.ReadJSON[response.Error](t, client.HttpPostFavorite(t, token, id), status)
}

// GetFavorites возвращает страницу избранного с фильтром по названию
func GetFavorites(t *testing.T, token, search string) models.List {
	t.Helper()
	t.Logf("[Шаг] Получаем избранное с фильтром %q", search)
	return response.ReadJSON[models.List](t, client.HttpGetFavorites(t, token, search, 1), http.StatusOK)
}

// RejectGetFavorites проверяет статус отказа в просмотре избранного
func RejectGetFavorites(t *testing.T, token string, status int) response.Error {
	t.Helper()
	t.Log("[Шаг] Проверяем отказ в получении избранного")
	return response.ReadJSON[response.Error](t, client.HttpGetFavorites(t, token, "", 1), status)
}

// DeleteFavorite удаляет запись из избранного с обязательным статусом 204
func DeleteFavorite(t *testing.T, token, id string) {
	t.Helper()
	t.Logf("[Шаг] Удаляем объявление %s из избранного", id)
	response.RequireStatus(t, client.HttpDeleteFavorite(t, token, id), http.StatusNoContent)
}

// CleanupFavorite очищает тестовую запись, допуская её предыдущее удаление
func CleanupFavorite(t *testing.T, token, id string) {
	t.Helper()
	t.Logf("[Очистка] Удаляем тестовую запись избранного %s", id)
	raw := response.RetryCleanup(t, func() *http.Response {
		return client.HttpDeleteFavorite(t, token, id)
	})
	response.RequireStatus(t, raw, http.StatusNoContent, http.StatusNotFound)
}

// RejectDeleteFavorite проверяет статус отказа в удалении записи из избранного
func RejectDeleteFavorite(t *testing.T, token, id string, status int) response.Error {
	t.Helper()
	t.Logf("[Шаг] Проверяем отказ в удалении %s из избранного", id)
	return response.ReadJSON[response.Error](t, client.HttpDeleteFavorite(t, token, id), status)
}
