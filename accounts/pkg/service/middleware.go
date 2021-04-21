package service

import (
	"context"
	domain "golangmicroservices/accounts/pkg/domain"

	log "github.com/go-kit/kit/log"
)

// Middleware describes a service middleware.
type Middleware func(AccountsService) AccountsService

type loggingMiddleware struct {
	logger log.Logger
	next   AccountsService
}

// LoggingMiddleware takes a logger as a dependency
// and returns a AccountsService Middleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next AccountsService) AccountsService {
		return &loggingMiddleware{logger, next}
	}

}

func (l loggingMiddleware) SignUp(ctx context.Context, account domain.Account) (d0 domain.Account, e1 error) {
	defer func() {
		l.logger.Log("method", "SignUp", "account", account, "d0", d0, "e1", e1)
	}()
	return l.next.SignUp(ctx, account)
}
func (l loggingMiddleware) SignIn(ctx context.Context, account domain.Auth) (d0 string, e1 error) {
	defer func() {
		l.logger.Log("method", "SignIn", "account", account, "d0", d0, "e1", e1)
	}()
	return l.next.SignIn(ctx, account)
}
func (l loggingMiddleware) Update(ctx context.Context, account domain.Account) (d0 domain.Account, e1 error) {
	defer func() {
		l.logger.Log("method", "Update", "account", account, "d0", d0, "e1", e1)
	}()
	return l.next.Update(ctx, account)
}
func (l loggingMiddleware) Delete(ctx context.Context) (e0 error) {
	defer func() {
		l.logger.Log("method", "Delete", "e0", e0)
	}()
	return l.next.Delete(ctx)
}
func (l loggingMiddleware) Get(ctx context.Context) (d0 []domain.Account, e1 error) {
	defer func() {
		l.logger.Log("method", "Get", "d0", d0, "e1", e1)
	}()
	return l.next.Get(ctx)
}
func (l loggingMiddleware) GetUserInfo(ctx context.Context, username string) (d0 domain.Account, e1 error) {
	defer func() {
		l.logger.Log("method", "GetUserInfo", "username", username, "d0", d0, "e1", e1)
	}()
	return l.next.GetUserInfo(ctx, username)
}
func (l loggingMiddleware) AddFunds(ctx context.Context, funds uint64) (d0 domain.Account, e1 error) {
	defer func() {
		l.logger.Log("method", "AddFunds", "funds", funds, "d0", d0, "e1", e1)
	}()
	return l.next.AddFunds(ctx, funds)
}

func (l loggingMiddleware) Me(ctx context.Context, token string) (d0 domain.Account, e1 error) {
	defer func() {
		l.logger.Log("method", "Me", "d0", d0, "e1", e1)
	}()
	return l.next.Me(ctx, token)
}
