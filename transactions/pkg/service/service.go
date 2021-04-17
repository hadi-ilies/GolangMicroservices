package service

import (
	"context"
	"golangmicroservices/transactions/pkg/domain"
)

// TransactionsService describes the service.
type TransactionsService interface {
	// Add your methods here
	// e.x: Foo(ctx context.Context,s string)(rs string, err error)
	// - Make an offer on an ad
	Create(ctx context.Context, transaction domain.Transaction) (domain.Transaction, error)
	// - Accept an offer on its own ad
	Accept(ctx context.Context, transaction domain.Transaction) (domain.Transaction, error)
	// - Refuse an offer on its own ad
	Reject(ctx context.Context, transaction domain.Transaction) error
	// - List all its own transaction
	GetAll(ctx context.Context) ([]domain.Transaction, error)
}

type basicTransactionsService struct{}

func (b *basicTransactionsService) Create(ctx context.Context, transaction domain.Transaction) (d0 domain.Transaction, e1 error) {
	// TODO implement the business logic of Create
	return d0, e1
}
func (b *basicTransactionsService) Accept(ctx context.Context, transaction domain.Transaction) (d0 domain.Transaction, e1 error) {
	// TODO implement the business logic of Accept
	return d0, e1
}
func (b *basicTransactionsService) Reject(ctx context.Context, transaction domain.Transaction) (e0 error) {
	// TODO implement the business logic of Reject
	return e0
}
func (b *basicTransactionsService) GetAll(ctx context.Context) (d0 []domain.Transaction, e1 error) {
	// TODO implement the business logic of GetAll
	return d0, e1
}

// NewBasicTransactionsService returns a naive, stateless implementation of TransactionsService.
func NewBasicTransactionsService() TransactionsService {
	return &basicTransactionsService{}
}

// New returns a TransactionsService with all of the expected middleware wired in.
func New(middleware []Middleware) TransactionsService {
	var svc TransactionsService = NewBasicTransactionsService()
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}
