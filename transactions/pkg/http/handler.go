package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	authEndpoint "golangmicroservices/auths/pkg/endpoint"
	endpoint "golangmicroservices/transactions/pkg/endpoint"
	"io/ioutil"
	http1 "net/http"
	"net/url"
	"time"

	http "github.com/go-kit/kit/transport/http"
	handlers "github.com/gorilla/handlers"
	mux "github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

func ExtractTokenMetadata(r http1.Request) (*authEndpoint.ExtractTokenMetadataResponse, error) {
	URL, _ := url.Parse("http://auths:8084/extract-token-metadata")
	r.URL.Scheme = URL.Scheme
	r.URL.Host = URL.Host
	r.URL.Path = "http://auths:8084/extract-token-metadata"
	r.RequestURI = ""
	r.Method = http1.MethodPost

	spaceClient := http1.Client{
		Timeout: time.Second * 20, // Timeout after 2 seconds
	}
	// Step 2: adjust Header
	r.Header.Set("X-Forwarded-For", r.RemoteAddr)
	res, getErr := spaceClient.Do(&r)
	if getErr != nil {
		fmt.Println("CALAMAR2")
		return nil, getErr
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		fmt.Println("CALAMAR3")
		return nil, readErr
	}
	myAccountResponse := authEndpoint.ExtractTokenMetadataResponse{}
	jsonErr := json.Unmarshal(body, &myAccountResponse)
	if jsonErr != nil {
		fmt.Println("CALAMAR4")
		return nil, jsonErr
	}

	fmt.Println("myAccountResponse.Details.UserID = ", myAccountResponse.Details.UserID)
	if myAccountResponse.Details.UserID == "" {
		return nil, fmt.Errorf("Error: auth Request failed")
	}
	return &myAccountResponse, nil
}

//FetchAuth() accepts the AccessDetails from the ExtractTokenMetadata function, then looks it up in mongodb.
//If the record is not found, it may mean the token has expired, hence an error is thrown.
func FetchAuth(authD *authEndpoint.FetchAuthRequest) (string, error) {
	url := "http://auths:8084/fetch-auth"
	spaceClient := http1.Client{
		Timeout: time.Second * 20, // Timeout after 2 seconds
	}
	s, _ := json.Marshal(*authD)
	req, err := http1.NewRequest(http1.MethodPost, url, bytes.NewReader(s))
	if err != nil {
		fmt.Println("CALAMAR1")
		return "", err
	}
	req.Header.Set("Access-Control-Allow-Origin", "*")

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		fmt.Println("CALAMAR2")
		return "", getErr
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		fmt.Println("CALAMAR3")
		return "", readErr
	}
	myAccountResponse := authEndpoint.FetchAuthResponse{}
	jsonErr := json.Unmarshal(body, &myAccountResponse)
	if jsonErr != nil {
		fmt.Println("CALAMAR4")
		return "", jsonErr
	}

	fmt.Println("myAccountResponse.UserID = ", myAccountResponse.UserID)
	return myAccountResponse.UserID, nil
}

func checkToken(r *http1.Request) (string, error) {
	var userID string
	//we MUST copy the body because it can be read Only once for each http.request
	buf, _ := ioutil.ReadAll(r.Body)
	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))

	r.Body = rdr1
	tokenAuth, err := ExtractTokenMetadata(*r)
	if err != nil {
		fmt.Println("UNAUTHORIZED: Token expired or not detected")
		return "", err
	}
	authD := authEndpoint.FetchAuthRequest{AuthD: tokenAuth.Details}
	userID, err = FetchAuth(&authD)
	if err != nil || userID == "" {
		fmt.Println("UNAUTHORIZED: Token deleted")
		return "", fmt.Errorf("Error: invalid Token")
	}
	fmt.Println("USERID = ", userID)
	return userID, err
}

// makeCreateHandler creates the handler logic
func makeCreateHandler(m *mux.Router, endpoints endpoint.Endpoints, options []http.ServerOption) {
	m.Methods("POST").Path("/").Handler(handlers.CORS(handlers.AllowedMethods([]string{"POST"}), handlers.AllowedOrigins([]string{"*"}))(http.NewServer(endpoints.CreateEndpoint, decodeCreateRequest, encodeCreateResponse, options...)))
}

// decodeCreateRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeCreateRequest(_ context.Context, r *http1.Request) (interface{}, error) {
	buf, _ := ioutil.ReadAll(r.Body)
	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))

	r.Body = rdr1
	userID, err := checkToken(r)
	if err != nil {
		return endpoint.CreateRequest{}, err
	}
	r.Body = rdr2
	req := endpoint.CreateRequest{}
	req.Transaction.AccountID = bson.ObjectIdHex(userID)
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
	buf, _ := ioutil.ReadAll(r.Body)
	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))

	r.Body = rdr1
	userID, err := checkToken(r)
	if err != nil {
		return endpoint.AcceptRequest{}, err
	}
	r.Body = rdr2
	req := endpoint.AcceptRequest{}
	req.Transaction.AccountID = bson.ObjectIdHex(userID)
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
	buf, _ := ioutil.ReadAll(r.Body)
	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))

	r.Body = rdr1
	userID, err := checkToken(r)
	if err != nil {
		return endpoint.AcceptRequest{}, err
	}
	r.Body = rdr2
	req := endpoint.RejectRequest{}
	req.Transaction.AccountID = bson.ObjectIdHex(userID)
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
	buf := []byte(`{}`)
	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))

	r.Body = rdr1
	userID, err := checkToken(r)
	if err != nil {
		return endpoint.AcceptRequest{}, err
	}
	r.Body = rdr2
	req := endpoint.GetAllRequest{AccountID: userID}
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
