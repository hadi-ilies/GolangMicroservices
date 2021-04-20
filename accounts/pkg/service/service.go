package service

import (
	"context"
	"fmt"
	"golangmicroservices/accounts/pkg/db"
	"golangmicroservices/accounts/pkg/domain"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"
)

var mySigningKey = []byte(os.Getenv("SECRET_KEY"))

func GetJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["client"] = "Krissanawat"
	claims["aud"] = "billing.jwtgo.io"
	claims["iss"] = "jwtgo.io"
	claims["exp"] = time.Now().Add(time.Minute * 1).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

// AccountsService describes the service.
type AccountsService interface {
	// Add your methods here
	// e.x: Foo(ctx context.Context,s string)(rs string, err error)
	//Create an account
	SignUp(ctx context.Context, account domain.Account) (domain.Account, error)
	//Login
	SignIn(ctx context.Context, auth domain.Auth) (string, error)
	//Update informations of its own account
	Update(ctx context.Context, account domain.Account) (domain.Account, error)
	//Delete its own account
	Delete(ctx context.Context) error
	//Fully read its own account
	Get(ctx context.Context) ([]domain.Account, error)
	//Partially read any user account
	GetUserInfo(ctx context.Context, username string) (domain.Account, error)
	//Add funds to it's own balance
	AddFunds(ctx context.Context, funds uint64) (domain.Account, error)
}

type basicAccountsService struct{}

func (b *basicAccountsService) SignUp(ctx context.Context, account domain.Account) (d0 domain.Account, e1 error) {
	// TODO implement the business logic of SignUp
	account.Id = bson.NewObjectId()
	account.CreatedAt = time.Now()
	session, err := db.GetMongoSession()
	if err != nil {
		return d0, err
	}
	defer session.Close()
	c := session.DB("my_store").C("accounts")
	e1 = c.Insert(&account)
	return account, e1
}
func (b *basicAccountsService) SignIn(ctx context.Context, account domain.Auth) (d0 string, e1 error) {
	// TODO implement the business logic of SignIn
	//check if account exiist
	var myAccount domain.Account
	session, err := db.GetMongoSession()
	if err != nil {
		return d0, err
	}
	defer session.Close()
	c := session.DB("my_store").C("accounts")
	e1 = c.Find(bson.M{"email": account.Email}).One(&myAccount)
	if err != nil {
		return d0, err
	}
	//account exist
	validToken, err := GetJWT()
	fmt.Println(validToken)
	return validToken, e1
}
func (b *basicAccountsService) Update(ctx context.Context, account domain.Account) (d0 domain.Account, e1 error) {
	// TODO implement the business logic of Update
	session, err := db.GetMongoSession()
	if err != nil {
		return d0, err
	}
	defer session.Close()
	data, err := bson.Marshal(account)
	e1 = err
	if e1 != nil {
		return d0, e1
	}
	c := session.DB("my_store").C("accounts")
	e1 = c.Update(bson.M{"_id": bson.ObjectIdHex(string(account.Id))}, data)
	e1 = c.Find(bson.M{"_id": bson.ObjectIdHex(string(account.Id))}).One(&d0)
	return d0, e1
}
func (b *basicAccountsService) Delete(ctx context.Context) (e0 error) {
	// TODO implement the business logic of Delete
	session, err := db.GetMongoSession()
	if err != nil {
		return err
	}
	defer session.Close()
	//TODO find way to get od from token
	// c := session.DB("my_store").C("accounts")
	// return c.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	return e0
}
func (b *basicAccountsService) Get(ctx context.Context) (d0 []domain.Account, e1 error) {
	// TODO implement the business logic of Get
	session, err := db.GetMongoSession()
	if err != nil {
		return d0, e1
	}
	defer session.Close()
	c := session.DB("my_store").C("accounts")
	e1 = c.Find(nil).All(&d0)
	return d0, e1
}
func (b *basicAccountsService) GetUserInfo(ctx context.Context, username string) (d0 domain.Account, e1 error) {
	// TODO implement the business logic of GetUserInfo
	session, err := db.GetMongoSession()
	if err != nil {
		return d0, e1
	}
	defer session.Close()
	c := session.DB("my_store").C("accounts")
	e1 = c.Find(bson.M{"username": username}).One(&d0)

	return d0, e1
}
func (b *basicAccountsService) AddFunds(ctx context.Context, funds uint64) (d0 domain.Account, e1 error) {
	// TODO implement the business logic of AddFunds
	return d0, e1
}

// NewBasicAccountsService returns a naive, stateless implementation of AccountsService.
func NewBasicAccountsService() AccountsService {
	return &basicAccountsService{}
}

// New returns a AccountsService with all of the expected middleware wired in.
func New(middleware []Middleware) AccountsService {
	var svc AccountsService = NewBasicAccountsService()
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}
