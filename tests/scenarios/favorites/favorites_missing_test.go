package favorites

import (
	"net/http"
	"testing"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/favorites"
	"api-tests-template/tests/fixtures"
)

func TestFavoritesNonexistentAdvertisement(t *testing.T) {
	user := fixtures.Login(t)
	t.Run("добавление несуществующего объявления", func(t *testing.T) {
		t.Log("[Подготовка] Используем отсутствующий UUID объявления")
		err := favorites.RejectAddFavorite(t, user.Token, fixtures.MissingID, http.StatusNotFound)
		assertions.APIError(t, err, "not_found", "Advertisement not found")
	})
	t.Run("удаление несуществующего объявления", func(t *testing.T) {
		t.Log("[Подготовка] Используем отсутствующий UUID объявления")
		err := favorites.RejectDeleteFavorite(t, user.Token, fixtures.MissingID, http.StatusNotFound)
		assertions.APIError(t, err, "not_found", "Favorite entry not found")
	})
}
