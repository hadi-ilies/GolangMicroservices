package service

import (
	"context"
	"fmt"
	"golangmicroservices/accounts/pkg/db"
	"golangmicroservices/accounts/pkg/domain"
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

var MySigningKey = []byte(os.Getenv("SECRET_KEY"))

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

func CreateAuth(userID string, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	session, err := db.GetMongoSession()
	if err != nil {
		return err
	}
	defer session.Close()
	c := session.DB("my_store").C("auths")
	// e1 = c.Insert(&account)
	errAccess := c.Insert(bson.M{"access_uuid": td.AccessUuid, "user_id": userID, "time": at.Sub(now)})
	if errAccess != nil {
		return errAccess
	}
	errRefresh := c.Insert(bson.M{"refresh_uuid": td.RefreshUuid, "user_id": userID, "time": rt.Sub(now)})
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func GetJWT(userID string) (*TokenDetails, error) {
	td := &TokenDetails{}
	var err error

	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewV4().String()
	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userID
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userID
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

//When a user logs out, we will instantly revoke/invalidate their JWT. This is achieved by deleting the JWT metadata from our mongodb.
func DeleteAuth(givenUuid string) error {
	session, err := db.GetMongoSession()
	if err != nil {
		return err
	}
	defer session.Close()
	c := session.DB("my_store").C("auths")
	err = c.Remove(bson.M{"access_uuid": givenUuid})
	fmt.Println("Auth RM\n")
	return err
}

// func GetJWT(userID string) (string, error) {
// 	var err error
//   //Creating Access Token
//   atClaims := jwt.MapClaims{}
//   atClaims["authorized"] = true
//   atClaims["user_id"] = userID
//   atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
//   at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
//   token, err := at.SignedString([]byte(os.Getenv("SECRET_KEY")))
//   if err != nil {
// 			fmt.Errorf("Something Went Wrong: %s", err.Error())
// 			return "", err
//   }
//   return token, nil
// }

func hashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool { // Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

// AccountsService describes the service.
type AccountsService interface {
	// Add your methods here
	// e.x: Foo(ctx context.Context,s string)(rs string, err error)
	//Create an account
	SignUp(ctx context.Context, account domain.Account) (domain.Account, error)
	//Login
	SignIn(ctx context.Context, auth domain.Auth) (map[string]string, error)
	//logout
	Logout(ctx context.Context, token string) error
	//Update informations of its own account
	Update(ctx context.Context, account domain.Account) (domain.Account, error)
	//Delete its own account
	Delete(ctx context.Context, token string) error
	//Fully read its own account
	Me(ctx context.Context, token string) (domain.Account, error)
	//get all accounts //tmp
	Get(ctx context.Context) ([]domain.Account, error)
	//Partially read any user account
	GetUserInfo(ctx context.Context, username string) (domain.Account, error)
	//Add funds to it's own balance
	AddFunds(ctx context.Context, token string, funds uint64) (domain.Account, error)
}

type basicAccountsService struct{}

func (b *basicAccountsService) SignUp(ctx context.Context, account domain.Account) (d0 domain.Account, e1 error) {
	// TODO implement the business logic of SignUp
	account.Id = bson.NewObjectId()
	account.CreatedAt = time.Now()
	account.Balance = 0
	account.Password = hashAndSalt([]byte(account.Password))
	session, err := db.GetMongoSession()
	if err != nil {
		return d0, err
	}
	defer session.Close()
	c := session.DB("my_store").C("accounts")
	e1 = c.Insert(&account)
	return account, e1
}
func (b *basicAccountsService) SignIn(ctx context.Context, account domain.Auth) (d0 map[string]string, e1 error) {
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

	//CHECK PASSWORD
	err = bcrypt.CompareHashAndPassword([]byte(myAccount.Password), []byte(account.Password))
	if err != nil {
		log.Println(err)
		return d0, fmt.Errorf(("Wrong Password"))
	}
	//account exist
	validToken, err := GetJWT(myAccount.Id.Hex())
	if err != nil {
		return d0, err
	}
	myAccount.Token = validToken.AccessToken
	//add auth
	saveErr := CreateAuth(myAccount.Id.Hex(), validToken)
	if saveErr != nil {
		return d0, saveErr
	}
	tokens := map[string]string{
		"access_token":  validToken.AccessToken,
		"refresh_token": validToken.RefreshToken,
	}
	//update current user Token
	_, e1 = b.Update(ctx, myAccount)
	return tokens, e1
}

func (b *basicAccountsService) Update(ctx context.Context, account domain.Account) (d0 domain.Account, e1 error) {
	// TODO implement the business logic of Update
	session, err := db.GetMongoSession()
	if err != nil {
		return d0, err
	}
	defer session.Close()
	c := session.DB("my_store").C("accounts")
	e1 = c.Update(bson.M{"_id": account.Id}, account)
	if e1 != nil {
		return account, e1
	}
	e1 = c.Find(bson.M{"_id": account.Id}).One(&d0)
	return d0, e1
}
func (b *basicAccountsService) Delete(ctx context.Context, token string) (e0 error) {
	// TODO implement the business logic of Delete
	session, err := db.GetMongoSession()
	if err != nil {
		return err
	}
	defer session.Close()
	//TODO find way to get od from token
	c := session.DB("my_store").C("accounts")
	e0 = c.Remove(bson.M{"token": token})
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
func (b *basicAccountsService) AddFunds(ctx context.Context, token string, funds uint64) (d0 domain.Account, e1 error) {
	// TODO implement the business logic of AddFunds
	session, err := db.GetMongoSession()
	if err != nil {
		return d0, e1
	}
	defer session.Close()
	c := session.DB("my_store").C("accounts")
	e1 = c.Find(bson.M{"token": token}).One(&d0)
	d0.Balance += funds
	d0, e1 = b.Update(ctx, d0)
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

func (b *basicAccountsService) Me(ctx context.Context, token string) (d0 domain.Account, e1 error) {
	// TODO implement the business logic of Me
	var myAccount domain.Account
	session, err := db.GetMongoSession()
	if err != nil {
		return d0, err
	}
	defer session.Close()
	c := session.DB("my_store").C("accounts")
	e1 = c.Find(bson.M{"token": token}).One(&myAccount)
	if err != nil {
		return d0, err
	}
	d0 = myAccount
	return d0, e1
}

//here token is the uuid not a real token (I need to change the name)
func (b *basicAccountsService) Logout(ctx context.Context, token string) (e0 error) {
	fmt.Println("PROUTAR\n")
	e0 = DeleteAuth(token)
	return e0
}
