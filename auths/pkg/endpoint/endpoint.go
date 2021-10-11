package endpoint

import (
	"context"
	endpoint "github.com/go-kit/kit/endpoint"
	domain "golangmicroservices/auths/pkg/domain"
	service "golangmicroservices/auths/pkg/service"
	"net/http"
)

// GetJWTRequest collects the request parameters for the GetJWT method.
type GetJWTRequest struct {
	UserID string `json:"user_id"`
}

// GetJWTResponse collects the response parameters for the GetJWT method.
type GetJWTResponse struct {
	Token domain.TokenDetails `json:"token"`
	Err   error               `json:"err"`
}

// MakeGetJWTEndpoint returns an endpoint that invokes GetJWT on the service.
func MakeGetJWTEndpoint(s service.AuthsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetJWTRequest)
		token, err := s.GetJWT(ctx, req.UserID)
		return GetJWTResponse{
			Err:   err,
			Token: token,
		}, nil
	}
}

// Failed implements Failer.
func (r GetJWTResponse) Failed() error {
	return r.Err
}

// CreateAuthRequest collects the request parameters for the CreateAuth method.
type CreateAuthRequest struct {
	UserID string              `json:"user_id"`
	Td     domain.TokenDetails `json:"td"`
}

// CreateAuthResponse collects the response parameters for the CreateAuth method.
type CreateAuthResponse struct {
	Log string `json:"log"`
	Err error  `json:"err"`
}

// MakeCreateAuthEndpoint returns an endpoint that invokes CreateAuth on the service.
func MakeCreateAuthEndpoint(s service.AuthsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateAuthRequest)
		log, err := s.CreateAuth(ctx, req.UserID, req.Td)
		return CreateAuthResponse{
			Err: err,
			Log: log,
		}, nil
	}
}

// Failed implements Failer.
func (r CreateAuthResponse) Failed() error {
	return r.Err
}

// DeleteAuthRequest collects the request parameters for the DeleteAuth method.
type DeleteAuthRequest struct {
	GivenUuid string `json:"given_uuid"`
}

// DeleteAuthResponse collects the response parameters for the DeleteAuth method.
type DeleteAuthResponse struct {
	Log string `json:"log"`
	Err error  `json:"err"`
}

// MakeDeleteAuthEndpoint returns an endpoint that invokes DeleteAuth on the service.
func MakeDeleteAuthEndpoint(s service.AuthsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteAuthRequest)
		log, err := s.DeleteAuth(ctx, req.GivenUuid)
		return DeleteAuthResponse{
			Err: err,
			Log: log,
		}, nil
	}
}

// Failed implements Failer.
func (r DeleteAuthResponse) Failed() error {
	return r.Err
}

// ExtractTokenMetadataRequest collects the request parameters for the ExtractTokenMetadata method.
type ExtractTokenMetadataRequest struct {
	R http.Request `json:"r"`
}

// ExtractTokenMetadataResponse collects the response parameters for the ExtractTokenMetadata method.
type ExtractTokenMetadataResponse struct {
	Details domain.AccessDetails `json:"details"`
	Err     error                `json:"err"`
}

// MakeExtractTokenMetadataEndpoint returns an endpoint that invokes ExtractTokenMetadata on the service.
func MakeExtractTokenMetadataEndpoint(s service.AuthsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ExtractTokenMetadataRequest)
		details, err := s.ExtractTokenMetadata(ctx, req.R)
		return ExtractTokenMetadataResponse{
			Details: details,
			Err:     err,
		}, nil
	}
}

// Failed implements Failer.
func (r ExtractTokenMetadataResponse) Failed() error {
	return r.Err
}

// FetchAuthRequest collects the request parameters for the FetchAuth method.
type FetchAuthRequest struct {
	AuthD domain.AccessDetails `json:"auth_d"`
}

// FetchAuthResponse collects the response parameters for the FetchAuth method.
type FetchAuthResponse struct {
	UserID string `json:"user_id"`
	Err    error  `json:"err"`
}

// MakeFetchAuthEndpoint returns an endpoint that invokes FetchAuth on the service.
func MakeFetchAuthEndpoint(s service.AuthsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(FetchAuthRequest)
		userID, err := s.FetchAuth(ctx, req.AuthD)
		return FetchAuthResponse{
			Err:    err,
			UserID: userID,
		}, nil
	}
}

// Failed implements Failer.
func (r FetchAuthResponse) Failed() error {
	return r.Err
}

// Failure is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so they've
// failed, and if so encode them using a separate write path based on the error.
type Failure interface {
	Failed() error
}

// GetJWT implements Service. Primarily useful in a client.
func (e Endpoints) GetJWT(ctx context.Context, userID string) (token domain.TokenDetails, err error) {
	request := GetJWTRequest{UserID: userID}
	response, err := e.GetJWTEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(GetJWTResponse).Token, response.(GetJWTResponse).Err
}

// CreateAuth implements Service. Primarily useful in a client.
func (e Endpoints) CreateAuth(ctx context.Context, userID string, td domain.TokenDetails) (log string, err error) {
	request := CreateAuthRequest{
		Td:     td,
		UserID: userID,
	}
	response, err := e.CreateAuthEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(CreateAuthResponse).Log, response.(CreateAuthResponse).Err
}

// DeleteAuth implements Service. Primarily useful in a client.
func (e Endpoints) DeleteAuth(ctx context.Context, givenUuid string) (log string, err error) {
	request := DeleteAuthRequest{GivenUuid: givenUuid}
	response, err := e.DeleteAuthEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(DeleteAuthResponse).Log, response.(DeleteAuthResponse).Err
}

// ExtractTokenMetadata implements Service. Primarily useful in a client.
func (e Endpoints) ExtractTokenMetadata(ctx context.Context, r http.Request) (details domain.AccessDetails, err error) {
	request := ExtractTokenMetadataRequest{R: r}
	response, err := e.ExtractTokenMetadataEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(ExtractTokenMetadataResponse).Details, response.(ExtractTokenMetadataResponse).Err
}

// FetchAuth implements Service. Primarily useful in a client.
func (e Endpoints) FetchAuth(ctx context.Context, authD domain.AccessDetails) (userID string, err error) {
	request := FetchAuthRequest{AuthD: authD}
	response, err := e.FetchAuthEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(FetchAuthResponse).UserID, response.(FetchAuthResponse).Err
}
