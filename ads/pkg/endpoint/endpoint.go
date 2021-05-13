package endpoint

import (
	"context"
	domain "golangmicroservices/ads/pkg/domain"
	service "golangmicroservices/ads/pkg/service"

	endpoint "github.com/go-kit/kit/endpoint"
)

// CreateRequest collects the request parameters for the Create method.
type CreateRequest struct {
	Ad domain.Ad `json:"ad"`
}

// CreateResponse collects the response parameters for the Create method.
type CreateResponse struct {
	D0 domain.Ad `json:"d0"`
	E1 error     `json:"e1"`
}

// MakeCreateEndpoint returns an endpoint that invokes Create on the service.
func MakeCreateEndpoint(s service.AdsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)
		d0, e1 := s.Create(ctx, req.Ad)
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

// UpdateRequest collects the request parameters for the Update method.
type UpdateRequest struct {
	Ad domain.Ad `json:"ad"`
}

// UpdateResponse collects the response parameters for the Update method.
type UpdateResponse struct {
	D0 domain.Ad `json:"d0"`
	E1 error     `json:"e1"`
}

// MakeUpdateEndpoint returns an endpoint that invokes Update on the service.
func MakeUpdateEndpoint(s service.AdsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateRequest)
		d0, e1 := s.Update(ctx, req.Ad)
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
	Ad domain.Ad `json:"ad"`
}

// DeleteResponse collects the response parameters for the Delete method.
type DeleteResponse struct {
	E1 error `json:"e1"`
}

// MakeDeleteEndpoint returns an endpoint that invokes Delete on the service.
func MakeDeleteEndpoint(s service.AdsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteRequest)
		e1 := s.Delete(ctx, req.Ad)
		return DeleteResponse{
			E1: e1,
		}, nil
	}
}

// Failed implements Failer.
func (r DeleteResponse) Failed() error {
	return r.E1
}

// GetRequest collects the request parameters for the Get method.
type GetRequest struct {
}

// GetResponse collects the response parameters for the Get method.
type GetResponse struct {
	D0 []domain.Ad `json:"d0"`
	E1 error       `json:"e1"`
}

// MakeGetEndpoint returns an endpoint that invokes Get on the service.
func MakeGetEndpoint(s service.AdsService) endpoint.Endpoint {
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

// GetAllByKeyWordRequest collects the request parameters for the GetAllByKeyWord method.
type GetAllByKeyWordRequest struct {
	Keywords string `json:"keywords"`
}

// GetAllByKeyWordResponse collects the response parameters for the GetAllByKeyWord method.
type GetAllByKeyWordResponse struct {
	D0 []domain.Ad `json:"d0"`
	E1 error       `json:"e1"`
}

// MakeGetAllByKeyWordEndpoint returns an endpoint that invokes GetAllByKeyWord on the service.
func MakeGetAllByKeyWordEndpoint(s service.AdsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetAllByKeyWordRequest)
		d0, e1 := s.GetAllByKeyWord(ctx, req.Keywords)
		return GetAllByKeyWordResponse{
			D0: d0,
			E1: e1,
		}, nil
	}
}

// Failed implements Failer.
func (r GetAllByKeyWordResponse) Failed() error {
	return r.E1
}

// GetAllByUserRequest collects the request parameters for the GetAllByUser method.
type GetAllByUserRequest struct {
	TargetAccountID string `json:"account_id"`
}

// GetAllByUserResponse collects the response parameters for the GetAllByUser method.
type GetAllByUserResponse struct {
	D0 []domain.Ad `json:"d0"`
	E1 error       `json:"e1"`
}

// MakeGetAllByUserEndpoint returns an endpoint that invokes GetAllByUser on the service.
func MakeGetAllByUserEndpoint(s service.AdsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetAllByUserRequest)
		d0, e1 := s.GetAllByUser(ctx, req.TargetAccountID)
		return GetAllByUserResponse{
			D0: d0,
			E1: e1,
		}, nil
	}
}

// Failed implements Failer.
func (r GetAllByUserResponse) Failed() error {
	return r.E1
}

// Failure is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so they've
// failed, and if so encode them using a separate write path based on the error.
type Failure interface {
	Failed() error
}

// Create implements Service. Primarily useful in a client.
func (e Endpoints) Create(ctx context.Context, ad domain.Ad) (d0 domain.Ad, e1 error) {
	request := CreateRequest{Ad: ad}
	response, err := e.CreateEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(CreateResponse).D0, response.(CreateResponse).E1
}

// Update implements Service. Primarily useful in a client.
func (e Endpoints) Update(ctx context.Context, ad domain.Ad) (d0 domain.Ad, e1 error) {
	request := UpdateRequest{Ad: ad}
	response, err := e.UpdateEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(UpdateResponse).D0, response.(UpdateResponse).E1
}

// Delete implements Service. Primarily useful in a client.
func (e Endpoints) Delete(ctx context.Context, ad domain.Ad) (e1 error) {
	request := DeleteRequest{Ad: ad}
	response, err := e.DeleteEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(DeleteResponse).E1
}

// Get implements Service. Primarily useful in a client.
func (e Endpoints) Get(ctx context.Context, ad domain.Ad) (d0 []domain.Ad, e1 error) {
	request := GetRequest{}
	response, err := e.GetEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(GetResponse).D0, response.(GetResponse).E1
}

// GetAllByKeyWord implements Service. Primarily useful in a client.
func (e Endpoints) GetAllByKeyWord(ctx context.Context, keywords string) (d0 []domain.Ad, e1 error) {
	request := GetAllByKeyWordRequest{Keywords: keywords}
	response, err := e.GetAllByKeyWordEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(GetAllByKeyWordResponse).D0, response.(GetAllByKeyWordResponse).E1
}

// GetAllByUser implements Service. Primarily useful in a client.
func (e Endpoints) GetAllByUser(ctx context.Context, targetAccountID string) (d0 []domain.Ad, e1 error) {
	request := GetAllByUserRequest{TargetAccountID: targetAccountID}
	response, err := e.GetAllByUserEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(GetAllByUserResponse).D0, response.(GetAllByUserResponse).E1
}
