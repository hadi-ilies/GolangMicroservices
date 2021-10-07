package endpoint

import (
	"context"
	domain "golangmicroservices/accounts/pkg/domain"
	service "golangmicroservices/accounts/pkg/service"

	endpoint "github.com/go-kit/kit/endpoint"
)

// SignUpRequest collects the request parameters for the SignUp method.
type SignUpRequest struct {
	Account domain.Account `json:"account"`
}

// SignUpResponse collects the response parameters for the SignUp method.
type SignUpResponse struct {
	D0 domain.Account `json:"d0"`
	E1 error          `json:"e1"`
}

// MakeSignUpEndpoint returns an endpoint that invokes SignUp on the service.
func MakeSignUpEndpoint(s service.AccountsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SignUpRequest)
		d0, e1 := s.SignUp(ctx, req.Account)
		return SignUpResponse{
			D0: d0,
			E1: e1,
		}, nil
	}
}

// Failed implements Failer.
func (r SignUpResponse) Failed() error {
	return r.E1
}

// SignInRequest collects the request parameters for the SignIn method.
type SignInRequest struct {
	Account domain.Auth `json:"auth"`
}

// SignInResponse collects the response parameters for the SignIn method.
type SignInResponse struct {
	D0 map[string]string `json:"d0"`
	E1 error             `json:"e1"`
}

// MakeSignInEndpoint returns an endpoint that invokes SignIn on the service.
func MakeSignInEndpoint(s service.AccountsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SignInRequest)
		d0, e1 := s.SignIn(ctx, req.Account)
		return SignInResponse{
			D0: d0,
			E1: e1,
		}, nil
	}
}

// Failed implements Failer.
func (r SignInResponse) Failed() error {
	return r.E1
}

// UpdateRequest collects the request parameters for the Update method.
type UpdateRequest struct {
	Account domain.Account `json:"account"`
}

// UpdateResponse collects the response parameters for the Update method.
type UpdateResponse struct {
	D0 domain.Account `json:"d0"`
	E1 error          `json:"e1"`
}

// MakeUpdateEndpoint returns an endpoint that invokes Update on the service.
func MakeUpdateEndpoint(s service.AccountsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateRequest)
		d0, e1 := s.Update(ctx, req.Account)
		return UpdateResponse{
			D0: d0,
			E1: e1,
		}, nil
	}
}

// Failed implements Failer.
func (r UpdateResponse) Failed() error {
	return r.E1
}

// DeleteRequest collects the request parameters for the Delete method.
type DeleteRequest struct {
	Token string `json:"-"`
}

// DeleteResponse collects the response parameters for the Delete method.
type DeleteResponse struct {
	E0 error `json:"e0"`
}

// MakeDeleteEndpoint returns an endpoint that invokes Delete on the service.
func MakeDeleteEndpoint(s service.AccountsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteRequest)
		e0 := s.Delete(ctx, req.Token)
		return DeleteResponse{E0: e0}, nil
	}
}

// Failed implements Failer.
func (r DeleteResponse) Failed() error {
	return r.E0
}

// GetRequest collects the request parameters for the Get method.
type GetRequest struct{}

// GetResponse collects the response parameters for the Get method.
type GetResponse struct {
	D0 []domain.Account `json:"d0"`
	E1 error            `json:"e1"`
}

// MakeGetEndpoint returns an endpoint that invokes Get on the service.
func MakeGetEndpoint(s service.AccountsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		d0, e1 := s.Get(ctx)
		return GetResponse{
			D0: d0,
			E1: e1,
		}, nil
	}
}

// Failed implements Failer.
func (r GetResponse) Failed() error {
	return r.E1
}

// GetUserInfoRequest collects the request parameters for the GetUserInfo method.
type GetUserInfoRequest struct {
	Username string `json:"username"`
}

// GetUserInfoResponse collects the response parameters for the GetUserInfo method.
type GetUserInfoResponse struct {
	D0 domain.Account `json:"d0"`
	E1 error          `json:"e1"`
}

// MakeGetUserInfoEndpoint returns an endpoint that invokes GetUserInfo on the service.
func MakeGetUserInfoEndpoint(s service.AccountsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetUserInfoRequest)
		d0, e1 := s.GetUserInfo(ctx, req.Username)
		return GetUserInfoResponse{
			D0: d0,
			E1: e1,
		}, nil
	}
}

// Failed implements Failer.
func (r GetUserInfoResponse) Failed() error {
	return r.E1
}

// AddFundsRequest collects the request parameters for the AddFunds method.
type AddFundsRequest struct {
	Token string `json:"-"`
	Funds uint64 `json:"funds"`
}

