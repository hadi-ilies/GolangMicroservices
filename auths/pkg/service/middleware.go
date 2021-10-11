package service

import (
	"context"
	log "github.com/go-kit/kit/log"
	domain "golangmicroservices/auths/pkg/domain"
	"net/http"
)

// Middleware describes a service middleware.
type Middleware func(AuthsService) AuthsService

type loggingMiddleware struct {
	logger log.Logger
	next   AuthsService
}

// LoggingMiddleware takes a logger as a dependency
// and returns a AuthsService Middleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next AuthsService) AuthsService {
		return &loggingMiddleware{logger, next}
	}

}

func (l loggingMiddleware) GetJWT(ctx context.Context, userID string) (token domain.TokenDetails, err error) {
	defer func() {
		l.logger.Log("method", "GetJWT", "userID", userID, "token", token, "err", err)
	}()
	return l.next.GetJWT(ctx, userID)
}
func (l loggingMiddleware) CreateAuth(ctx context.Context, userID string, td domain.TokenDetails) (log string, err error) {
	defer func() {
		l.logger.Log("method", "CreateAuth", "userID", userID, "td", td, "log", log, "err", err)
	}()
	return l.next.CreateAuth(ctx, userID, td)
}
func (l loggingMiddleware) DeleteAuth(ctx context.Context, givenUuid string) (log string, err error) {
	defer func() {
		l.logger.Log("method", "DeleteAuth", "givenUuid", givenUuid, "log", log, "err", err)
	}()
	return l.next.DeleteAuth(ctx, givenUuid)
}
func (l loggingMiddleware) ExtractTokenMetadata(ctx context.Context, r http.Request) (details domain.AccessDetails, err error) {
	defer func() {
		l.logger.Log("method", "ExtractTokenMetadata", "r", r, "details", details, "err", err)
	}()
	return l.next.ExtractTokenMetadata(ctx, r)
}
func (l loggingMiddleware) FetchAuth(ctx context.Context, authD domain.AccessDetails) (userID string, err error) {
	defer func() {
		l.logger.Log("method", "FetchAuth", "authD", authD, "userID", userID, "err", err)
	}()
	return l.next.FetchAuth(ctx, authD)
}
