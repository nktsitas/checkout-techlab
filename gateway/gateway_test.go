package gateway

import (
	"testing"
	"errors"
	"io/ioutil"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	
	"github.com/nktsitas/checkout-techlab/db"
	"github.com/nktsitas/checkout-techlab/bank"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testAuthorizations map[string]*Authorization
var testAuthorizationStrings map[string][]byte

type MockDB struct {
	mock.Mock
}

func (m *MockDB) StoreItem(id string, item interface{}) {
	m.Called()
}

func (m *MockDB) FetchItem(id string) interface{} {
	args := m.Called(id)
	return args.Get(0)
}

func (m *MockDB) DeleteItem(id string) {
	m.Called()
}

// --- --- ---

// We predefine a series of Auth objects in order to observe the correct behavior when
// combinations of Capture,Void,Refund calls are made on (NewAuthorization) created authorizations

func init() {
	log.SetOutput(ioutil.Discard)

	testCreditCards := map[string]*bank.CreditCard{
		"OK": &bank.CreditCard{
			Number: "4000 0000 0000 0123",
			Expiry: "12/22",
			Cvv: "123",
		},
		"OK_void1": &bank.CreditCard{
			Number: "4000 0000 0000 0123",
			Expiry: "12/22",
			Cvv: "123",
		},
		"OK_void2": &bank.CreditCard{
			Number: "4000 0000 0000 0123",
			Expiry: "12/22",
			Cvv: "123",
		},
		"OK_refund": &bank.CreditCard{
			Number: "4000 0000 0000 0123",
			Expiry: "12/22",
			Cvv: "123",
		},
		"void": &bank.CreditCard{
			Number: "4000 0000 0000 0123",
			Expiry: "12/22",
			Cvv: "123",
		},
		"AuthFailure": &bank.CreditCard{
			Number: "4000 0000 0000 0119",
			Expiry: "12/22",
			Cvv: "123",
		},
		"CaptureFailure": &bank.CreditCard{
			Number: "4000 0000 0000 0259",
			Expiry: "12/22",
			Cvv: "123",
		},
		"RefundFailure": &bank.CreditCard{
			Number: "4000 0000 0000 3238",
			Expiry: "12/22",
			Cvv: "123",
		},
		"NoCvv": &bank.CreditCard{
			Number: "4000 0000 0000 0123",
			Expiry: "12/22",
		},
		"NoExp": &bank.CreditCard{
			Number: "4000 0000 0000 0123",
			Cvv: "123",
		},
		"NoNumber": &bank.CreditCard{
			Expiry: "12/22",
			Cvv: "123",
		},
		"InvalidCvv": &bank.CreditCard{
			Number: "4000 0000 0000 0123",
			Expiry: "12/22",
			Cvv: "aaa",
		},
		"InvalidExpiry": &bank.CreditCard{
			Number: "4000 0000 0000 0123",
			Expiry: "aaa",
			Cvv: "123",
		},
		"InvalidNumber": &bank.CreditCard{
			Number: "aaa",
			Expiry: "12/22",
			Cvv: "123",
		},
		"InvalidNumber2": &bank.CreditCard{
			Number: "4917 4845 8989 7107",
			Expiry: "01/23",
			Cvv: "123",
		},
	}

	testAuthorizations = make(map[string]*Authorization)
	testAuthorizationStrings = make(map[string][]byte)

	for key, iterCreditCard := range testCreditCards {
		auth, authJSON := getNewTestAuth(iterCreditCard)

		if key == "void" {
			auth.void = true
		}

		testAuthorizations[key] = auth
		testAuthorizationStrings[key] = authJSON
	}
}

func getNewTestAuth(cc *bank.CreditCard) (*Authorization, []byte) {
	auth := Authorization{
		Amount: 200.00,
		Currency: "EUR",
		CreditCard: cc,
	}

	authJSON, err := json.Marshal(auth)
	if err != nil {
		log.Fatal(err)
	}

	return &auth, authJSON
}

func getNewTestCapture(auth *Authorization, amount float64) *Capture {
	return &Capture{
		Authorization: auth,
		Amount: amount,
	}
}

func getNewTestRefund(auth *Authorization, amount float64) *Refund {
	return &Refund{
		Authorization: auth,
		Amount: amount,
	}
}

// --- --- ---
	
func TestNewAuthorization(t *testing.T) {
	assert := assert.New(t)

	tests := []struct{
		input []byte
		expected *Authorization
		err error
		description string
	}{
		{
			testAuthorizationStrings["OK"],
			testAuthorizations["OK"],
			nil,
			"OK - Create authorization",
		},
		{
			nil,
			nil,
			errors.New("Error Unmarshaling JSON - unexpected end of JSON input"),
			"Error - No Body",
		},
		{
			testAuthorizationStrings["AuthFailure"],
			nil,
			errors.New("Authorization failure - Unknown Error"),
			"Error - Manually triggered authorization failure.",
		},
		{
			testAuthorizationStrings["NoCvv"],
			nil,
			errors.New("Invalid CreditCard - No Cvv provided"),
			"Error - No Cvv Provided",
		},
		{
			testAuthorizationStrings["NoExp"],
			nil,
			errors.New("Invalid CreditCard - No Expiry provided"),
			"Error - No Expiry Provided",
		},
		{
			testAuthorizationStrings["NoNumber"],
			nil,
			errors.New("Invalid CreditCard - No Number provided"),
			"Error - No Number Provided",
		},
	}

	for _, iterTest := range tests {
		testDB := new(MockDB)
		db.DB = testDB
		testDB.On("StoreItem", mock.Anything).Return()

		testGateway := new(GatewayS)

		auth, err := testGateway.NewAuthorization(iterTest.input, "salt")

		if iterTest.expected != nil {
			iterTest.expected.Id = generateID(iterTest.input, "salt")	
			testDB.AssertNumberOfCalls(t, "StoreItem", 1)	
		}

		assert.Equal(iterTest.expected, auth, iterTest.description)
		assert.Equal(iterTest.err, err, iterTest.description)
	}
}

// ---

func TestVoid(t *testing.T) {
	assert := assert.New(t)

	testAuthorizations["OK_void2"].Capture(10.00, "EUR")

	tests := []struct{
		input *Authorization
		expected *Authorization
		err error
		description string
	}{
		{
			testAuthorizations["OK_void1"],
			testAuthorizations["void"],
			nil,
			"OK - Void a transaction",
		},
		{
			testAuthorizations["void"],
			testAuthorizations["void"],
			errors.New("Void Failure - Transaction already void"),
			"Error - Try void a void transaction",
		},
		{
			testAuthorizations["OK_void2"],
			testAuthorizations["OK_void2"],
			errors.New("Void Failure - Cannot void transaction with captured amount"),
			"Error - Try void a transaction with captures",
		},
	}

	for _, iterTest := range tests {
		log.Info(iterTest.input.void)
		err := iterTest.input.Void()
		log.Info(iterTest.input.void)

		assert.Equal(iterTest.expected.void, iterTest.input.void, iterTest.description)
		assert.Equal(iterTest.err, err, iterTest.description)
	}
}

func TestCapture(t *testing.T) {
	assert := assert.New(t)

	testCapture := getNewTestCapture(testAuthorizations["OK"], 100.00)

	tests := []struct{
		authorization *Authorization
		amount float64
		currency string
		expected *Capture
		err error
		description string
	}{
		{
			testAuthorizations["OK"],
			100.00,
			"EUR",
			testCapture,
			nil,
			"OK - Capture Created for amount",
		},
		{
			testAuthorizations["OK"],
			300.00,
			"EUR",
			nil,
			errors.New("Capture failure - Cannot capture amount that exceeds authorization's availability."),
			"Error - Try Capture more than amount",
		},
		{
			testAuthorizations["CaptureFailure"],
			50.00,
			"EUR",
			nil,
			errors.New("Capture failure - Unknown Error"),
			"Error - Manually triggered capture failure.",
		},
		{
			testAuthorizations["OK"],
			150.00,
			"EUR",
			nil,
			errors.New("Capture failure - Cannot capture more than the remaining amount"),
			"Error - Try Capture more than remaining amount",
		},
		{
			testAuthorizations["void"],
			10.00,
			"EUR",
			nil,
			errors.New("Capture failure - Cannot capture on void transaction"),
			"Error - Try Capture on void transaction",
		},
	}

	for _, iterTest := range tests {
		capture, err := iterTest.authorization.Capture(iterTest.amount, iterTest.currency)

		assert.Equal(iterTest.expected, capture, iterTest.description)
		assert.Equal(iterTest.err, err, iterTest.description)
	}
}

func TestRefund(t *testing.T) {
	assert := assert.New(t)

	testAuthorizations["OK_refund"].Capture(10.00, "EUR")
	testAuthorizations["RefundFailure"].Capture(100.00, "EUR")

	testRefund := getNewTestRefund(testAuthorizations["OK_refund"], 5.00)

	tests := []struct{
		authorization *Authorization
		amount float64
		currency string
		expected *Refund
		err error
		description string
	}{
		{
			testAuthorizations["OK_refund"],
			5.00,
			"EUR",
			testRefund,
			nil,
			"OK - Refund an amount from what is captured",
		},
		{
			testAuthorizations["OK_refund"],
			15.00,
			"EUR",
			nil,
			errors.New("Refund failure - Cannot refund more than total captured amount"),
			"Error - Try refund more than total captured amount",
		},
		{
			testAuthorizations["void"],
			10.00,
			"EUR",
			nil,
			errors.New("Refund failure - Cannot refund on void transaction"),
			"Error - Try Capture on void transaction",
		},
		{
			testAuthorizations["RefundFailure"],
			50.00,
			"EUR",
			nil,
			errors.New("Refund failure - Unknown Error"),
			"Error - Manually triggered refund failure.",
		},
	}
	
	for _, iterTest := range tests {

		refund, err := iterTest.authorization.Refund(iterTest.amount, iterTest.currency)

		assert.Equal(iterTest.expected, refund, iterTest.description)
		assert.Equal(iterTest.err, err, iterTest.description)
	}
}
