package service

import (
	"context"
	"golangmicroservices/transactions/pkg/domain"

	"gitlab.com/CamillePolice/golangmicroservices/transactions/pkg/domain"
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
