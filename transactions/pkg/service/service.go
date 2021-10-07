package service

import (
	"context"
	"fmt"
	"golangmicroservices/transactions/pkg/db"
	"golangmicroservices/transactions/pkg/domain"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// TransactionsService describes the service.
type TransactionsService interface {
	// Add your methods here
	// e.x: Foo(ctx context.Context,s string)(rs string, err error)
	// - Make an offer on an ad
	Create(ctx context.Context, transaction domain.Transaction) (domain.Transaction, error)
	// - Accept an offer on its own ad
	Accept(ctx context.Context, transaction domain.Transaction) (domain.Transaction, error)
	// - Refuse an offer on its own ad
	Reject(ctx context.Context, transaction domain.Transaction) error
	// - List all its own transaction
	GetAll(ctx context.Context, accountID string) ([]domain.Transaction, error)
}

type basicTransactionsService struct{}

func (b *basicTransactionsService) Create(ctx context.Context, transaction domain.Transaction) (d0 domain.Transaction, e1 error) {
	transaction.Id = bson.NewObjectId()
	transaction.CreatedAt = time.Now()
	transaction.Accepted = false
	transaction.Rejected = false
	session, err := db.GetMongoSession()
	if err != nil {
		return d0, err
	}
	defer session.Close()
	c := session.DB("my_store").C("transactions")
	e1 = c.Insert(&transaction)
	return transaction, e1
}
func (b *basicTransactionsService) Accept(ctx context.Context, transaction domain.Transaction) (d0 domain.Transaction, e1 error) {
	//implement the business logic of Accept
	transaction.Accepted = true
	transaction.Rejected = false
	session, err := db.GetMongoSession()
	if err != nil {
		return d0, e1
	}
	defer session.Close()
	c := session.DB("my_store").C("transactions")
	e1 = c.Update(bson.M{"_id": transaction.Id, "account_id": transaction.AccountID, "ad_id": transaction.AdID}, transaction)
	if e1 != nil {
		return d0, e1
	}
	e1 = c.Find(bson.M{"_id": transaction.Id, "account_id": transaction.AccountID, "ad_id": transaction.AdID}).One(&d0)
	return d0, e1
}
func (b *basicTransactionsService) Reject(ctx context.Context, transaction domain.Transaction) (e0 error) {
	//implement the business logic of Reject
	transaction.Rejected = true
	transaction.Accepted = false
	session, err := db.GetMongoSession()
	if err != nil {
		return e0
	}
	defer session.Close()
	c := session.DB("my_store").C("transactions")
	e0 = c.Update(bson.M{"_id": transaction.Id, "account_id": transaction.AccountID, "ad_id": transaction.AdID}, transaction)
	if e0 != nil {
		return e0
	}
	e0 = c.Find(bson.M{"_id": transaction.Id, "account_id": transaction.AccountID, "ad_id": transaction.AdID}).One(&transaction)
	return e0
}
func (b *basicTransactionsService) GetAll(ctx context.Context, accountID string) (d0 []domain.Transaction, e1 error) {
	//implement the business logic of GetAll
	session, err := db.GetMongoSession()
	if err != nil {
		return d0, e1
	}
	defer session.Close()
	c := session.DB("my_store").C("transactions")
	//should be passed in request
	fmt.Println("DEBUGOS = ", accountID)
	e1 = c.Find(bson.M{"account_id": bson.ObjectIdHex(accountID)}).All(&d0)
	return d0, e1
}

// NewBasicTransactionsService returns a naive, stateless implementation of TransactionsService.
func NewBasicTransactionsService() TransactionsService {
	return &basicTransactionsService{}
}

// New returns a TransactionsService with all of the expected middleware wired in.
func New(middleware []Middleware) TransactionsService {
	var svc TransactionsService = NewBasicTransactionsService()
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}
