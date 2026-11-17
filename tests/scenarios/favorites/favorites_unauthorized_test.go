package favorites

import (
	"net/http"
	"testing"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/favorites"
	"api-tests-template/tests/fixtures"
)

func TestFavoritesUnauthorized(t *testing.T) {
	owner, viewer := fixtures.LoginPair(t)
	ad := fixtures.CreateAdvertisement(t, owner.Token, fixtures.NewCreateRequest(t))
	fixtures.AddFavorite(t, viewer.Token, ad.ID)
	cases := []struct {
		name    string
		token   string
		message string
	}{
		{"без токена", "", fixtures.MissingAuthMessage},
		{"неверный токен", fixtures.InvalidToken, fixtures.InvalidAuthMessage},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("GET", func(t *testing.T) {
				err := favorites.RejectGetFavorites(t, tc.token, http.StatusUnauthorized)
				assertions.APIError(t, err, fixtures.Unauthorized, tc.message)
			})
			t.Run("POST", func(t *testing.T) {
				err := favorites.RejectAddFavorite(t, tc.token, ad.ID, http.StatusUnauthorized)
				assertions.APIError(t, err, fixtures.Unauthorized, tc.message)
			})
			t.Run("DELETE", func(t *testing.T) {
				err := favorites.RejectDeleteFavorite(t, tc.token, ad.ID, http.StatusUnauthorized)
				assertions.APIError(t, err, fixtures.Unauthorized, tc.message)
			})
		})
	}
}
