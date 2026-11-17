package advertisements

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	client "api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/response"
)

// CreateAdvertisement создаёт объявление и возвращает его параметры
func CreateAdvertisement(t *testing.T, token string, data models.CreateRequest) models.Advertisement {
	t.Helper()
	t.Logf("[Шаг] Создаём объявление %q", data.Fields["title"])
	return response.ReadJSON[models.Advertisement](t, client.HttpPostAdvertisement(t, token, data), http.StatusCreated)
}

// RejectCreateAdvertisement возвращает ошибку создания и удаляет объявление при успехе
func RejectCreateAdvertisement(t *testing.T, token string, data models.CreateRequest, status int) response.Error {
	t.Helper()
	t.Log("[Шаг] Проверяем отказ в создании объявления")
	raw := client.HttpPostAdvertisement(t, token, data)
	if raw.StatusCode == http.StatusCreated {
		ad := response.ReadJSON[models.Advertisement](t, raw, http.StatusCreated)
		if ad.ID != "" {
			CleanupAdvertisement(t, token, ad.ID)
		}
		require.FailNow(t, "Некорректное объявление было создано", "Ожидался HTTP %d", status)
	}
	return response.ReadJSON[response.Error](t, raw, status)
}

// GetAdvertisement возвращает параметры объявления по идентификатору
func GetAdvertisement(t *testing.T, token, id string) models.Advertisement {
	t.Helper()
	t.Logf("[Шаг] Получаем объявление %s", id)
	return response.ReadJSON[models.Advertisement](t, client.HttpGetAdvertisement(t, token, id), http.StatusOK)
}

// DeleteAdvertisement удаляет объявление с обязательным статусом 204
func DeleteAdvertisement(t *testing.T, token, id string) {
	t.Helper()
	t.Logf("[Шаг] Удаляем объявление %s", id)
	response.RequireStatus(t, client.HttpDeleteAdvertisement(t, token, id), http.StatusNoContent)
}

// CleanupAdvertisement удаляет тестовый ресурс и проверяет результат
func CleanupAdvertisement(t *testing.T, token, id string) {
	t.Helper()
	t.Logf("[Очистка] Удаляем тестовое объявление %s", id)
	raw := response.RetryCleanup(t, func() *http.Response {
		return client.HttpDeleteAdvertisement(t, token, id)
	})
	if raw.StatusCode == http.StatusForbidden {
		response.ReadJSON[response.Error](t, raw, http.StatusForbidden)
	} else {
		response.RequireStatus(t, raw, http.StatusNoContent, http.StatusNotFound)
	}
	verification := response.RetryCleanup(t, func() *http.Response {
		return client.HttpGetAdvertisement(t, token, id)
	})
	response.RequireStatus(t, verification, http.StatusNotFound)
	t.Logf("[Результат] Тестовое объявление %s удалено, GET вернул 404", id)
}

// RequireDeleted проверяет, что удалённое объявление недоступно
func RequireDeleted(t *testing.T, token, id string) {
	t.Helper()
	response.RequireStatus(t, client.HttpGetAdvertisement(t, token, id), http.StatusNotFound)
}

// GetAdvertisementPhotos возвращает типизированный список загруженных фотографий
func GetAdvertisementPhotos(t *testing.T, token, id string) []models.Photo {
	t.Helper()
	t.Logf("[Шаг] Получаем фотографии объявления %s", id)
	return response.ReadJSON[[]models.Photo](t, client.HttpGetAdvertisementPhotos(t, token, id), http.StatusOK)
}
