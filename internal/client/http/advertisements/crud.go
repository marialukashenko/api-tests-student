package advertisements

import (
	"fmt"
	"net/http"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/advertisements/models"
)

// HttpPostAdvertisement отправляет multipart-запрос с полями объявления и фотографиями
func HttpPostAdvertisement(t *testing.T, token string, data models.CreateRequest) *http.Response {
	t.Helper()
	req := apiRunner.GetRunner().Auth(token).Create().Post(path.AdvertisementPath)
	for name, value := range data.Fields {
		req.MultipartFormData(name, value)
	}
	if len(data.Photos) > 0 {
		req.MultipartFile("photos", data.Photos...)
	}
	return req.Expect(t).End().Response
}

// HttpGetAdvertisement получает объявление по идентификатору
func HttpGetAdvertisement(t *testing.T, token, id string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Get(path.AdvertisementPath).
		Query("id", id).Expect(t).End().Response
}

// HttpDeleteAdvertisement удаляет объявление пользователя-владельца
func HttpDeleteAdvertisement(t *testing.T, token, id string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Delete(path.AdvertisementPath).
		Query("id", id).Expect(t).End().Response
}

// HttpGetAdvertisementPhotos получает фотографии объявления по его идентификатору
func HttpGetAdvertisementPhotos(t *testing.T, token, id string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Get(fmt.Sprintf(path.AdvertisementPhotos, id)).
		Expect(t).End().Response
}
