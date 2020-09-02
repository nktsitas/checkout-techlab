package gateway

import (
	"encoding/json"
	"fmt"
	"errors"
	"crypto/sha256"
	"sync"
	"time"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/nktsitas/checkout-techlab/bank"
	"github.com/nktsitas/checkout-techlab/db"
)

// Create a GatewayI interface as well as an AuthorizationI interface
// to have the ability of mocking the actions of this package in handlers testing
type GatewayI interface{
	NewAuthorization([]byte, string) (*Authorization, error)
	GetSalt() string
}

type GatewayS struct {}
var Gateway GatewayI

type AuthorizationI interface{
	Void() error
	Capture(float64, string) (*Capture, error)
	Refund(float64, string) (*Refund, error)
	GetCurrency() string
}

type Authorization struct {
	Id string										 `json:"id" swaggerignore:"true"`
	CreditCard *bank.CreditCard  `json:"credit_card"`
	Amount float64							 `json:"amount" example:"100.00"`
	Currency string							 `json:"currency" example:"EUR"`

	captures []*Capture						
	refunds []*Refund							

	void bool					

	mu sync.Mutex
}

type Capture struct {
	Authorization *Authorization
	Amount float64
}

type Refund struct {
	Authorization *Authorization
	Amount float64
}

func (g *GatewayS) NewAuthorization(req_body []byte, salt string) (*Authorization, error) {
	var newAuth Authorization
	err := json.Unmarshal(req_body, &newAuth)
	if err != nil {
		log.WithField("err", err).Error("NewAuthorization - Error Reading request body")
		return nil, fmt.Errorf("Error Unmarshaling JSON - %s", err.Error())
	}

	if newAuth.CreditCard == nil {
		log.WithField("err", err).Error("NewAuthorization - No Credit Card provided")
		return nil, err
	}
	
	if err := newAuth.CreditCard.Validate(); err != nil {
		log.WithField("err", err).Error("NewAuthorization - Invalid Credit Card provided")
		return nil, err
	}

	if newAuth.CreditCard.Number == "4000 0000 0000 0119" {
		log.Error("NewAuthorization - Manually triggered auth error")
		return nil, errors.New("Authorization failure - Unknown Error")
	}

	newAuth.Id = generateID(req_body, salt)

	db.DB.StoreItem(newAuth.Id, &newAuth)

	log.WithField("newAuth", newAuth).Debug("New Authorization Successfully created")

	return &newAuth, nil
}

func (g *GatewayS) GetSalt() string {
	now := time.Now().UnixNano()
	nowString := strconv.FormatInt(now, 10)

	return nowString
}

func generateID(req_body []byte, salt string) string {
	bodyStr := string(req_body)
	theString := bodyStr + salt
	return fmt.Sprintf("%x", sha256.Sum256([]byte(theString)))
}

// ---

func (auth *Authorization) GetCurrency() string {
	return auth.Currency
}

func (auth *Authorization) Void() error {
	auth.mu.Lock()
	defer auth.mu.Unlock()

	if auth.void {
		return errors.New("Void Failure - Transaction already void")
	}

	if auth.TotalCapturedAmount() > 0 {
		return errors.New("Void Failure - Cannot void transaction with captured amount")
	}

	auth.void = true

	log.WithField("auth", auth).Debug("Void Successfully executed.")

	return nil
}

func (auth *Authorization) Capture(amount float64, currency string) (*Capture, error) {
	auth.mu.Lock()
	defer auth.mu.Unlock()

	if auth.void {
		log.Error("Authorization.Capture - transaction is void")

		return nil, errors.New("Capture failure - Cannot capture on void transaction")
	}

	if amount > auth.Amount {
		log.Error("Authorization.Capture - Capture amount is greater than Auth amount")

		return nil, errors.New("Capture failure - Cannot capture amount that exceeds authorization's availability.")
	}

	if amount > auth.Balance() {
		log.Error("Authorization.Capture - Capture amount is greater than remaining Auth amount")

		return nil, errors.New("Capture failure - Cannot capture more than the remaining amount")
	}

	err := auth.CreditCard.Transaction("charge")
	if err != nil {
		log.WithField("err", err).Error("Authorization.Capture - Error trying to charge CC")
		
		return nil, err
	}

	newCapture := &Capture{
		Authorization: auth,
		Amount: amount,
	}

	auth.captures = append(auth.captures, newCapture)

	log.WithField("newCapture", newCapture).Debug("New Capture Successfully created")

	return newCapture, nil
}

func (auth *Authorization) Refund(amount float64, currency string) (*Refund, error) {
	auth.mu.Lock()
	defer auth.mu.Unlock()
	
	if auth.void {
		log.Error("Authorization.Refund - transaction is void")

		return nil, errors.New("Refund failure - Cannot refund on void transaction")
	}

	if amount > auth.TotalCapturedAmount() {
		log.Error("Authorization.Refund - Trying to refund more than total captured amount")

		return nil, errors.New("Refund failure - Cannot refund more than total captured amount")
	}

	err := auth.CreditCard.Transaction("refund")
	if err != nil {
		log.WithField("err", err).Error("Authorization.Refund - Error trying to charge CC")
		
		return nil, err
	}

	newRefund := &Refund{
		Authorization: auth,
		Amount: amount,
	}

	auth.refunds = append(auth.refunds, newRefund)

	log.WithField("newRefund", newRefund).Debug("New Refund Successfully created")

	return newRefund, nil
}

func (auth *Authorization) Balance() float64 {
	capturedAmount := 0.0
	for _, iterCapture := range auth.captures {
		capturedAmount += iterCapture.Amount
	}

	refundedAmount := 0.0
	for _, iterRefund := range auth.refunds {
		refundedAmount += iterRefund.Amount
	}

	return auth.Amount - capturedAmount + refundedAmount
}

func (auth *Authorization) TotalCapturedAmount() float64 {
	capturedAmount := 0.0
	for _, iterCapture := range auth.captures {
		capturedAmount += iterCapture.Amount
	}

	refundedAmount := 0.0
	for _, iterRefund := range auth.refunds {
		refundedAmount += iterRefund.Amount
	}

	return capturedAmount - refundedAmount
}