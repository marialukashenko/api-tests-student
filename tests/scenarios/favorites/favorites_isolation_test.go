package favorites

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/favorites"
	"api-tests-template/tests/fixtures"
)

func TestCannotDeleteAnotherUsersFavorite(t *testing.T) {
	owner, viewer := fixtures.LoginPair(t)
	ad := fixtures.CreateAdvertisement(t, owner.Token, fixtures.NewCreateRequest(t))
	added := fixtures.AddFavorite(t, viewer.Token, ad.ID)

	t.Log("[Подготовка] Объявление есть в избранном второго пользователя, но не автора")
	initial := favorites.GetFavorites(t, owner.Token, ad.Title)
	require.Zero(t, initial.Total, "У автора неожиданно есть тестовая запись избранного")
	err := favorites.RejectDeleteFavorite(t, owner.Token, ad.ID, http.StatusNotFound)
	assertions.APIError(t, err, "not_found", "Favorite entry not found")

	list := favorites.GetFavorites(t, viewer.Token, ad.Title)
	require.Equal(t, 1, list.Total, "Запрос другого пользователя изменил чужое избранное")
	require.Len(t, list.Items, 1, "Чужая запись избранного была удалена")
	require.Equal(t, added, list.Items[0], "Чужая запись избранного изменилась")
	t.Log("[Результат] Пользователь не может удалить запись из чужого избранного")
}
