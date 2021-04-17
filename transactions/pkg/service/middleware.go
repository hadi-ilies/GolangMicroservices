package service

import (
	"context"
	log "github.com/go-kit/kit/log"
	domain "golangmicroservices/transactions/pkg/domain"
)

// Middleware describes a service middleware.
type Middleware func(TransactionsService) TransactionsService

type loggingMiddleware struct {
	logger log.Logger
	next   TransactionsService
}

// LoggingMiddleware takes a logger as a dependency
// and returns a TransactionsService Middleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next TransactionsService) TransactionsService {
		return &loggingMiddleware{logger, next}
	}

}

func (l loggingMiddleware) Create(ctx context.Context, transaction domain.Transaction) (d0 domain.Transaction, e1 error) {
	defer func() {
		l.logger.Log("method", "Create", "transaction", transaction, "d0", d0, "e1", e1)
	}()
	return l.next.Create(ctx, transaction)
}
func (l loggingMiddleware) Accept(ctx context.Context, transaction domain.Transaction) (d0 domain.Transaction, e1 error) {
	defer func() {
		l.logger.Log("method", "Accept", "transaction", transaction, "d0", d0, "e1", e1)
	}()
	return l.next.Accept(ctx, transaction)
}
func (l loggingMiddleware) Reject(ctx context.Context, transaction domain.Transaction) (e0 error) {
	defer func() {
		l.logger.Log("method", "Reject", "transaction", transaction, "e0", e0)
	}()
	return l.next.Reject(ctx, transaction)
}
func (l loggingMiddleware) GetAll(ctx context.Context) (d0 []domain.Transaction, e1 error) {
	defer func() {
		l.logger.Log("method", "GetAll", "d0", d0, "e1", e1)
	}()
	return l.next.GetAll(ctx)
}
