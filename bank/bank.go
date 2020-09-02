package bank

import (
	"strconv"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"

)

type Transaction struct {
	message string
	err error
}

type CreditCard struct {
	Number string		 `json:"number"`
	Expiry string		 `json:"expiry"`
	Cvv string			 `json:"cvv"`

	transaction chan Transaction
}

// This can be either Charge or Refund, etc
func (cc *CreditCard) Transaction(action string) error {
	if cc.transaction == nil {
		cc.transaction = make(chan Transaction)
	}

	// communicate with CreditCard service and wait to receive response.
	// This can be an API call or a message broker receiver, etc
	go cc.simulateCreditCardTransaction(action)
	resp := <- cc.transaction

	if resp.err != nil {
		return resp.err
	}

	return nil
}

func (cc *CreditCard) simulateCreditCardTransaction(action string) {
	time.Sleep(200*time.Millisecond)

	if cc.Number == "4000 0000 0000 0259" && action == "charge" {
		log.Error("CreditCard.Charge - Manually triggered capture failure.")
		
		cc.transaction <- Transaction{
			"Capture Failure",
			errors.New("Capture failure - Unknown Error"),
		}
	} else if cc.Number == "4000 0000 0000 3238" && action == "refund" {
		log.Error("CreditCard.Charge - Manually triggered refund failure.")
		
		cc.transaction <- Transaction{
			"Capture Failure",
			errors.New("Refund failure - Unknown Error"),
		}
	} else {
		cc.transaction <- Transaction{
			"Capture Successful!",
			nil,
		}
	}
}

func (cc *CreditCard) Validate() error {
	if cc.Number == "" {
		return errors.New("Invalid CreditCard - No Number provided")
	}
	if cc.Expiry == "" {
		return errors.New("Invalid CreditCard - No Expiry provided")
	}
	if cc.Cvv == "" {
		return errors.New("Invalid CreditCard - No Cvv provided")
	}

	if _, err := strconv.Atoi(cc.Cvv); err != nil {
		return errors.New("Invalid CreditCard - Cvv is not valid")
	}

	return nil
}

// func (cc *CreditCard) Luhn() bool {
	// Out of time :)

	// Idea here would be to separate our credit card number into odd and even digits.
	// Starting from the right (after the check digit - the last one) we multiply the even digits by 2.
	// If the result is greater than 10, we add the digits together.
	// We add all the results -> A1

	// We add all odd digits -> A2

	// Finally we calculate: (10 - (A1 + A2) % 10) % 10 <-- this is on top of my head, might get improvements.
	// The idea is: what's the number that we need to add in order to have: A1 + A2 + number % 10 = 0
// }