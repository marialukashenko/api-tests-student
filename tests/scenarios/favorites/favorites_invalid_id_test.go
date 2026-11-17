package favorites

import (
	"net/http"
	"testing"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/favorites"
	"api-tests-template/tests/fixtures"
)

func TestFavoritesInvalidAdvertisementID(t *testing.T) {
	user := fixtures.Login(t)
	for _, id := range []string{"", "not-a-uuid"} {
		t.Run("добавление_"+id, func(t *testing.T) {
			t.Log("[Подготовка] Используем некорректный идентификатор объявления")
			err := favorites.RejectAddFavorite(t, user.Token, id, http.StatusBadRequest)
			message := "invalid UUID length: 10"
			if id == "" {
				message = "invalid UUID length: 0"
			}
			assertions.APIError(t, err, "invalid_request", message)
		})
	}
}