// AddFundsResponse collects the response parameters for the AddFunds method.
type AddFundsResponse struct {
	D0 domain.Account `json:"d0"`
	E1 error          `json:"e1"`
}

// MakeAddFundsEndpoint returns an endpoint that invokes AddFunds on the service.
func MakeAddFundsEndpoint(s service.AccountsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AddFundsRequest)
		d0, e1 := s.AddFunds(ctx, req.Token, req.Funds)
		return AddFundsResponse{
			D0: d0,
			E1: e1,
		}, nil
	}
}

// Failed implements Failer.
func (r AddFundsResponse) Failed() error {
	return r.E1
}

// Failure is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so they've
// failed, and if so encode them using a separate write path based on the error.
type Failure interface {
	Failed() error
}

// SignUp implements Service. Primarily useful in a client.
func (e Endpoints) SignUp(ctx context.Context, account domain.Account) (d0 domain.Account, e1 error) {
	request := SignUpRequest{Account: account}
	response, err := e.SignUpEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(SignUpResponse).D0, response.(SignUpResponse).E1
}

// SignIn implements Service. Primarily useful in a client.
func (e Endpoints) SignIn(ctx context.Context, account domain.Auth) (d0 map[string]string, e1 error) {
	request := SignInRequest{Account: account}
	response, err := e.SignInEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(SignInResponse).D0, response.(SignInResponse).E1
}

// Update implements Service. Primarily useful in a client.
func (e Endpoints) Update(ctx context.Context, account domain.Account) (d0 domain.Account, e1 error) {
	request := UpdateRequest{Account: account}
	response, err := e.UpdateEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(UpdateResponse).D0, response.(UpdateResponse).E1
}

// Delete implements Service. Primarily useful in a client.
func (e Endpoints) Delete(ctx context.Context) (e0 error) {
	request := DeleteRequest{}
	response, err := e.DeleteEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(DeleteResponse).E0
}

// Get implements Service. Primarily useful in a client.
func (e Endpoints) Get(ctx context.Context) (d0 []domain.Account, e1 error) {
	request := GetRequest{}
	response, err := e.GetEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(GetResponse).D0, response.(GetResponse).E1
}

// GetUserInfo implements Service. Primarily useful in a client.
func (e Endpoints) GetUserInfo(ctx context.Context, username string) (d0 domain.Account, e1 error) {
	request := GetUserInfoRequest{Username: username}
	response, err := e.GetUserInfoEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(GetUserInfoResponse).D0, response.(GetUserInfoResponse).E1
}

// AddFunds implements Service. Primarily useful in a client.
func (e Endpoints) AddFunds(ctx context.Context, funds uint64) (d0 domain.Account, e1 error) {
	request := AddFundsRequest{Funds: funds}
	response, err := e.AddFundsEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(AddFundsResponse).D0, response.(AddFundsResponse).E1
}

// MeRequest collects the request parameters for the Me method.
type MeRequest struct {
	Token string `json:"-"`
}

// MeResponse collects the response parameters for the Me method.
type MeResponse struct {
	D0 domain.Account `json:"d0"`
	E1 error          `json:"e1"`
}

// MakeMeEndpoint returns an endpoint that invokes Me on the service.
func MakeMeEndpoint(s service.AccountsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(MeRequest)
		d0, e1 := s.Me(ctx, req.Token)
		return MeResponse{
			D0: d0,
			E1: e1,
		}, nil
	}
}

// Failed implements Failer.
func (r MeResponse) Failed() error {
	return r.E1
}

// Me implements Service. Primarily useful in a client.
func (e Endpoints) Me(ctx context.Context) (d0 domain.Account, e1 error) {
	request := MeRequest{}
	response, err := e.MeEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(MeResponse).D0, response.(MeResponse).E1
}

// LogoutRequest collects the request parameters for the Logout method.
type LogoutRequest struct {
	Token string `json:"token"`
}

// LogoutResponse collects the response parameters for the Logout method.
type LogoutResponse struct {
	E0 error `json:"e0"`
}

// MakeLogoutEndpoint returns an endpoint that invokes Logout on the service.
func MakeLogoutEndpoint(s service.AccountsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LogoutRequest)
		e0 := s.Logout(ctx, req.Token)
		return LogoutResponse{E0: e0}, nil
	}
}

// Failed implements Failer.
func (r LogoutResponse) Failed() error {
	return r.E0
}

// Logout implements Service. Primarily useful in a client.
func (e Endpoints) Logout(ctx context.Context, token string) (e0 error) {
	request := LogoutRequest{Token: token}
	response, err := e.LogoutEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(LogoutResponse).E0
}
