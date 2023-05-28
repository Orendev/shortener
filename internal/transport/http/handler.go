package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Orendev/shortener/internal/auth"
	"github.com/Orendev/shortener/internal/logger"
	"github.com/Orendev/shortener/internal/models"
	"github.com/Orendev/shortener/internal/random"
	"github.com/Orendev/shortener/internal/storage"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Handler struct {
	shortLinkStorage storage.ShortLinkStorage
	baseURL          string
}

func NewHandler(storage storage.ShortLinkStorage, baseURL string) Handler {
	return Handler{shortLinkStorage: storage, baseURL: baseURL}
}

func (h *Handler) ShortLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	code := strings.TrimPrefix(r.URL.Path, "/")

	shortLink, err := h.shortLinkStorage.GetByCode(r.Context(), code)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if shortLink.DeletedFlag {
		// Целевой запрос больше не доступен
		w.WriteHeader(http.StatusGone)
		return
	}

	w.Header().Add("Location", shortLink.OriginalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) ShortLinkAdd(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req := models.ShortLinkRequest{}
	req.URL = string(body)

	if err = req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := auth.GetAuthIdentifier(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	logger.Log.Info("userID", zap.String("userID", userID))
	code := random.Strn(8)
	shortLink := &models.ShortLink{
		UUID:        uuid.New().String(),
		UserID:      userID,
		Code:        code,
		OriginalURL: req.URL,
		ShortURL:    fmt.Sprintf("%s/%s", strings.TrimPrefix(h.baseURL, "/"), code),
		DeletedFlag: false,
	}

	err = h.shortLinkStorage.Save(r.Context(), *shortLink)

	if err != nil && !errors.Is(err, storage.ErrConflict) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if errors.Is(err, storage.ErrConflict) {
		w.WriteHeader(http.StatusConflict)
		shortLink, err = h.shortLinkStorage.GetByOriginalURL(r.Context(), req.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(shortLink.ShortURL))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
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
	if err := dec.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := req.Validate(); err != nil {
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
	err = h.shortLinkStorage.Save(r.Context(), *shortLink)

	if err != nil && !errors.Is(err, storage.ErrConflict) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if errors.Is(err, storage.ErrConflict) {
		w.WriteHeader(http.StatusConflict)
		shortLink, err = h.shortLinkStorage.GetByOriginalURL(r.Context(), req.URL)
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

func (h *Handler) ShortenBatch(w http.ResponseWriter, r *http.Request) {
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

		model, err := h.shortLinkStorage.GetByID(r.Context(), req.CorrelationID)

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
	err = h.shortLinkStorage.InsertBatch(r.Context(), shortLinksInsert)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.shortLinkStorage.UpdateBatch(r.Context(), shortLinksUpdate)
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

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	err := h.shortLinkStorage.Ping(ctx)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UserUrls(w http.ResponseWriter, r *http.Request) {

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

	shortLinks, err := h.shortLinkStorage.ShortLinksByUserID(r.Context(), userID, limit)
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

func (h *Handler) DeleteUserUrls(w http.ResponseWriter, r *http.Request) {
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
	logger.Log.Info("userID", zap.String("userID", userID))
	logger.Log.Info("reqData", zap.Any("reqData", reqData))
	inputCh := generator(reqData)
	channels := h.fanOut(r.Context(), inputCh, userID)
	resultCh := fanIn(channels...)

	go h.flushShortLink(r.Context(), resultCh)

	// Запрос получен, но еще не обработан
	w.WriteHeader(http.StatusAccepted)

}

// generator создаем каналы
func generator(input []string) chan string {
	inputCh := make(chan string)

	go func() {
		defer close(inputCh)
		for _, data := range input {
			inputCh <- data
		}
	}()

	return inputCh
}

func (h *Handler) fanOut(ctx context.Context, inputCh chan string, userID string) []chan models.ShortLink {
	// количество горутин add
	numWorkers := 10
	// каналы, в которые отправляются результаты
	channels := make([]chan models.ShortLink, numWorkers)

	for i := 0; i < numWorkers; i++ {
		addResultCh := h.getShortLink(ctx, inputCh, userID)
		channels[i] = addResultCh
	}

	// возвращаем слайс каналов
	return channels
}

// getShortLinkCode принимает на вход конткст для прекращения работы и канал с входными данными для работы,
// а возвращает канал, в который будет отправляться результат запроса чтения из БД.
// На фоне будет запущена горутина, выполняющая запрос чтения из БД до момента закрытия doneCh.
func (h *Handler) getShortLink(ctx context.Context, inputCh chan string, userID string) chan models.ShortLink {
	// канал с результатом
	resultCh := make(chan models.ShortLink)

	// горутина, в которой добавляем к значению из inputCh единицу и отправляем результат в addRes
	go func() {
		// закрываем канал, когда горутина завершается
		defer close(resultCh)

		// берём из канала inputCh значения, которые надо изменить
		for data := range inputCh {

			model, err := h.shortLinkStorage.GetByCode(ctx, data)
			if err != nil {
				logger.Log.Debug("cannot get shortLink", zap.Error(err))
				continue
			}

			if model.UserID == userID {
				model.DeletedFlag = true
			}

			resultCh <- *model
		}
	}()

	// возвращаем канал для результатов вычислений
	return resultCh

}

func (h *Handler) flushShortLink(_ context.Context, resultCh chan models.ShortLink) {

	var shortLinks []models.ShortLink

	ctx, cancel := context.WithTimeout(context.Background(), 55*time.Second)
	defer cancel()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case shortLink, ok := <-resultCh:
			if !ok {
				continue
			}
			shortLinks = append(shortLinks, shortLink)
		case <-ticker.C:
			if len(shortLinks) == 0 {
				continue
			}
			err := h.shortLinkStorage.UpdateBatch(ctx, shortLinks)
			if err != nil {
				logger.Log.Info("cannot save shortLink", zap.Error(err))
				continue
			}
			shortLinks = nil
		}
	}

}

// Merge объединяет несколько каналов resultChs в один.
func fanIn(resultChs ...chan models.ShortLink) chan models.ShortLink {
	// конечный выходной канал в который отправляем данные из всех каналов из слайса, назовём его результирующим
	finalCh := make(chan models.ShortLink)
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
				finalCh <- data
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
