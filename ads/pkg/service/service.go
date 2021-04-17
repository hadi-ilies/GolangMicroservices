package service

import (
	"context"
	"golangmicroservices/ads/domain"

	"gitlab.com/CamillePolice/golangmicroservices/ads/pkg/domain"
)

// AdsService describes the service.
type AdsService interface {
	// Add your methods here
	// e.x: Foo(ctx context.Context,s string)(rs string, err error)
	// - Create an ad
	Create(ctx context.Context, ad domain.Ad) (domain.Ad, error)
	// - Update one of its own ad
	Update(ctx context.Context, ad domain.Ad) (domain.Ad, error)
	// - Delete one of its own ad
	Delete(ctx context.Context, ad domain.Ad) (domain.Ad, error)
	// - Read any ad
	Get(ctx context.Context, ad domain.Ad) (domain.Ad, error)
	// - Get a list of ads searching by keywords
	GetAllByKeyWord(ctx context.Context, keywords string) ([]domain.Ad, error)
	// - Get a list of all the ads of a user
	GetAllByUser(ctx context.Context, username string) ([]domain.Ad, error)
}
