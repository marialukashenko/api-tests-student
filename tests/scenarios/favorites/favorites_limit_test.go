package favorites

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/favorites"
	"api-tests-template/tests/fixtures"
)

const favoritesLimit = 50

func TestFavoritesLimit(t *testing.T) {
	owner, viewer := fixtures.LoginPair(t)
	initial := favorites.GetFavorites(t, viewer.Token, "")
	require.LessOrEqual(t, initial.Total, favoritesLimit, "Исходное избранное превышает лимит")
	t.Cleanup(func() {
		restored := favorites.GetFavorites(t, viewer.Token, "")
		require.Equal(t, initial, restored, "После очистки изменилось исходное избранное пользователя")
		t.Log("[Результат] Исходное избранное восстановлено после проверки лимита")
	})

	t.Logf("[Подготовка] Заполняем %d свободных мест избранного", favoritesLimit-initial.Total)
	for i := initial.Total; i < favoritesLimit; i++ {
		ad := fixtures.CreateAdvertisement(t, owner.Token, fixtures.NewCreateRequest(t))
		favorites.AddFavorite(t, viewer.Token, ad.ID)
	}
	full := favorites.GetFavorites(t, viewer.Token, "")
	require.Equal(t, favoritesLimit, full.Total, "Не удалось подготовить избранное на границе лимита")

	ad := fixtures.CreateAdvertisement(t, owner.Token, fixtures.NewCreateRequest(t))
	err := favorites.RejectAddFavorite(t, viewer.Token, ad.ID, http.StatusBadRequest)
	assertions.APIError(t, err, "limit_reached", "Favorites limit of 50 reached")
	list := favorites.GetFavorites(t, viewer.Token, ad.Title)
	require.Zero(t, list.Total, "Объявление добавилось сверх лимита")
	require.Empty(t, list.Items, "Объявление сверх лимита присутствует в избранном")
	after := favorites.GetFavorites(t, viewer.Token, "")
	require.Equal(t, favoritesLimit, after.Total, "Отказ изменил общее число записей избранного")
	t.Log("[Результат] Добавление сверх лимита 50 записей отклонено")
}
