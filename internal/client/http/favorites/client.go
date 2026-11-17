package favorites

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/favorites/models"
)

// HttpPostFavorite добавляет объявление в избранное пользователя
func HttpPostFavorite(t *testing.T, token, id string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Post(path.FavoritesPath).
		JSON(models.AddRequest{AdvertisementID: id}).Expect(t).End().Response
}

// HttpGetFavorites получает страницу избранного, отфильтрованную по названию объявления
func HttpGetFavorites(t *testing.T, token, search string, page int) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Get(path.FavoritesPath).
		Query("search", search).Query("page", strconv.Itoa(page)).Expect(t).End().Response
}

// HttpDeleteFavorite удаляет объявление из избранного по его идентификатору
func HttpDeleteFavorite(t *testing.T, token, id string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Delete(fmt.Sprintf(path.FavoritePath, id)).
		Expect(t).End().Response
}