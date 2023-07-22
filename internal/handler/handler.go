package handler

import (
	"context"
	"sync"
	"time"

	"github.com/Orendev/shortener/internal/logger"
	"github.com/Orendev/shortener/internal/models"
	"github.com/Orendev/shortener/internal/repository"
	"go.uber.org/zap"
)

type Handler struct {
	repo    repository.Storage
	baseURL string
}

// NewHandler конструктор создает структуру Handler
func NewHandler(repo repository.Storage, baseURL string) Handler {
	return Handler{repo: repo, baseURL: baseURL}
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
		channels[i] = h.getShortLink(ctx, inputCh, userID)
	}

	// возвращаем слайс каналов
	return channels
}

// getShortLinkCode принимает на вход конткст для прекращения работы и канал с входными данными для работы,
// а возвращает канал, в который будет отправляться результат запроса чтения из БД.
// На фоне будет запущена горутина, выполняющая запрос чтения из БД до момента закрытия doneCh.
func (h *Handler) getShortLink(ctx context.Context, inputCh chan string, _ string) chan models.ShortLink {
	// канал с результатом
	resultCh := make(chan models.ShortLink)

	// горутина, в которой добавляем к значению из inputCh единицу и отправляем результат в addRes
	go func() {
		// закрываем канал, когда горутина завершается
		defer close(resultCh)
		// берём из канала inputCh значения, которые надо изменить
		for data := range inputCh {

			model, err := h.repo.GetByCode(ctx, data)
			if err != nil {
				logger.Log.Info("cannot get shortLink", zap.Error(err))
				continue
			}

			//if model.UserID == userID {
			//	model.DeletedFlag = true
			//}
			model.DeletedFlag = true

			select {
			case <-ctx.Done():
				return
			case resultCh <- *model:
			}
		}
	}()

	// возвращаем канал для результатов вычислений
	return resultCh

}

func (h *Handler) flushShortLink(ctx context.Context, resultCh chan models.ShortLink) {

	var shortLinks []models.ShortLink

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
			err := h.repo.UpdateBatch(ctx, shortLinks)
			if err != nil {
				logger.Log.Info("cannot save shortLink", zap.Error(err))
				continue
			}
			shortLinks = nil
		}
	}

}

// fanIn объединяет несколько каналов resultChs в один.
func (h *Handler) fanIn(ctx context.Context, resultChs ...chan models.ShortLink) chan models.ShortLink {
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
