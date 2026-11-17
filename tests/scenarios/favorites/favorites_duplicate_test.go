package favorites

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/favorites"
	"api-tests-template/tests/fixtures"
)

func TestCannotFavoriteAdvertisementTwice(t *testing.T) {
	owner, viewer := fixtures.LoginPair(t)
	ad := fixtures.CreateAdvertisement(t, owner.Token, fixtures.NewCreateRequest(t))
	original := fixtures.AddFavorite(t, viewer.Token, ad.ID)
	err := favorites.RejectAddFavorite(t, viewer.Token, ad.ID, http.StatusConflict)
	assertions.APIError(t, err, "already_exists", "Advertisement is already in favorites")
	list := favorites.GetFavorites(t, viewer.Token, ad.Title)
	require.Equal(t, 1, list.Total, "Повторное добавление изменило число записей избранного")
	require.Len(t, list.Items, 1, "В избранном появилась дублирующая запись")
	require.Equal(t, original.ID, list.Items[0].ID, "Повторное добавление заменило исходную запись")
	t.Log("[Результат] Повторное добавление отклонено без изменения избранного")
}
