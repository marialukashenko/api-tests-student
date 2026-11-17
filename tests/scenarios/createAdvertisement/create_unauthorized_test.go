package createAdvertisement

import (
	"net/http"
	"testing"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/utils"
	"api-tests-template/tests/fixtures"
)

func TestCreateAdvertisementUnauthorized(t *testing.T) {
	utils.LoadEnv()
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
			t.Log("[Подготовка] Используем корректное объявление и невалидную авторизацию")
			err := advertisements.RejectCreateAdvertisement(t, tc.token, fixtures.NewCreateRequest(t), http.StatusUnauthorized)
			assertions.APIError(t, err, fixtures.Unauthorized, tc.message)
		})
	}
}
