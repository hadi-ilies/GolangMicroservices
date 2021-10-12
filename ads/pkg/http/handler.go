package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	endpoint "golangmicroservices/ads/pkg/endpoint"
	"io/ioutil"
	http1 "net/http"
	"net/url"
	"time"

	authEndpoint "golangmicroservices/auths/pkg/endpoint"

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
	m.Methods("POST").Path("/createAd").Handler(handlers.CORS(handlers.AllowedMethods([]string{"POST"}), handlers.AllowedOrigins([]string{"*"}))(http.NewServer(endpoints.CreateEndpoint, decodeCreateRequest, encodeCreateResponse, options...)))
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
	req.Ad.AccountID = bson.ObjectIdHex(userID)
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
	buf, _ := ioutil.ReadAll(r.Body)
	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))

	r.Body = rdr1
	userID, err := checkToken(r)
	if err != nil {
		return endpoint.UpdateRequest{}, err
	}
	r.Body = rdr2
	req := endpoint.UpdateRequest{}
	req.Ad.AccountID = bson.ObjectIdHex(userID)
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
	buf, _ := ioutil.ReadAll(r.Body)
	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))

	r.Body = rdr1
	userID, err := checkToken(r)
	if err != nil {
		return endpoint.DeleteRequest{}, err
	}
	r.Body = rdr2
	req := endpoint.DeleteRequest{}
	req.Ad.AccountID = bson.ObjectIdHex(userID)
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
	//since auth request is a POST we need to add a body
	buf := []byte(`{}`)
	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))

	r.Body = rdr1
	_, err := checkToken(r)
	if err != nil {
		return endpoint.GetRequest{}, err
	}
	r.Body = rdr2
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
	buf, _ := ioutil.ReadAll(r.Body)
	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))

	r.Body = rdr1
	_, err := checkToken(r)
	if err != nil {
		return endpoint.GetAllByUserRequest{}, err
	}
	r.Body = rdr2
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
