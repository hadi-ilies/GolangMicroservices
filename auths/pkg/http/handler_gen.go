// THIS FILE IS AUTO GENERATED BY GK-CLI DO NOT EDIT!!
package http

import (
	http "github.com/go-kit/kit/transport/http"
	mux "github.com/gorilla/mux"
	endpoint "golangmicroservices/auths/pkg/endpoint"
	http1 "net/http"
)

// NewHTTPHandler returns a handler that makes a set of endpoints available on
// predefined paths.
func NewHTTPHandler(endpoints endpoint.Endpoints, options map[string][]http.ServerOption) http1.Handler {
	m := mux.NewRouter()
	makeGetJWTHandler(m, endpoints, options["GetJWT"])
	makeCreateAuthHandler(m, endpoints, options["CreateAuth"])
	makeDeleteAuthHandler(m, endpoints, options["DeleteAuth"])
	makeExtractTokenMetadataHandler(m, endpoints, options["ExtractTokenMetadata"])
	makeFetchAuthHandler(m, endpoints, options["FetchAuth"])
	return m
}
