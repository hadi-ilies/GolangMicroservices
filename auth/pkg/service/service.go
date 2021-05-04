package service

import (
	http1 "net/http"
)

// AuthService describes the service.
type AuthService interface {
	// Add your methods here
	// e.x: Foo(ctx context.Context,s string)(rs string, err error)
	//get JWT
	GetJWT() (token string, err error)
	IsAuthorized(r *http1.Request) (token string, err error)
}
