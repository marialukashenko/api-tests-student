package favorites

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/favorites"
	"api-tests-template/tests/fixtures"
)

func TestCannotFavoriteOwnAdvertisement(t *testing.T) {
	owner := fixtures.Login(t)
	ad := fixtures.CreateAdvertisement(t, owner.Token, fixtures.NewCreateRequest(t))
	err := favorites.RejectAddFavorite(t, owner.Token, ad.ID, http.StatusForbidden)
	assertions.APIError(t, err, "forbidden", "Cannot add own advertisement to favorites")
	list := favorites.GetFavorites(t, owner.Token, ad.Title)
	require.Zero(t, list.Total, "Своё объявление появилось в избранном после отказа")
	require.Empty(t, list.Items, "Своё объявление присутствует в items после отказа")
}
