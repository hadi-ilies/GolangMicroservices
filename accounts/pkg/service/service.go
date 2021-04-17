package service

import (
	"context"
	"golangmicroservices/accounts/pkg/domain"
)

// AccountsService describes the service.
type AccountsService interface {
	// Add your methods here
	// e.x: Foo(ctx context.Context,s string)(rs string, err error)
	//Create an account
	SignUp(ctx context.Context, account domain.Account) (domain.Account, error)
	//Login
	SignIn(ctx context.Context, account domain.Account) (domain.Account, error)
	//Update informations of its own account
	Update(ctx context.Context, account domain.Account) (domain.Account, error)
	//Delete its own account
	Delete(ctx context.Context) error
	//Fully read its own account
	Get(ctx context.Context) (domain.Account, error)
	//Partially read any user account
	GetUserInfo(ctx context.Context, username string) (domain.Account, error)
	//Add funds to it's own balance
	AddFunds(ctx context.Context, funds uint64) (domain.Account, error)
}

type basicAccountsService struct{}

func (b *basicAccountsService) SignUp(ctx context.Context, account domain.Account) (d0 domain.Account, e1 error) {
	// TODO implement the business logic of SignUp
	return d0, e1
}
func (b *basicAccountsService) SignIn(ctx context.Context, account domain.Account) (d0 domain.Account, e1 error) {
	// TODO implement the business logic of SignIn
	return d0, e1
}
func (b *basicAccountsService) Update(ctx context.Context, account domain.Account) (d0 domain.Account, e1 error) {
	// TODO implement the business logic of Update
	return d0, e1
}
func (b *basicAccountsService) Delete(ctx context.Context) (e0 error) {
	// TODO implement the business logic of Delete
	return e0
}
func (b *basicAccountsService) Get(ctx context.Context) (d0 domain.Account, e1 error) {
	// TODO implement the business logic of Get
	return d0, e1
}
func (b *basicAccountsService) GetUserInfo(ctx context.Context, username string) (d0 domain.Account, e1 error) {
	// TODO implement the business logic of GetUserInfo
	return d0, e1
}
func (b *basicAccountsService) AddFunds(ctx context.Context, funds uint64) (d0 domain.Account, e1 error) {
	// TODO implement the business logic of AddFunds
	return d0, e1
}

// NewBasicAccountsService returns a naive, stateless implementation of AccountsService.
func NewBasicAccountsService() AccountsService {
	return &basicAccountsService{}
}

// New returns a AccountsService with all of the expected middleware wired in.
func New(middleware []Middleware) AccountsService {
	var svc AccountsService = NewBasicAccountsService()
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}
