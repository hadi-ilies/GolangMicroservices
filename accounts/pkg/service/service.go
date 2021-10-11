package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"golangmicroservices/accounts/pkg/db"
	"golangmicroservices/accounts/pkg/domain"
	authEndpoint "golangmicroservices/auths/pkg/endpoint"
	"io/ioutil"
	"log"
	http1 "net/http"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

var MySigningKey = []byte(os.Getenv("SECRET_KEY"))

// CreateAuth: save auth tokens inside mongodb by calling the auth service
func CreateAuth(userID string, td *authEndpoint.GetJWTResponse) error {
	url := "http://auths:8084/create-auth"
	spaceClient := http1.Client{
		Timeout: time.Second * 20, // Timeout after 2 seconds
	}
	jsonData := map[string]interface{}{"user_id": userID, "td": td.Token}
	s, _ := json.Marshal(jsonData)
	req, err := http1.NewRequest(http1.MethodPost, url, bytes.NewReader(s))
	if err != nil {
		fmt.Println("CALAMAR1")
		return err
	}
	req.Header.Set("Access-Control-Allow-Origin", "*")

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		fmt.Println("CALAMAR2")
		return getErr
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		fmt.Println("CALAMAR3")
		return readErr
	}
	myAccountResponse := authEndpoint.CreateAuthResponse{}
	jsonErr := json.Unmarshal(body, &myAccountResponse)
	if jsonErr != nil {
		fmt.Println("CALAMAR4")
		return jsonErr
	}

	fmt.Println("CREATE AUTH LOG = ", myAccountResponse.Log)
	return nil
}

// GetJWT: generate a access and refresh tokens by calling the auth service
func GetJWT(userID string) (*authEndpoint.GetJWTResponse, error) {
	url := "http://auths:8084/get-jwt"
	spaceClient := http1.Client{
		Timeout: time.Second * 20, // Timeout after 2 seconds
	}
	jsonData := map[string]string{"user_id": userID}
	s, _ := json.Marshal(jsonData)
	req, err := http1.NewRequest(http1.MethodPost, url, bytes.NewReader(s))
	if err != nil {
		fmt.Println("CALAMAR1")
		return nil, err
	}
	req.Header.Set("Access-Control-Allow-Origin", "*")

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
	myAccountResponse := authEndpoint.GetJWTResponse{}
	jsonErr := json.Unmarshal(body, &myAccountResponse)
	if jsonErr != nil {
		fmt.Println("CALAMAR4")
		return nil, jsonErr
	}

	fmt.Println("AccessToken = ", myAccountResponse.Token.AccessToken)
	return &myAccountResponse, nil
}

//When a user logs out, we will instantly revoke/invalidate their JWT. This is achieved by deleting the JWT metadata from our mongodb.
func DeleteAuth(givenUuid string) error {
	url := "http://auths:8084/delete-auth"
	spaceClient := http1.Client{
		Timeout: time.Second * 20, // Timeout after 2 seconds
	}
	jsonData := map[string]interface{}{"given_uuid": givenUuid}
	s, _ := json.Marshal(jsonData)
	req, err := http1.NewRequest(http1.MethodPost, url, bytes.NewReader(s))
	if err != nil {
		fmt.Println("CALAMAR1")
		return err
	}
	req.Header.Set("Access-Control-Allow-Origin", "*")

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		fmt.Println("CALAMAR2")
		return getErr
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		fmt.Println("CALAMAR3")
		return readErr
	}
	myAccountResponse := authEndpoint.DeleteAuthResponse{}
	jsonErr := json.Unmarshal(body, &myAccountResponse)
	if jsonErr != nil {
		fmt.Println("CALAMAR4")
		return jsonErr
	}

	fmt.Println("DELETE AUTH LOG = ", myAccountResponse.Log, "| ", myAccountResponse.Err)
	return nil
}

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

// Since we'll be getting the hashed password from the DB it
func comparePasswords(hashedPwd string, plainPwd []byte) bool {
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
	Logout(ctx context.Context, accessUUID string) error
	//Update informations of its own account
	Update(ctx context.Context, account domain.Account) (domain.Account, error)
	//Delete its own account
	Delete(ctx context.Context, userID string) error
	//Fully read its own account
	Me(ctx context.Context, userID string) (domain.Account, error)
	//get all accounts //tmp
	Get(ctx context.Context) ([]domain.Account, error)
	//Partially read any user account
	GetUserInfo(ctx context.Context, username string) (domain.Account, error)
	//Add funds to it's own balance
	AddFunds(ctx context.Context, userID string, funds uint64) (domain.Account, error)
}

type basicAccountsService struct{}

func (b *basicAccountsService) SignUp(ctx context.Context, account domain.Account) (d0 domain.Account, e1 error) {
	// the business logic of SignUp
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
	// the business logic of SignIn
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
	myAccount.Token = validToken.Token.AccessToken
	//add auth
	saveErr := CreateAuth(myAccount.Id.Hex(), validToken)
	if saveErr != nil {
		return d0, saveErr
	}
	tokens := map[string]string{
		"access_token":  validToken.Token.AccessToken,
		"refresh_token": validToken.Token.RefreshToken,
	}
	//update current user Token
	_, e1 = b.Update(ctx, myAccount)
	return tokens, e1
}

func (b *basicAccountsService) Update(ctx context.Context, account domain.Account) (d0 domain.Account, e1 error) {
	// the business logic of Update
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

func (b *basicAccountsService) Delete(ctx context.Context, userID string) (e0 error) {
	//the business logic of Delete
	session, err := db.GetMongoSession()
	if err != nil {
		return err
	}
	defer session.Close()
	c := session.DB("my_store").C("accounts")
	e0 = c.Remove(bson.M{"_id": bson.ObjectIdHex(userID)})
	return e0
}
func (b *basicAccountsService) Get(ctx context.Context) (d0 []domain.Account, e1 error) {
	//the business logic of Get
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
	// the business logic of GetUserInfo
	session, err := db.GetMongoSession()
	if err != nil {
		return d0, e1
	}
	defer session.Close()
	c := session.DB("my_store").C("accounts")
	e1 = c.Find(bson.M{"username": username}).One(&d0)

	return d0, e1
}
func (b *basicAccountsService) AddFunds(ctx context.Context, userID string, funds uint64) (d0 domain.Account, e1 error) {
	// the business logic of AddFunds
	session, err := db.GetMongoSession()
	if err != nil {
		return d0, e1
	}
	defer session.Close()
	c := session.DB("my_store").C("accounts")
	e1 = c.Find(bson.M{"_id": bson.ObjectIdHex(userID)}).One(&d0)
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

func (b *basicAccountsService) Me(ctx context.Context, userID string) (d0 domain.Account, e1 error) {
	// the business logic of Me
	var myAccount domain.Account
	session, err := db.GetMongoSession()
	if err != nil {
		return d0, err
	}
	defer session.Close()
	c := session.DB("my_store").C("accounts")
	e1 = c.Find(bson.M{"_id": bson.ObjectIdHex(userID)}).One(&myAccount)
	if err != nil {
		return d0, err
	}
	d0 = myAccount
	return d0, e1
}

//here token is the uuid not a real token (I need to change the name)
func (b *basicAccountsService) Logout(ctx context.Context, accessUUID string) (e0 error) {
	e0 = DeleteAuth(accessUUID)
	return e0
}
