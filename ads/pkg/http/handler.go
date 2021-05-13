package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	AccountDomain "golangmicroservices/accounts/pkg/domain"
	endpoint "golangmicroservices/ads/pkg/endpoint"
	"io/ioutil"
	http1 "net/http"
	"os"
	"time"

	accountEndpoint "golangmicroservices/accounts/pkg/endpoint"

	"github.com/dgrijalva/jwt-go"
	http "github.com/go-kit/kit/transport/http"
	handlers "github.com/gorilla/handlers"
	mux "github.com/gorilla/mux"
)

func getAccount(token string) (*AccountDomain.Account, error) {
	fmt.Println("am I here ?")
	url := "http://accounts:8081/me"
	spaceClient := http1.Client{
		Timeout: time.Second * 20, // Timeout after 2 seconds
	}
	req, err := http1.NewRequest(http1.MethodGet, url, nil)
	if err != nil {
		fmt.Println("CALAMAR1")
		return nil, err
	}
	req.Header.Set("Access-Control-Allow-Origin", "*")
	req.Header.Set("Authorization", token)
	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		fmt.Println("CALAMAR2")
		return nil, getErr
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		fmt.Println("CALAMAR3")
		return nil, readErr
	}
	fmt.Println("BODY = ", string(body))
	var myAccountResponse accountEndpoint.MeResponse = accountEndpoint.MeResponse{}
	jsonErr := json.Unmarshal(body, &myAccountResponse)
	if jsonErr != nil {
		fmt.Println("CALAMAR4")
		return nil, jsonErr
	}

	fmt.Println("USERNAME = ", myAccountResponse.D0.Username)
	return &myAccountResponse.D0, nil
}

//TODO move this inside auth microservice
func IsAuthorized(r *http1.Request) (string, error) {
	if r.Header["Authorization"] != nil {
		token, err := jwt.Parse(r.Header["Authorization"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf(("Invalid Signing Method"))
			}
			aud := "billing.jwtgo.io"
			checkAudience := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
			if !checkAudience {
				return nil, fmt.Errorf(("invalid aud"))
			}
			// verify iss claim
			iss := "jwtgo.io"
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				return nil, fmt.Errorf(("invalid iss"))
			}
			return []byte(os.Getenv("SECRET_KEY")), nil
		})
		if err != nil {
			return "", err
		}

		if token.Valid {
			return r.Header["Authorization"][0], nil
		}

	}
	return "", fmt.Errorf(("no Token detected"))
}

// makeCreateHandler creates the handler logic
func makeCreateHandler(m *mux.Router, endpoints endpoint.Endpoints, options []http.ServerOption) {
	m.Methods("POST").Path("/createAd").Handler(handlers.CORS(handlers.AllowedMethods([]string{"POST"}), handlers.AllowedOrigins([]string{"*"}))(http.NewServer(endpoints.CreateEndpoint, decodeCreateRequest, encodeCreateResponse, options...)))
}

// decodeCreateRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeCreateRequest(_ context.Context, r *http1.Request) (interface{}, error) {
	token, err := IsAuthorized(r)

	if err != nil {
		return nil, err
	}
	//TODO make request to account microservice
	myAccount, err := getAccount(token)
	if err != nil {
		return nil, err
	}
	fmt.Println("myAccount = ", myAccount.Id)
	req := endpoint.CreateRequest{}
	req.Ad.AccountID = myAccount.Id
	err = json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// encodeCreateResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer
func encodeCreateResponse(ctx context.Context, w http1.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(response)
	return
}

// makeUpdateHandler creates the handler logic
func makeUpdateHandler(m *mux.Router, endpoints endpoint.Endpoints, options []http.ServerOption) {
	m.Methods("PUT", "OPTIONS").Path("/").Handler(
		handlers.CORS(
			handlers.AllowedHeaders([]string{"Content-Type", "Content-Length"}),
			handlers.AllowedMethods([]string{"PUT"}),
			handlers.AllowedOrigins([]string{"*"}),
		)(http.NewServer(endpoints.UpdateEndpoint, decodeUpdateRequest, encodeUpdateResponse, options...)))

}

// decodeUpdateRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeUpdateRequest(_ context.Context, r *http1.Request) (interface{}, error) {
	token, err := IsAuthorized(r)

	if err != nil {
		return nil, err
	}
	//TODO make request to account microservice
	myAccount, err := getAccount(token)
	if err != nil {
		return nil, err
	}
	fmt.Println("myAccount = ", myAccount.Id)
	req := endpoint.UpdateRequest{}
	req.Ad.AccountID = myAccount.Id
	err = json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// encodeUpdateResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer
func encodeUpdateResponse(ctx context.Context, w http1.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(response)
	return
}

// makeDeleteHandler creates the handler logic
func makeDeleteHandler(m *mux.Router, endpoints endpoint.Endpoints, options []http.ServerOption) {
	m.Methods("DELETE").Path("/").Handler(handlers.CORS(handlers.AllowedMethods([]string{"DELETE"}), handlers.AllowedOrigins([]string{"*"}))(http.NewServer(endpoints.DeleteEndpoint, decodeDeleteRequest, encodeDeleteResponse, options...)))
}

// decodeDeleteRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeDeleteRequest(_ context.Context, r *http1.Request) (interface{}, error) {
	token, err := IsAuthorized(r)

	if err != nil {
		return nil, err
	}
	//TODO make request to account microservice
	myAccount, err := getAccount(token)
	if err != nil {
		return nil, err
	}
	fmt.Println("myAccount = ", myAccount.Id)
	req := endpoint.DeleteRequest{}
	req.Ad.AccountID = myAccount.Id
	err = json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// encodeDeleteResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer
func encodeDeleteResponse(ctx context.Context, w http1.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(response)
	return
}

// makeGetHandler creates the handler logic
func makeGetHandler(m *mux.Router, endpoints endpoint.Endpoints, options []http.ServerOption) {
	m.Methods("GET").Path("/").Handler(handlers.CORS(handlers.AllowedMethods([]string{"GET"}), handlers.AllowedOrigins([]string{"*"}))(http.NewServer(endpoints.GetEndpoint, decodeGetRequest, encodeGetResponse, options...)))
}

// decodeGetRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeGetRequest(_ context.Context, r *http1.Request) (interface{}, error) {
	_, err := IsAuthorized(r)

	if err != nil {
		return nil, err
	}
	req := endpoint.GetRequest{}
	return req, err
}

// encodeGetResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer
func encodeGetResponse(ctx context.Context, w http1.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(response)
	return
}

// makeGetAllByKeyWordHandler creates the handler logic
func makeGetAllByKeyWordHandler(m *mux.Router, endpoints endpoint.Endpoints, options []http.ServerOption) {
	m.Methods("POST").Path("/get-all-by-key-word").Handler(handlers.CORS(handlers.AllowedMethods([]string{"POST"}), handlers.AllowedOrigins([]string{"*"}))(http.NewServer(endpoints.GetAllByKeyWordEndpoint, decodeGetAllByKeyWordRequest, encodeGetAllByKeyWordResponse, options...)))
}

// decodeGetAllByKeyWordRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeGetAllByKeyWordRequest(_ context.Context, r *http1.Request) (interface{}, error) {
	req := endpoint.GetAllByKeyWordRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// encodeGetAllByKeyWordResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer
func encodeGetAllByKeyWordResponse(ctx context.Context, w http1.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(response)
	return
}

// makeGetAllByUserHandler creates the handler logic
func makeGetAllByUserHandler(m *mux.Router, endpoints endpoint.Endpoints, options []http.ServerOption) {
	m.Methods("POST").Path("/get-all-by-user").Handler(handlers.CORS(handlers.AllowedMethods([]string{"POST"}), handlers.AllowedOrigins([]string{"*"}))(http.NewServer(endpoints.GetAllByUserEndpoint, decodeGetAllByUserRequest, encodeGetAllByUserResponse, options...)))
}

// decodeGetAllByUserRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeGetAllByUserRequest(_ context.Context, r *http1.Request) (interface{}, error) {
	_, err := IsAuthorized(r)

	if err != nil {
		return nil, err
	}
	req := endpoint.GetAllByUserRequest{}
	err = json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// encodeGetAllByUserResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer
func encodeGetAllByUserResponse(ctx context.Context, w http1.ResponseWriter, response interface{}) (err error) {
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
