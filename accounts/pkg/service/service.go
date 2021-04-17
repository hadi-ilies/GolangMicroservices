package service

import (
	"context"
	"golangmicroservices/accounts/pkg/domain"

	"gitlab.com/CamillePolice/golangmicroservices/accounts/pkg/domain"
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
