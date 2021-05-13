package endpoint

import (
	"context"
	domain "golangmicroservices/transactions/pkg/domain"
	service "golangmicroservices/transactions/pkg/service"

	endpoint "github.com/go-kit/kit/endpoint"
)

// CreateRequest collects the request parameters for the Create method.
type CreateRequest struct {
	Transaction domain.Transaction `json:"transaction"`
}

// CreateResponse collects the response parameters for the Create method.
type CreateResponse struct {
	D0 domain.Transaction `json:"d0"`
	E1 error              `json:"e1"`
}

// MakeCreateEndpoint returns an endpoint that invokes Create on the service.
func MakeCreateEndpoint(s service.TransactionsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)
		d0, e1 := s.Create(ctx, req.Transaction)
		return CreateResponse{
			D0: d0,
			E1: e1,
		}, nil
	}
}

// Failed implements Failer.
func (r CreateResponse) Failed() error {
	return r.E1
}

// AcceptRequest collects the request parameters for the Accept method.
type AcceptRequest struct {
	Transaction domain.Transaction `json:"transaction"`
}

// AcceptResponse collects the response parameters for the Accept method.
type AcceptResponse struct {
	D0 domain.Transaction `json:"d0"`
	E1 error              `json:"e1"`
}

// MakeAcceptEndpoint returns an endpoint that invokes Accept on the service.
func MakeAcceptEndpoint(s service.TransactionsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AcceptRequest)
		d0, e1 := s.Accept(ctx, req.Transaction)
		return AcceptResponse{
			D0: d0,
			E1: e1,
		}, nil
	}
}

// Failed implements Failer.
func (r AcceptResponse) Failed() error {
	return r.E1
}

// RejectRequest collects the request parameters for the Reject method.
type RejectRequest struct {
	Transaction domain.Transaction `json:"transaction"`
}

// RejectResponse collects the response parameters for the Reject method.
type RejectResponse struct {
	E0 error `json:"e0"`
}

// MakeRejectEndpoint returns an endpoint that invokes Reject on the service.
func MakeRejectEndpoint(s service.TransactionsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(RejectRequest)
		e0 := s.Reject(ctx, req.Transaction)
		return RejectResponse{E0: e0}, nil
	}
}

// Failed implements Failer.
func (r RejectResponse) Failed() error {
	return r.E0
}

// GetAllRequest collects the request parameters for the GetAll method.
type GetAllRequest struct {
	AccountID string
}

// GetAllResponse collects the response parameters for the GetAll method.
type GetAllResponse struct {
	D0 []domain.Transaction `json:"d0"`
	E1 error                `json:"e1"`
}

// MakeGetAllEndpoint returns an endpoint that invokes GetAll on the service.
func MakeGetAllEndpoint(s service.TransactionsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetAllRequest)
		d0, e1 := s.GetAll(ctx, req.AccountID)
		return GetAllResponse{
			D0: d0,
			E1: e1,
		}, nil
	}
}

// Failed implements Failer.
func (r GetAllResponse) Failed() error {
	return r.E1
}

// Failure is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so they've
// failed, and if so encode them using a separate write path based on the error.
type Failure interface {
	Failed() error
}

// Create implements Service. Primarily useful in a client.
func (e Endpoints) Create(ctx context.Context, transaction domain.Transaction) (d0 domain.Transaction, e1 error) {
	request := CreateRequest{Transaction: transaction}
	response, err := e.CreateEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(CreateResponse).D0, response.(CreateResponse).E1
}

// Accept implements Service. Primarily useful in a client.
func (e Endpoints) Accept(ctx context.Context, transaction domain.Transaction) (d0 domain.Transaction, e1 error) {
	request := AcceptRequest{Transaction: transaction}
	response, err := e.AcceptEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(AcceptResponse).D0, response.(AcceptResponse).E1
}

// Reject implements Service. Primarily useful in a client.
func (e Endpoints) Reject(ctx context.Context, transaction domain.Transaction) (e0 error) {
	request := RejectRequest{Transaction: transaction}
	response, err := e.RejectEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(RejectResponse).E0
}

// GetAll implements Service. Primarily useful in a client.
func (e Endpoints) GetAll(ctx context.Context, accountID string) (d0 []domain.Transaction, e1 error) {
	request := GetAllRequest{AccountID: accountID}
	response, err := e.GetAllEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(GetAllResponse).D0, response.(GetAllResponse).E1
}
