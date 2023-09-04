package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Orendev/shortener/internal/auth"
	"github.com/Orendev/shortener/internal/dedupe"
	"github.com/Orendev/shortener/internal/logger"
	"github.com/Orendev/shortener/internal/models"
	"github.com/Orendev/shortener/internal/random"
	"github.com/Orendev/shortener/internal/repository"
	"github.com/Orendev/shortener/internal/utils"
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

	go func() {
		for _, code := range reqData {
			h.msgDeleteUserUrlsChan <- models.Message{
				UserID: userID,
				Code:   code,
			}
		}
	}()

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

// GetAPIStats statistics on the short link service.
func (h Handler) GetAPIStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ip, err := utils.ResolveIP(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	check, err := utils.CidrRangeContains(h.trustedSubnet, ip.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !check {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	urls, err := h.repo.UrlsStats(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	users, err := h.repo.UsersStats(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	stats := models.StatsResponse{
		Urls:  urls,
		Users: users,
	}

	// заполняем модель ответа
	enc, err := json.Marshal(stats)
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

func (h *Handler) flushDeleteShortLink() {

	messages := make(map[string][]string)

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case message := <-h.msgDeleteUserUrlsChan:
			messages[message.UserID] = dedupe.DedupeStrings(append(messages[message.UserID], message.Code))
		case <-ticker.C:
			if len(messages) == 0 {
				continue
			}

			for userID, codes := range messages {
				err := h.repo.DeleteFlagBatch(context.Background(), codes, userID)
				if err != nil {
					logger.Log.Error("cannot delete shortLink", zap.Error(err))
					continue
				}
			}

			// почистим срез сообщений на удаления
			for k := range messages {
				delete(messages, k)
			}
		}
	}

}
