package service

import (
	"context"
	"fmt"
	"golangmicroservices/auths/pkg/db"
	"golangmicroservices/auths/pkg/domain"
	http1 "net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
	"gopkg.in/mgo.v2/bson"
)

// AuthsService describes the service.
type AuthsService interface {
	// Add your methods here
	// e.x: Foo(ctx context.Context,s string)(rs string, err error)
	//get JWT
	GetJWT(ctx context.Context, userID string) (token domain.TokenDetails, err error)
	//generate auth in db
	CreateAuth(ctx context.Context, userID string, td domain.TokenDetails) (log string, err error)
	//delete auth from db
	DeleteAuth(ctx context.Context, givenUuid string) (log string, err error)
	//check if valid token return token
	ExtractTokenMetadata(ctx context.Context, r http1.Request) (details domain.AccessDetails, err error)
	//get userID --> account
	FetchAuth(ctx context.Context, authD domain.AccessDetails) (userID string, err error)
}

type basicAuthsService struct{}

func (b *basicAuthsService) GetJWT(ctx context.Context, userID string) (token domain.TokenDetails, err error) {
	token = domain.TokenDetails{}

	token.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	token.AccessUuid = uuid.NewV4().String()
	token.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	token.RefreshUuid = uuid.NewV4().String()
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = token.AccessUuid
	atClaims["user_id"] = userID
	atClaims["exp"] = token.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return domain.TokenDetails{}, err
	}
	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = token.RefreshUuid
	rtClaims["user_id"] = userID
	rtClaims["exp"] = token.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	token.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return domain.TokenDetails{}, err
	}
	return token, err
}

func (b *basicAuthsService) CreateAuth(ctx context.Context, userID string, td domain.TokenDetails) (log string, err error) {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()
	session, err2 := db.GetMongoSession()

	if err2 != nil {
		return "Error GetMongoSession", err2
	}
	defer session.Close()
	c := session.DB("my_store").C("auths")
	errAccess := c.Insert(bson.M{"access_uuid": td.AccessUuid, "user_id": userID, "time": at.Sub(now)})
	if errAccess != nil {
		return "Error Insert access_uuid", errAccess
	}
	errRefresh := c.Insert(bson.M{"refresh_uuid": td.RefreshUuid, "user_id": userID, "time": rt.Sub(now)})
	if errRefresh != nil {
		return "Error Insert refresh_uuid", errRefresh
	}
	return "success", err
}

func (b *basicAuthsService) DeleteAuth(ctx context.Context, givenUuid string) (log string, err error) {
	session, err := db.GetMongoSession()
	if err != nil {
		return "Error GetMongoSession", err
	}
	defer session.Close()
	c := session.DB("my_store").C("auths")
	err = c.Remove(bson.M{"access_uuid": givenUuid})
	if err != nil {
		return "Error Remove", err
	}
	fmt.Println("givenUuid = ", givenUuid)
	log = "success"
	return log, err
}

func extractToken(r *http1.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	fmt.Println("BEARER = ", bearToken)
	fmt.Println(strArr)
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func verifyToken(r *http1.Request) (*jwt.Token, error) {
	tokenString := extractToken(r)
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

func tokenValid(r *http1.Request) error {
	token, err := verifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		fmt.Println("Error: TOKEN CLAIMS\n")
		return err
	}
	return nil
}

func (b *basicAuthsService) ExtractTokenMetadata(ctx context.Context, r http1.Request) (details domain.AccessDetails, err error) {
	fmt.Println("DEBUG")
	token, err := verifyToken(&r)
	if err != nil {
		return domain.AccessDetails{}, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return domain.AccessDetails{}, err
		}
		userID, exist := claims["user_id"]
		fmt.Println("USERID = ", userID)
		if !exist {
			return domain.AccessDetails{}, fmt.Errorf(("claims error"))
		}
		return domain.AccessDetails{
			AccessUuid: accessUuid,
			UserID:     userID.(string),
		}, nil
	}
	fmt.Println("Error spotted3\n")
	return domain.AccessDetails{}, err
}
func (b *basicAuthsService) FetchAuth(ctx context.Context, authD domain.AccessDetails) (userID string, err error) {
	session, err := db.GetMongoSession()

	if err != nil {
		fmt.Println("Error spotted1\n")
		return "", err
	}
	defer session.Close()
	c := session.DB("my_store").C("auths")
	accessDetailField := domain.AccessDetails{}
	fmt.Println("AUTH authD.AccessUuid = ", authD.AccessUuid)
	err = c.Find(bson.M{"access_uuid": authD.AccessUuid}).One(&accessDetailField)
	if err != nil {
		fmt.Println("Error spotted2\n")
		return "", err
	}
	return accessDetailField.UserID, nil
}

// NewBasicAuthsService returns a naive, stateless implementation of AuthsService.
func NewBasicAuthsService() AuthsService {
	return &basicAuthsService{}
}

// New returns a AuthsService with all of the expected middleware wired in.
func New(middleware []Middleware) AuthsService {
	var svc AuthsService = NewBasicAuthsService()
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}
