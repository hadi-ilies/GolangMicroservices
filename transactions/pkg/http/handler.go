package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	AccountDomain "golangmicroservices/accounts/pkg/domain"
	accountEndpoint "golangmicroservices/accounts/pkg/endpoint"

	endpoint "golangmicroservices/transactions/pkg/endpoint"
	"io/ioutil"
	http1 "net/http"
	"os"
	"time"

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
	m.Methods("POST").Path("/").Handler(handlers.CORS(handlers.AllowedMethods([]string{"POST"}), handlers.AllowedOrigins([]string{"*"}))(http.NewServer(endpoints.CreateEndpoint, decodeCreateRequest, encodeCreateResponse, options...)))
}

// decodeCreateRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeCreateRequest(_ context.Context, r *http1.Request) (interface{}, error) {
	token, err := IsAuthorized(r)

	if err != nil {
		return nil, err
	}
	myAccount, err := getAccount(token)
	if err != nil {
		return nil, err
	}
	req := endpoint.CreateRequest{}
	req.Transaction.AccountID = myAccount.Id
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

// makeAcceptHandler creates the handler logic
func makeAcceptHandler(m *mux.Router, endpoints endpoint.Endpoints, options []http.ServerOption) {
	m.Methods("POST").Path("/accept").Handler(handlers.CORS(handlers.AllowedMethods([]string{"POST"}), handlers.AllowedOrigins([]string{"*"}))(http.NewServer(endpoints.AcceptEndpoint, decodeAcceptRequest, encodeAcceptResponse, options...)))
}

// decodeAcceptRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeAcceptRequest(_ context.Context, r *http1.Request) (interface{}, error) {
	token, err := IsAuthorized(r)

	if err != nil {
		return nil, err
	}
	myAccount, err := getAccount(token)
	req := endpoint.AcceptRequest{}
	req.Transaction.AccountID = myAccount.Id
	err = json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// encodeAcceptResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer
func encodeAcceptResponse(ctx context.Context, w http1.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(response)
	return
}

// makeRejectHandler creates the handler logic
func makeRejectHandler(m *mux.Router, endpoints endpoint.Endpoints, options []http.ServerOption) {
	m.Methods("POST").Path("/reject").Handler(handlers.CORS(handlers.AllowedMethods([]string{"POST"}), handlers.AllowedOrigins([]string{"*"}))(http.NewServer(endpoints.RejectEndpoint, decodeRejectRequest, encodeRejectResponse, options...)))
}

// decodeRejectRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeRejectRequest(_ context.Context, r *http1.Request) (interface{}, error) {
	token, err := IsAuthorized(r)

	if err != nil {
		return nil, err
	}
	myAccount, err := getAccount(token)
	req := endpoint.RejectRequest{}
	req.Transaction.AccountID = myAccount.Id
	err = json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// encodeRejectResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer
func encodeRejectResponse(ctx context.Context, w http1.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(response)
	return
}

// makeGetAllHandler creates the handler logic
func makeGetAllHandler(m *mux.Router, endpoints endpoint.Endpoints, options []http.ServerOption) {
	m.Methods("GET").Path("/get-all").Handler(handlers.CORS(handlers.AllowedMethods([]string{"GET"}), handlers.AllowedOrigins([]string{"*"}))(http.NewServer(endpoints.GetAllEndpoint, decodeGetAllRequest, encodeGetAllResponse, options...)))
}

// decodeGetAllRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeGetAllRequest(_ context.Context, r *http1.Request) (interface{}, error) {
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
	req := endpoint.GetAllRequest{AccountID: myAccount.Id.Hex()}
	return req, err
}

// encodeGetAllResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer
func encodeGetAllResponse(ctx context.Context, w http1.ResponseWriter, response interface{}) (err error) {
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
