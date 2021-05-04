package service

import (
	"context"
	domain "golangmicroservices/ads/pkg/domain"

	log "github.com/go-kit/kit/log"
)

// Middleware describes a service middleware.
type Middleware func(AdsService) AdsService

type loggingMiddleware struct {
	logger log.Logger
	next   AdsService
}

// LoggingMiddleware takes a logger as a dependency
// and returns a AdsService Middleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next AdsService) AdsService {
		return &loggingMiddleware{logger, next}
	}

}

func (l loggingMiddleware) Create(ctx context.Context, ad domain.Ad) (d0 domain.Ad, e1 error) {
	defer func() {
		l.logger.Log("method", "Create", "ad", ad, "d0", d0, "e1", e1)
	}()
	return l.next.Create(ctx, ad)
}
func (l loggingMiddleware) Update(ctx context.Context, ad domain.Ad) (d0 domain.Ad, e1 error) {
	defer func() {
		l.logger.Log("method", "Update", "ad", ad, "d0", d0, "e1", e1)
	}()
	return l.next.Update(ctx, ad)
}
func (l loggingMiddleware) Delete(ctx context.Context, ad domain.Ad) (e1 error) {
	defer func() {
		l.logger.Log("method", "Delete", "ad", ad, "e1", e1)
	}()
	return l.next.Delete(ctx, ad)
}
func (l loggingMiddleware) Get(ctx context.Context, ad domain.Ad) (d0 domain.Ad, e1 error) {
	defer func() {
		l.logger.Log("method", "Get", "ad", ad, "d0", d0, "e1", e1)
	}()
	return l.next.Get(ctx, ad)
}
func (l loggingMiddleware) GetAllByKeyWord(ctx context.Context, keywords string) (d0 []domain.Ad, e1 error) {
	defer func() {
		l.logger.Log("method", "GetAllByKeyWord", "keywords", keywords, "d0", d0, "e1", e1)
	}()
	return l.next.GetAllByKeyWord(ctx, keywords)
}
func (l loggingMiddleware) GetAllByUser(ctx context.Context, username string) (d0 []domain.Ad, e1 error) {
	defer func() {
		l.logger.Log("method", "GetAllByUser", "username", username, "d0", d0, "e1", e1)
	}()
	return l.next.GetAllByUser(ctx, username)
}
