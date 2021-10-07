package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golangmicroservices/accounts/pkg/db"
	endpoint "golangmicroservices/accounts/pkg/endpoint"
	http1 "net/http"
	"os"

	"gopkg.in/mgo.v2/bson"

	"strings"

	"github.com/dgrijalva/jwt-go"
	http "github.com/go-kit/kit/transport/http"
	handlers "github.com/gorilla/handlers"
	mux "github.com/gorilla/mux"
)

// ---- IS AUTH ------

type AccessDetails struct {
	AccessUuid string `bson:"access_uuid"`
	UserID     string `bson:"user_id"`
}

func ExtractToken(r *http1.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func VerifyToken(r *http1.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func TokenValid(r *http1.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		fmt.Println("Error: TOKEN CLAIMS\n")
		return err
	}
	return nil
}

func ExtractTokenMetadata(r *http1.Request) (*AccessDetails, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userID, exist := claims["user_id"]
		if !exist {
			return nil, fmt.Errorf(("claims error"))
		}
		return &AccessDetails{
			AccessUuid: accessUuid,
			UserID:     userID.(string),
		}, nil
	}
	fmt.Println("Error spotted3\n")
	return nil, err
}

//FetchAuth() accepts the AccessDetails from the ExtractTokenMetadata function, then looks it up in mongodb.
//If the record is not found, it may mean the token has expired, hence an error is thrown.
func FetchAuth(authD *AccessDetails) (string, error) {
	session, err := db.GetMongoSession()

	if err != nil {
		fmt.Println("Error spotted1\n")
		return "", err
	}
	defer session.Close()
	c := session.DB("my_store").C("auths")
	accessDetailField := AccessDetails{}
	err = c.Find(bson.M{"access_uuid": authD.AccessUuid}).One(&accessDetailField)
	if err != nil {
		fmt.Println("Error spotted2\n")
		return "", err
	}
	return accessDetailField.UserID, nil
}

//TODO DELETE
func IsAuthorized(r *http1.Request) (string, error) {
	return "", nil
}

// func IsAuthorized(r *http1.Request) (string, error) {
// 	//extract token
// 	if r.Header["Authorization"] != nil {
// 		token, err := jwt.Parse(r.Header["Authorization"][0], func(token *jwt.Token) (interface{}, error) {
// 			//Make sure that the token method conform to "SigningMethodHMAC"
// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, fmt.Errorf(("Invalid Signing Method"))
// 			}
// 			aud := "billing.jwtgo.io"
// 			checkAudience := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
// 			if !checkAudience {
// 				return nil, fmt.Errorf(("invalid aud"))
// 			}
// 			// verify iss claim
// 			iss := "jwtgo.io"
// 			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
// 			if !checkIss {
// 				return nil, fmt.Errorf(("invalid iss"))
// 			}
// 			return []byte(os.Getenv("SECRET_KEY")), nil
// 		})
// 		if err != nil {
// 			return "", err
// 		}

// 		if token.Valid {
// 			return r.Header["Authorization"][0], nil
// 		}

// 	}
// 	return "", fmt.Errorf(("no Token detected"))
// }

// ----END IS AUTH----

// makeSignUpHandler creates the handler logic
func makeSignUpHandler(m *mux.Router, endpoints endpoint.Endpoints, options []http.ServerOption) {
	m.Methods("POST", "OPTIONS").Path("/signUp").Handler(
		handlers.CORS(
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedHeaders([]string{"Content-Type", "Content-Length"}),
			handlers.AllowedMethods([]string{"POST"}),
		)(http.NewServer(endpoints.SignUpEndpoint, decodeSignUpRequest, encodeSignUpResponse, options...)))
}

// decodeSignUpRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeSignUpRequest(_ context.Context, r *http1.Request) (interface{}, error) {
	req := endpoint.SignUpRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// encodeSignUpResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer
func encodeSignUpResponse(ctx context.Context, w http1.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(response)
	return
}

// makeSignInHandler creates the handler logic
func makeSignInHandler(m *mux.Router, endpoints endpoint.Endpoints, options []http.ServerOption) {
	m.Methods("POST").Path("/signIn").Handler(handlers.CORS(handlers.AllowedMethods([]string{"POST"}), handlers.AllowedOrigins([]string{"*"}))(http.NewServer(endpoints.SignInEndpoint, decodeSignInRequest, encodeSignInResponse, options...)))
}

// decodeSignInRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeSignInRequest(_ context.Context, r *http1.Request) (interface{}, error) {
	req := endpoint.SignInRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// encodeSignInResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer
func encodeSignInResponse(ctx context.Context, w http1.ResponseWriter, response interface{}) (err error) {
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
	_, err := IsAuthorized(r)

	if err != nil {
		return nil, err
	}
	req := endpoint.UpdateRequest{}
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
	m.Methods("DELETE", "OPTIONS").Path("/").Handler(
		handlers.CORS(
			handlers.AllowedMethods([]string{"DELETE"}),
			handlers.AllowedHeaders([]string{"Content-Type", "Content-Length"}),
			handlers.AllowedOrigins([]string{"*"}),
		)(http.NewServer(endpoints.DeleteEndpoint, decodeDeleteRequest, encodeDeleteResponse, options...)))
}

// decodeDeleteRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeDeleteRequest(_ context.Context, r *http1.Request) (interface{}, error) {
	token, err := IsAuthorized(r)

	if err != nil {
		return nil, err
	}
	req := endpoint.DeleteRequest{Token: token}
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
	return req, nil
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

// makeGetUserInfoHandler creates the handler logic
func makeGetUserInfoHandler(m *mux.Router, endpoints endpoint.Endpoints, options []http.ServerOption) {
	m.Methods("POST", "OPTIONS").Path("/GetUserInfo").Handler(
		handlers.CORS(
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedHeaders([]string{"Content-Type", "Content-Length"}),
			handlers.AllowedMethods([]string{"POST"}),
		)(http.NewServer(endpoints.GetUserInfoEndpoint, decodeGetUserInfoRequest, encodeGetUserInfoResponse, options...)))
}

// decodeGetUserInfoRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeGetUserInfoRequest(_ context.Context, r *http1.Request) (interface{}, error) {
	var userID string
	tokenAuth, err := ExtractTokenMetadata(r)
	if err != nil {
		fmt.Println("tokenAuth ERROR")
		return nil, err
	}
	userID, err = FetchAuth(tokenAuth)
	if err != nil {
		fmt.Println("UNAUTHORIZED")
		return nil, err
	}
	fmt.Println("USERID = ", userID)
	//you can proceed

	// _, err := IsAuthorized(r)

	// if err != nil {
	// 	return nil, err
	// }
	req := endpoint.GetUserInfoRequest{}
	err = json.NewDecoder(r.Body).Decode(&req)
	fmt.Println("Lolilol")
	return req, err
}

// encodeGetUserInfoResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer
func encodeGetUserInfoResponse(ctx context.Context, w http1.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(response)
	return
}

// makeAddFundsHandler creates the handler logic
func makeAddFundsHandler(m *mux.Router, endpoints endpoint.Endpoints, options []http.ServerOption) {
	m.Methods("POST").Path("/add-funds").Handler(handlers.CORS(handlers.AllowedMethods([]string{"POST"}), handlers.AllowedOrigins([]string{"*"}))(http.NewServer(endpoints.AddFundsEndpoint, decodeAddFundsRequest, encodeAddFundsResponse, options...)))
}

// decodeAddFundsRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeAddFundsRequest(_ context.Context, r *http1.Request) (interface{}, error) {
	token, err := IsAuthorized(r)

	if err != nil {
		return nil, err
	}
	req := endpoint.AddFundsRequest{Token: token}
	err = json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// encodeAddFundsResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer
func encodeAddFundsResponse(ctx context.Context, w http1.ResponseWriter, response interface{}) (err error) {
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

// makeMeHandler creates the handler logic
func makeMeHandler(m *mux.Router, endpoints endpoint.Endpoints, options []http.ServerOption) {
	m.Methods("GET").Path("/me").Handler(handlers.CORS(handlers.AllowedMethods([]string{"GET"}), handlers.AllowedOrigins([]string{"*"}))(http.NewServer(endpoints.MeEndpoint, decodeMeRequest, encodeMeResponse, options...)))
}

// decodeMeRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeMeRequest(_ context.Context, r *http1.Request) (interface{}, error) {
	token, err := IsAuthorized(r)

	if err != nil {
		return nil, err
	}
	req := endpoint.MeRequest{Token: token}
	return req, err
}

// encodeMeResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer
func encodeMeResponse(ctx context.Context, w http1.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(response)
	return
}

// makeLogoutHandler creates the handler logic
func makeLogoutHandler(m *mux.Router, endpoints endpoint.Endpoints, options []http.ServerOption) {
	m.Methods("POST").Path("/logout").Handler(handlers.CORS(handlers.AllowedMethods([]string{"POST"}), handlers.AllowedOrigins([]string{"*"}))(http.NewServer(endpoints.LogoutEndpoint, decodeLogoutRequest, encodeLogoutResponse, options...)))
}

// decodeLogoutRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeLogoutRequest(_ context.Context, r *http1.Request) (interface{}, error) {
	tokenAuth, err := ExtractTokenMetadata(r)
	if err != nil {
		fmt.Println("tokenAuth ERROR")
		return nil, err
	}
	req := endpoint.LogoutRequest{}
	fmt.Println("tokenAuth.AccessUuid = ", tokenAuth.AccessUuid)
	//todo change req.token name to access_uuid
	req.Token = tokenAuth.AccessUuid
	err = json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// encodeLogoutResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer
func encodeLogoutResponse(ctx context.Context, w http1.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(response)
	return
}
