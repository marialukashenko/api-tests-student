package createAdvertisement

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/tests/fixtures"
)

func TestCreateAdvertisementAllParameters(t *testing.T) {
	owner := fixtures.Login(t)
	for _, count := range []int{1, 3} {
		t.Run(photoCountName(count), func(t *testing.T) {
			data := fixtures.NewCreateRequest(t)
			for len(data.Photos) < count {
				data.Photos = append(data.Photos, data.Photos[0])
			}
			ad := fixtures.CreateAdvertisement(t, owner.Token, data)
			require.Equal(t, data.Fields["title"], ad.Title, "Название было изменено при создании")
			require.Equal(t, fixtures.Description, ad.Description, "Описание было изменено при создании")
			require.Equal(t, fixtures.PriceValue, ad.Price, "Цена была изменена при создании")
			require.Equal(t, fixtures.QuantityValue, ad.Quantity, "Количество было изменено при создании")
			require.Equal(t, owner.User.Id, ad.UserID, "Созданное объявление создано другим пользователем")
			require.NotEmpty(t, ad.CreatedAt, "Отсутствует время создания")
			require.NotEmpty(t, ad.UpdatedAt, "Отсутствует время обновления")

			got := advertisements.GetAdvertisement(t, owner.Token, ad.ID)
			require.Equal(t, ad, got, "GET /advertisement вернул другие характеристики объявления")
			t.Log("[Результат] Все характеристики созданного объявления совпадают")

			t.Log("[Шаг] Ищем объявление по полному названию")
			list := advertisements.GetAdvertisementsList(t, owner.Token, ad.Title)
			require.Len(t, list.Items, 1, "По уникальному полному названию в выводе должно быть минимум одно объявление")
			require.Equal(t, ad.ID, list.Items[0].ID, "Поиск вернул другое объявление")
			require.Equal(t, ad.Title, list.Items[0].Title, "Поиск вернул другое название")
			t.Log("[Результат] Созданное объявление доступно в поиске")

			photos := advertisements.GetAdvertisementPhotos(t, owner.Token, ad.ID)
			require.Len(t, photos, count, "Количество загруженных фотографий не сходится")
			require.Equal(t, ad.Photos, photos, "Список фотографий отличается от ответа создания")
			for index, photo := range photos {
				require.NotEmpty(t, photo.ID, "У фотографии нет идентификатора")
				require.NotEmpty(t, photo.CreatedAt, "У фотографии нет времени создания")
				require.Equal(t, index, photo.SortOrder, "Неверный порядок фотографий")
				require.Regexp(t, `^https://.+`, photo.URL, "Нет URL загруженной фотографии")
			}
			t.Log("[Результат] Фотографии загружены и доступны через API")
		})
	}
}

// photoCountName даёт понятное имя проверке минимального и максимального числа фотографий
func photoCountName(count int) string {
	if count == 1 {
		return "одна фотография"
	}
	return "три фотографии"
}
