package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Orendev/shortener/internal/auth"
	"github.com/Orendev/shortener/internal/logger"
	"github.com/Orendev/shortener/internal/models"
	"github.com/Orendev/shortener/internal/random"
	"github.com/Orendev/shortener/internal/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// PostAPIShorten save the link and return the short link.
func (h *Handler) PostAPIShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userID, err := auth.GetAuthIdentifier(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.ShortLinkRequest
	dec := json.NewDecoder(r.Body)
	// читаем тело запроса и декодируем
	if err = dec.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	code := random.Strn(8)
	shortLink := &models.ShortLink{
		UUID:        uuid.New().String(),
		UserID:      userID,
		Code:        code,
		OriginalURL: req.URL,
		ShortURL:    fmt.Sprintf("%s/%s", strings.TrimPrefix(h.baseURL, "/"), code),
		DeletedFlag: false,
	}

	// Сохраним модель
	err = h.repo.Save(r.Context(), *shortLink)

	if err != nil && !errors.Is(err, repository.ErrConflict) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if errors.Is(err, repository.ErrConflict) {
		w.WriteHeader(http.StatusConflict)
		shortLink, err = h.repo.GetByOriginalURL(r.Context(), req.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	// заполняем модель ответа
	resp := models.ShortLinkResponse{
		Result: shortLink.ShortURL,
	}

	enc, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	_, err = w.Write(enc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// PostAPIShortenBatch save the link and return the short link.
func (h *Handler) PostAPIShortenBatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var reqData []models.ShortLinkBatchRequest
	dec := json.NewDecoder(r.Body)
	// читаем тело запроса и декодируем
	if err := dec.Decode(&reqData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	shortLinksInsert := make([]models.ShortLink, 0, len(reqData))
	shortLinksUpdate := make([]models.ShortLink, 0, len(reqData))
	shortLinkBatchResponse := make([]models.ShortLinkBatchResponse, 0, len(reqData))

	userID, err := auth.GetAuthIdentifier(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	for _, req := range reqData {
		code := random.Strn(8)
		var model *models.ShortLink

		model, err = h.repo.GetByID(r.Context(), req.CorrelationID)

		if err != nil {
			model = &models.ShortLink{
				UUID:        req.CorrelationID,
				UserID:      userID,
				Code:        code,
				OriginalURL: req.OriginalURL,
				ShortURL:    fmt.Sprintf("%s/%s", strings.TrimPrefix(h.baseURL, "/"), code),
				DeletedFlag: false,
			}

			shortLinksInsert = append(shortLinksInsert, *model)

		} else {

			model.OriginalURL = req.OriginalURL
			model.DeletedFlag = false
			shortLinksUpdate = append(shortLinksUpdate, *model)
		}

		// заполняем модель ответа
		shortLinkBatchResponse = append(shortLinkBatchResponse, models.ShortLinkBatchResponse{
			CorrelationID: model.UUID,
			ShortURL:      model.ShortURL,
		})
	}

	// Сохраним модель
	err = h.repo.InsertBatch(r.Context(), shortLinksInsert)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.repo.UpdateBatch(r.Context(), shortLinksUpdate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// заполняем модель ответа
	enc, err := json.Marshal(shortLinkBatchResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	_, err = w.Write(enc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// DeleteAPIUserUrls delete the user's link.
func (h *Handler) DeleteAPIUserUrls(w http.ResponseWriter, r *http.Request) {
	// Проверим HTTP Method
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Проверим если у пользователя права
	userID, err := auth.GetAuthIdentifier(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Получим тело запроса
	var reqData []string
	dec := json.NewDecoder(r.Body)
	// читаем тело запроса и декодируем
	if err := dec.Decode(&reqData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	go h.deleteUserUrlsCodes(context.Background(), reqData, userID)
	// Запрос получен, но еще не обработан
	w.WriteHeader(http.StatusAccepted)

}

// GetAPIUserUrls we will get the user's links.
func (h *Handler) GetAPIUserUrls(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	limit := 100
	w.Header().Set("Content-Type", "application/json")
	shortLinkUserResponse := make([]models.ShortLinkUserResponse, 0, limit)

	userID, err := auth.GetAuthIdentifier(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	shortLinks, err := h.repo.ShortLinksByUserID(r.Context(), userID, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(shortLinks) == 0 {
		w.WriteHeader(http.StatusNoContent)
	}

	for _, model := range shortLinks {
		// заполняем модель ответа
		shortLinkUserResponse = append(shortLinkUserResponse, models.ShortLinkUserResponse{
			OriginalURL: model.OriginalURL,
			ShortURL:    model.ShortURL,
		})
	}

	// заполняем модель ответа
	enc, err := json.Marshal(shortLinkUserResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(enc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (h *Handler) deleteUserUrlsCodes(_ context.Context, codes []string, userID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	logger.Log.Info("codes", zap.Any("codes", codes))
	channels := h.fanOut(ctx, func(input []string) chan string {
		inputCh := make(chan string)

		go func() {
			defer close(inputCh)
			for _, data := range input {
				inputCh <- data
			}
		}()

		return inputCh
	}(codes))

	resultCh := h.fanIn(ctx, channels...)

	var urls []string
	for url := range resultCh {
		urls = append(urls, url)
	}

	err := h.repo.DeleteFlagBatch(ctx, urls, userID)
	if err != nil {
		logger.Log.Error("error", zap.Error(err))
	}
}

func (h *Handler) fanOut(ctx context.Context, inputCh chan string) []chan string {
	// количество горутин add
	numWorkers := 10
	// каналы, в которые отправляются результаты
	channels := make([]chan string, numWorkers)

	for i := 0; i < numWorkers; i++ {
		channels = append(channels, h.getShortLink(ctx, inputCh))
	}

	// возвращаем слайс каналов
	return channels
}

// fanIn объединяет несколько каналов resultChs в один.
func (h *Handler) fanIn(ctx context.Context, resultChs ...chan string) chan string {
	// конечный выходной канал в который отправляем данные из всех каналов из слайса, назовём его результирующим
	finalCh := make(chan string)
	// понадобится для ожидания всех горутин
	var wg sync.WaitGroup

	// перебираем все входящие каналы
	for _, ch := range resultChs {
		// в горутину передавать переменную цикла нельзя, поэтому делаем так
		chClosure := ch

		// инкрементируем счётчик горутин, которые нужно подождать
		wg.Add(1)

		go func() {

			// откладываем сообщение о том, что горутина завершилась
			defer wg.Done()

			// получаем данные из канала
			for data := range chClosure {
				select {
				// выходим из горутины, если канал закрылся
				case <-ctx.Done():
					return
				// если не закрылся, отправляем данные в конечный выходной канал
				case finalCh <- data:
				}
			}
		}()
	}

	go func() {
		// ждём завершения всех горутин
		wg.Wait()
		// когда все горутины завершились, закрываем результирующий канал
		close(finalCh)
	}()

	// возвращаем результирующий канал
	return finalCh
}

// getShortLinkCode принимает на вход конткст для прекращения работы и канал с входными данными для работы,
// а возвращает канал, в который будет отправляться результат запроса чтения из БД.
// На фоне будет запущена горутина, выполняющая запрос чтения из БД до момента закрытия doneCh.
func (h *Handler) getShortLink(ctx context.Context, inputCh chan string) chan string {
	// канал с результатом
	resultCh := make(chan string)

	// горутина, в которой добавляем к значению из inputCh единицу и отправляем результат в addRes
	go func() {
		// закрываем канал, когда горутина завершается
		defer close(resultCh)
		// берём из канала inputCh значения, которые надо изменить
		for data := range inputCh {
			/**
			получим модель shortLink по коду
			*/
			model, err := h.repo.GetByCode(ctx, data)
			if err != nil {
				logger.Log.Info("cannot get shortLink", zap.Error(err))
				continue
			}

			select {
			case <-ctx.Done():
				return
			case resultCh <- model.Code:
			}
		}
	}()

	// возвращаем канал для результатов вычислений
	return resultCh

}
