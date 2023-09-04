package grpc

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Orendev/shortener/internal/models"
	pb "github.com/Orendev/shortener/internal/pkg/grpc/proto"
	"github.com/Orendev/shortener/internal/random"
	"github.com/Orendev/shortener/internal/repository"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPC struct {
	pb.UnimplementedShortenerServiceServer
	repo          repository.Storage
	baseURL       string
	trustedSubnet string
}

func NewGRPC(repo repository.Storage, baseURL, trustedSubnet string) *GRPC {
	return &GRPC{repo: repo, baseURL: baseURL, trustedSubnet: trustedSubnet}
}

func (g *GRPC) GetAPIUserUrls(ctx context.Context, reg *pb.APIUserUrlsRequest) (*pb.APIUserUrlsResponse, error) {
	var response pb.APIUserUrlsResponse

	shortLinks, err := g.repo.ShortLinksByUserID(ctx, reg.UserID, 100)
	if err != nil {
		return nil, status.Error(codes.NotFound, "shorten url not found")
	}

	if len(shortLinks) == 0 {
		return nil, status.Error(codes.NotFound, "no content")
	}

	userUrls := make([]*pb.UserUrl, len(shortLinks))
	for _, model := range shortLinks {
		userUrls = append(userUrls, &pb.UserUrl{
			OriginalUrl: model.OriginalURL,
			ShortUrl:    model.ShortURL,
		})
	}

	response.UserUrls = userUrls

	return &response, nil
}

func (g *GRPC) GetAPIStats(ctx context.Context, _ *pb.APIStatsRequest) (*pb.APIStatsResponse, error) {

	urls, err := g.repo.UrlsStats(ctx)
	if err != nil {
		return nil, status.Error(codes.NotFound, "not found")
	}
	users, err := g.repo.UsersStats(ctx)
	if err != nil {
		return nil, status.Error(codes.NotFound, "not found")
	}

	response := pb.APIStatsResponse{
		Urls:  int64(urls),
		Users: int64(users),
	}

	return &response, nil
}

func (g *GRPC) SaveAPIShorten(ctx context.Context, reg *pb.APIShortenRequest) (*pb.APIShortenResponse, error) {
	var response pb.APIShortenResponse

	code := random.Strn(8)
	shortLink := &models.ShortLink{
		UUID:        uuid.New().String(),
		UserID:      uuid.New().String(),
		Code:        code,
		OriginalURL: reg.URL,
		ShortURL:    fmt.Sprintf("%s/%s", strings.TrimPrefix(g.baseURL, "/"), code),
		DeletedFlag: false,
	}

	response.Result = shortLink.ShortURL
	// Сохраним модель
	err := g.repo.Save(ctx, *shortLink)

	if err != nil && !errors.Is(err, repository.ErrConflict) {
		return nil, status.Error(codes.Internal, "something went wrong")
	}

	if errors.Is(err, repository.ErrConflict) {
		shortLink, err = g.repo.GetByOriginalURL(ctx, reg.URL)
		if err != nil {
			return nil, status.Error(codes.Internal, "something went wrong")
		}
		return nil, status.Error(codes.AlreadyExists, shortLink.ShortURL)
	}

	return &response, nil
}

func (g *GRPC) Ping(ctx context.Context, reg *pb.PingRequest) (*pb.PingResponse, error) {
	var response pb.PingResponse

	err := g.repo.Ping(ctx)

	if err != nil {
		return nil, status.Error(codes.Internal, "something went wrong")
	}
	response.Result = "Ping"
	return &response, nil
}
