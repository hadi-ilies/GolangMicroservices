package service

import (
	"context"
	"golangmicroservices/ads/pkg/domain"
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

type basicAdsService struct{}

func (b *basicAdsService) Create(ctx context.Context, ad domain.Ad) (d0 domain.Ad, e1 error) {
	// TODO implement the business logic of Create
	return d0, e1
}
func (b *basicAdsService) Update(ctx context.Context, ad domain.Ad) (d0 domain.Ad, e1 error) {
	// TODO implement the business logic of Update
	return d0, e1
}
func (b *basicAdsService) Delete(ctx context.Context, ad domain.Ad) (d0 domain.Ad, e1 error) {
	// TODO implement the business logic of Delete
	return d0, e1
}
func (b *basicAdsService) Get(ctx context.Context, ad domain.Ad) (d0 domain.Ad, e1 error) {
	// TODO implement the business logic of Get
	return d0, e1
}
func (b *basicAdsService) GetAllByKeyWord(ctx context.Context, keywords string) (d0 []domain.Ad, e1 error) {
	// TODO implement the business logic of GetAllByKeyWord
	return d0, e1
}
func (b *basicAdsService) GetAllByUser(ctx context.Context, username string) (d0 []domain.Ad, e1 error) {
	// TODO implement the business logic of GetAllByUser
	return d0, e1
}

// NewBasicAdsService returns a naive, stateless implementation of AdsService.
func NewBasicAdsService() AdsService {
	return &basicAdsService{}
}

// New returns a AdsService with all of the expected middleware wired in.
func New(middleware []Middleware) AdsService {
	var svc AdsService = NewBasicAdsService()
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}
