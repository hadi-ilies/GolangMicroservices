package http

import (
	"context"
	"encoding/json"
	"errors"
	endpoint "golangmicroservices/auths/pkg/endpoint"
	http1 "net/http"

	http "github.com/go-kit/kit/transport/http"
	handlers "github.com/gorilla/handlers"
	mux "github.com/gorilla/mux"
)

// makeGetJWTHandler creates the handler logic
func makeGetJWTHandler(m *mux.Router, endpoints endpoint.Endpoints, options []http.ServerOption) {
	m.Methods("POST").Path("/get-jwt").Handler(handlers.CORS(handlers.AllowedMethods([]string{"POST"}), handlers.AllowedOrigins([]string{"*"}))(http.NewServer(endpoints.GetJWTEndpoint, decodeGetJWTRequest, encodeGetJWTResponse, options...)))
}

// decodeGetJWTRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeGetJWTRequest(_ context.Context, r *http1.Request) (interface{}, error) {
	req := endpoint.GetJWTRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// encodeGetJWTResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer
func encodeGetJWTResponse(ctx context.Context, w http1.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(response)
	return
}

// makeCreateAuthHandler creates the handler logic
func makeCreateAuthHandler(m *mux.Router, endpoints endpoint.Endpoints, options []http.ServerOption) {
	m.Methods("POST").Path("/create-auth").Handler(handlers.CORS(handlers.AllowedMethods([]string{"POST"}), handlers.AllowedOrigins([]string{"*"}))(http.NewServer(endpoints.CreateAuthEndpoint, decodeCreateAuthRequest, encodeCreateAuthResponse, options...)))
}

// decodeCreateAuthRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeCreateAuthRequest(_ context.Context, r *http1.Request) (interface{}, error) {
	req := endpoint.CreateAuthRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// encodeCreateAuthResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer
func encodeCreateAuthResponse(ctx context.Context, w http1.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(response)
	return
}

// makeDeleteAuthHandler creates the handler logic
func makeDeleteAuthHandler(m *mux.Router, endpoints endpoint.Endpoints, options []http.ServerOption) {
	m.Methods("POST").Path("/delete-auth").Handler(handlers.CORS(handlers.AllowedMethods([]string{"POST"}), handlers.AllowedOrigins([]string{"*"}))(http.NewServer(endpoints.DeleteAuthEndpoint, decodeDeleteAuthRequest, encodeDeleteAuthResponse, options...)))
}

// decodeDeleteAuthRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeDeleteAuthRequest(_ context.Context, r *http1.Request) (interface{}, error) {
	req := endpoint.DeleteAuthRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// encodeDeleteAuthResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer
func encodeDeleteAuthResponse(ctx context.Context, w http1.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(response)
	return
}

// makeExtractTokenMetadataHandler creates the handler logic
func makeExtractTokenMetadataHandler(m *mux.Router, endpoints endpoint.Endpoints, options []http.ServerOption) {
	m.Methods("POST").Path("/extract-token-metadata").Handler(handlers.CORS(handlers.AllowedMethods([]string{"POST"}), handlers.AllowedOrigins([]string{"*"}))(http.NewServer(endpoints.ExtractTokenMetadataEndpoint, decodeExtractTokenMetadataRequest, encodeExtractTokenMetadataResponse, options...)))
}

// decodeExtractTokenMetadataRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeExtractTokenMetadataRequest(_ context.Context, r *http1.Request) (interface{}, error) {
	req := endpoint.ExtractTokenMetadataRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	req.R = *r
	return req, err
}

// encodeExtractTokenMetadataResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer
func encodeExtractTokenMetadataResponse(ctx context.Context, w http1.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(response)
	return
}

// makeFetchAuthHandler creates the handler logic
func makeFetchAuthHandler(m *mux.Router, endpoints endpoint.Endpoints, options []http.ServerOption) {
	m.Methods("POST").Path("/fetch-auth").Handler(handlers.CORS(handlers.AllowedMethods([]string{"POST"}), handlers.AllowedOrigins([]string{"*"}))(http.NewServer(endpoints.FetchAuthEndpoint, decodeFetchAuthRequest, encodeFetchAuthResponse, options...)))
}

// decodeFetchAuthRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeFetchAuthRequest(_ context.Context, r *http1.Request) (interface{}, error) {
	req := endpoint.FetchAuthRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// encodeFetchAuthResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer
func encodeFetchAuthResponse(ctx context.Context, w http1.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(response)
	return
}
func ErrorEncoder(_ context.Context, err error, w http1.ResponseWriter) {
	w.WriteHeader(err2code(err))
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}
func ErrorDecoder(r *http1.Response) error {
	var w errorWrapper
	if err := json.NewDecoder(r.Body).Decode(&w); err != nil {
		return err
	}
	return errors.New(w.Error)
}

// This is used to set the http status, see an example here :
// https://github.com/go-kit/kit/blob/master/examples/addsvc/pkg/addtransport/http.go#L133
func err2code(err error) int {
	return http1.StatusInternalServerError
}

type errorWrapper struct {
	Error string `json:"error"`
}
