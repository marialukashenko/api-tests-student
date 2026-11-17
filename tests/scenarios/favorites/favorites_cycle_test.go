package favorites

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/favorites"
	"api-tests-template/tests/fixtures"
)

func TestFavoritesFullCycle(t *testing.T) {
	owner, viewer := fixtures.LoginPair(t)
	ad := fixtures.CreateAdvertisement(t, owner.Token, fixtures.NewCreateRequest(t))

	initial := favorites.GetFavorites(t, viewer.Token, ad.Title)
	require.Zero(t, initial.Total, "Новое объявление уже присутствует в избранном")
	added := fixtures.AddFavorite(t, viewer.Token, ad.ID)
	require.NotEmpty(t, added.ID, "В ответе отсутствует идентификатор записи избранного")
	require.NotEmpty(t, added.AddedAt, "В ответе отсутствует дата добавления в избранное")
	require.Equal(t, ad, added.Advertisement, "В избранное добавлено другое объявление")

	list := favorites.GetFavorites(t, viewer.Token, ad.Title)
	require.Equal(t, 1, list.Total, "Объявление не появилось в избранном")
	require.Len(t, list.Items, 1, "Ожидалась одна запись избранного")
	require.Equal(t, added, list.Items[0], "GET /favorites вернул другую запись")
	t.Log("[Результат] Добавленное объявление присутствует в избранном")

	favorites.DeleteFavorite(t, viewer.Token, ad.ID)
	list = favorites.GetFavorites(t, viewer.Token, ad.Title)
	require.Zero(t, list.Total, "Удалённая запись осталась в избранном")
	require.Empty(t, list.Items, "Удалённое из избранного объявление возвращается в items")
	t.Log("[Результат] Удалённое из избранного объявление отсутствует в списке")

	fixtures.AddFavorite(t, viewer.Token, ad.ID)
	list = favorites.GetFavorites(t, viewer.Token, ad.Title)
	require.Equal(t, 1, list.Total, "Повторное добавление после удаления не восстановило запись")
	advertisements.DeleteAdvertisement(t, owner.Token, ad.ID)
	advertisements.RequireDeleted(t, owner.Token, ad.ID)
	list = favorites.GetFavorites(t, viewer.Token, ad.Title)
	require.Zero(t, list.Total, "Удалённое объявление осталось в избранном")
	require.Empty(t, list.Items, "Удаление объявления не очистило избранное")
	t.Log("[Результат] Удаление объявления также удаляет его из избранного")
}
