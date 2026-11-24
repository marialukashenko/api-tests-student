package fixtures

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
	adModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/managers/favorites"
	favoriteModels "api-tests-template/internal/managers/favorites/models"
	"api-tests-template/internal/utils"
)

const (
	Description        = "тестовое описание объявления"
	Price              = "52.67"
	PriceValue         = 52.67
	Quantity           = "16"
	QuantityValue      = 16
	MissingID          = "00000000-1337-0000-6767-000000000000"
	InvalidToken       = "incorrect_token"
	Unauthorized       = "unauthorized"
	MissingAuthMessage = "Authorization header required"
	InvalidAuthMessage = "Invalid or expired token"
)

// Login загружает настройки и авторизует основной тестовый аккаунт
func Login(t *testing.T) authModels.LoginOkResponse {
	t.Helper()
	utils.LoadEnv()
	t.Log("[Подготовка] Авторизуем основной тестовый аккаунт")
	return auth.Login(t, os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

// LoginPair авторизует автора объявления и отдельного пользователя избранного
func LoginPair(t *testing.T) (authModels.LoginOkResponse, authModels.LoginOkResponse) {
	t.Helper()
	owner := Login(t)
	require.NotEmpty(t, os.Getenv("SECOND_TEST_LOGIN"), "Задайте SECOND_TEST_LOGIN в .env.override")
	require.NotEmpty(t, os.Getenv("SECOND_TEST_PASSWORD"), "Задайте SECOND_TEST_PASSWORD в .env.override")
	t.Log("[Подготовка] Авторизуем второй тестовый аккаунт")
	viewer := auth.Login(t, os.Getenv("SECOND_TEST_LOGIN"), os.Getenv("SECOND_TEST_PASSWORD"))
	require.NotEqual(t, owner.User.Id, viewer.User.Id, "Для проверки избранного нужны два разных аккаунта")
	return owner, viewer
}

// NewCreateRequest формирует уникальное объявление со всеми полями и фотографиями
//
//	func NewCreateRequest(t *testing.T) adModels.CreateRequest {
//		t.Helper()
//		photoPath := filepath.Join(t.TempDir(), "photo.png")
//		file, err := os.Create(photoPath)
//		require.NoError(t, err, "Не удалось создать тестовую фотографию")
//		encodeErr := png.Encode(file, image.NewRGBA(image.Rect(0, 0, 2, 2)))
//		closeErr := file.Close()
//		require.NoError(t, encodeErr, "Не удалось записать PNG")
//		require.NoError(t, closeErr, "Не удалось закрыть PNG")
//		return adModels.CreateRequest{
//			Fields: map[string]string{
//				"title":       "API-test-" + utils.RandomString(16),
//				"description": Description,
//				"price":       Price,
//				"quantity":    Quantity,
//			},
//			Photos: []string{photoPath},
//		}
//	}
//
// NewCreateRequest формирует уникальное объявление
// со всеми полями и тестовой фотографией.
func NewCreateRequest(t *testing.T) adModels.CreateRequest {
	t.Helper()

	_, filename, _, ok := runtime.Caller(0)
	require.True(t, ok, "Не удалось определить путь к фикстуре")

	photoPath := filepath.Join(
		filepath.Dir(filename),
		"..",
		"testdata",
		"photo.png",
	)

	require.FileExists(t, photoPath, "Тестовая фотография не найдена")

	return adModels.CreateRequest{
		Fields: map[string]string{
			"title":       "API-test-" + utils.RandomString(16),
			"description": Description,
			"price":       Price,
			"quantity":    Quantity,
		},
		Photos: []string{photoPath},
	}
}

// CreateAdvertisement регистрирует очистку сразу после создания тестового объявления
func CreateAdvertisement(t *testing.T, token string, data adModels.CreateRequest) adModels.Advertisement {
	t.Helper()
	ad := advertisements.CreateAdvertisement(t, token, data)
	if ad.ID != "" {
		t.Cleanup(func() { advertisements.CleanupAdvertisement(t, token, ad.ID) })
	}
	require.NotEmpty(t, ad.ID, "Сервер не вернул идентификатор созданного объявления")
	return ad
}

// AddFavorite заранее регистрирует очистку тестовой записи избранного
func AddFavorite(t *testing.T, token, id string) favoriteModels.Favorite {
	t.Helper()
	t.Cleanup(func() { favorites.CleanupFavorite(t, token, id) })
	return favorites.AddFavorite(t, token, id)
}
