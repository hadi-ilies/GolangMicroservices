package service

import (
	"context"
	"fmt"
	"golangmicroservices/ads/pkg/db"
	"golangmicroservices/ads/pkg/domain"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// AdsService describes the service.
type AdsService interface {
	// Add your methods here
	// e.x: Foo(ctx context.Context,s string)(rs string, err error)
	// - Create an ad
	Create(ctx context.Context, ad domain.Ad) (domain.Ad, error)
	// - Update one of its own ad
	Update(ctx context.Context, ad domain.Ad) (domain.Ad, error)
	// - Delete one of its own ad
	Delete(ctx context.Context, ad domain.Ad) error
	// - Read any ad (get all ads)
	Get(ctx context.Context) ([]domain.Ad, error)
	// - Get a list of ads searching by keywords
	GetAllByKeyWord(ctx context.Context, keywords string) ([]domain.Ad, error)
	// - Get a list of all the ads of a user by clientID
	GetAllByUser(ctx context.Context, targetAccountID string) ([]domain.Ad, error)
}

type basicAdsService struct{}

func (b *basicAdsService) Create(ctx context.Context, ad domain.Ad) (d0 domain.Ad, e1 error) {
	// TODO implement the business logic of Create
	ad.Id = bson.NewObjectId()
	ad.CreatedAt = time.Now()
	session, err := db.GetMongoSession()
	if err != nil {
		return d0, err
	}
	defer session.Close()
	c := session.DB("my_store").C("ads")
	e1 = c.Insert(&ad)
	return ad, e1
}
func (b *basicAdsService) Update(ctx context.Context, ad domain.Ad) (d0 domain.Ad, e1 error) {
	// TODO implement the business logic of Update
	session, err := db.GetMongoSession()
	if err != nil {
		return d0, err
	}
	defer session.Close()
	c := session.DB("my_store").C("ads")
	e1 = c.Update(bson.M{"_id": ad.Id, "account_id": ad.AccountID}, ad)
	if e1 != nil {
		return ad, e1
	}
	e1 = c.Find(bson.M{"_id": ad.Id, "account_id": ad.AccountID}).One(&d0)
	return d0, e1
}
func (b *basicAdsService) Delete(ctx context.Context, ad domain.Ad) (e1 error) {
	// TODO implement the business logic of Delete
	session, err := db.GetMongoSession()
	if err != nil {
		return err
	}
	defer session.Close()
	//TODO find way to get od from token
	c := session.DB("my_store").C("ads")
	e1 = c.Remove(bson.M{"account_id": ad.AccountID, "_id": ad.Id})
	return e1
}
func (b *basicAdsService) Get(ctx context.Context) (d0 []domain.Ad, e1 error) {
	// TODO implement the business logic of Get
	session, err := db.GetMongoSession()
	if err != nil {
		return d0, e1
	}
	defer session.Close()
	c := session.DB("my_store").C("ads")
	e1 = c.Find(nil).All(&d0)
	return d0, e1
}

func (b *basicAdsService) GetAllByKeyWord(ctx context.Context, keywords string) (d0 []domain.Ad, e1 error) {
	// TODO implement the business logic of GetAllByKeyWord
	return d0, e1
}
func (b *basicAdsService) GetAllByUser(ctx context.Context, targetAccountID string) (d0 []domain.Ad, e1 error) {
	// TODO implement the business logic of GetAllByUser
	session, err := db.GetMongoSession()
	if err != nil {
		return d0, e1
	}
	defer session.Close()
	c := session.DB("my_store").C("ads")
	//TODO should be passed in request
	fmt.Println("DEBUGOS = ", targetAccountID)
	e1 = c.Find(bson.M{"account_id": bson.ObjectIdHex(targetAccountID)}).All(&d0)

	return d0, e1
}

// NewBasicAdsService returns a naive, stateless implementation of AdsService.
func NewBasicAdsService() AdsService {
	return &basicAdsService{}
}

// New returns a AdsService with all of the expected middleware wired in.
func New(middleware []Middleware) AdsService {
	var svc AdsService = NewBasicAdsService()
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}
