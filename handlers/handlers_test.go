package handlers

import (
		"net/http"
		"net/http/httptest"
		"testing"
		"encoding/json"
		"bytes"
		"io/ioutil"
		"errors"
		"github.com/stretchr/testify/assert"
		"github.com/stretchr/testify/mock"

		"github.com/nktsitas/checkout-techlab/gateway"
		"github.com/nktsitas/checkout-techlab/bank"
		"github.com/nktsitas/checkout-techlab/db"

		log "github.com/sirupsen/logrus"
)

// Mocks

type MockGateway struct {
	mock.Mock
}

func (m *MockGateway) NewAuthorization(req_body []byte, salt string) (*gateway.Authorization, error) {
	args := m.Called(req_body, salt)

	return args.Get(0).(*gateway.Authorization), args.Error(1)
}

func (m *MockGateway) GetSalt() string {
	args := m.Called()

	return args.String(0)
}

// ---

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

// ---

type MockAuthorization struct {
	mock.Mock
}

func (m *MockAuthorization) Void() error {
	args := m.Called()

	return args.Error(0)
}

func (m *MockAuthorization) Capture(amount float64, currency string) (*gateway.Capture, error) {
	args := m.Called(amount, currency)

	return args.Get(0).(*gateway.Capture), args.Error(1)
}

func (m *MockAuthorization) Refund(amount float64, currency string) (*gateway.Refund, error) {
	args := m.Called(amount, currency)

	return args.Get(0).(*gateway.Refund), args.Error(1)
}

func (m *MockAuthorization) GetCurrency() string {
	args := m.Called()

	return args.String(0)
}

// --- --- ---

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestPing(t *testing.T) {
	assert := assert.New(t)

	req, err := http.NewRequest("GET", "/ping", nil)
	assert.NoError(err, "Pong Response")

	w := httptest.NewRecorder()
	Ping(w, req)

	assert.Equal(200, w.Code, "Pong")
	assert.Equal("Pong!", w.Body.String(), "Pong")
}

// in handlers tests we mock out the gateway functionality (NewAuthorization, Capture, etc)
// as it is thoroughly tested in gateway package. Here we check that responses are as they should and errors are caught

func TestCreateAuthorizationHandler(t *testing.T) {
	assert := assert.New(t)

	testAuth := &gateway.Authorization{
		Amount: 100.00,
		Currency: "EUR",
		CreditCard: &bank.CreditCard{
			Number: "4000 0000 0000 0123",
			Expiry: "12/22",
			Cvv: "123",
		},
	}

	testResp := &authResponse{
		Id: "test",
		Amount: 100.00,
		Currency: "EUR",
	}

	testAuthJSON, _ := json.Marshal(testAuth)
	testRespJSON, _ := json.Marshal(testResp)

	tests := []struct{
		body []byte
		authCreated *gateway.Authorization
		err error
		expectedCode int
		expectedBody string
		description string
	}{
		{
			testAuthJSON,
			&gateway.Authorization{
				Id: "test",
				Amount: 100.00,
				Currency: "EUR",
			},
			nil,
			200,
			string(testRespJSON),
			"OK - Authorization Created",
		},
		{
			testAuthJSON,
			nil,
			errors.New("Something went wrong"),
			400,
			"Something went wrong\n",
			"Error",
		},
	}

	for _, iterTest := range tests {
		req, err := http.NewRequest("POST", "/authorize", bytes.NewBuffer(iterTest.body))
		assert.NoError(err)
	
		mockGateway := new(MockGateway)
		gateway.Gateway = mockGateway
	
		mockGateway.On("GetSalt").Return("test_salt", nil)
		mockGateway.On("NewAuthorization", testAuthJSON, "test_salt").Return(iterTest.authCreated, iterTest.err)
	
		mockGateway.MethodCalled("NewAuthorization", testAuthJSON, "test_salt")
		mockGateway.AssertNumberOfCalls(t, "NewAuthorization", 1)
	
		w := httptest.NewRecorder()
		CreateAuthorizationHandler(w, req)
	
		assert.Equal(iterTest.expectedCode, w.Code, iterTest.description)
		assert.Equal(iterTest.expectedBody, w.Body.String(), iterTest.description)
	}
}

func TestCaptureHandler(t *testing.T) {
	assert := assert.New(t)

	testAmount := 100.00

	testAuth := &gateway.Authorization{
		Id: "test",
		Amount: testAmount,
		Currency: "EUR",
		CreditCard: &bank.CreditCard{
			Number: "4000 0000 0000 0123",
			Expiry: "12/22",
			Cvv: "123",
		},
	}

	testCaptureRequest := &requestParams{
		Id: "test",
		Amount: testAmount,
	}

	testResp := &actionsResponse{
		Amount: testAmount,
		Currency: "EUR",
	}

	testCaptureRequestJSON, _ := json.Marshal(testCaptureRequest)
	testRespJSON, _ := json.Marshal(testResp)

	tests := []struct{
		body []byte
		authReturned *MockAuthorization
		captureCreated *gateway.Capture
		err error
		expectedCode int
		expectedBody string
		description string
	}{
		{
			testCaptureRequestJSON,
			new(MockAuthorization),
			&gateway.Capture{
				Authorization: testAuth,
				Amount: testAmount,
			},
			nil,
			200,
			string(testRespJSON),
			"OK - Capture Created",
		},
		{
			testCaptureRequestJSON,
			nil,
			nil,
			nil,
			400,
			"Wrong auth Id\n",
			"Error - Auth Id nil, wrong Id",
		},
		{
			testCaptureRequestJSON,
			new(MockAuthorization),
			nil,
			errors.New("Capture Error - Something went wrong"),
			400,
			"Capture Error - Something went wrong\n",
			"Error - Capture Error - Something went wrong",
		},
	}

	for _, iterTest := range tests {
		req, err := http.NewRequest("POST", "/capture", bytes.NewBuffer(iterTest.body))
		assert.NoError(err)

		mockAuth := iterTest.authReturned
		if mockAuth != nil {
			mockAuth.On("GetCurrency").Return("EUR")
			mockAuth.On("Capture", testAmount, "EUR").Return(iterTest.captureCreated, iterTest.err)

			mockAuth.MethodCalled("Capture", testAmount, "EUR")
			mockAuth.AssertNumberOfCalls(t, "Capture", 1)
		}
		
		testDB := new(MockDB)
		db.DB = testDB

		testDB.On("FetchItem", mock.Anything).Return(mockAuth)
	
		w := httptest.NewRecorder()
		CaptureHandler(w, req)
	
		assert.Equal(iterTest.expectedCode, w.Code, iterTest.description)
		assert.Equal(iterTest.expectedBody, w.Body.String(), iterTest.description)
	}
}

func TestRefundHandler(t *testing.T) {
	assert := assert.New(t)

	testAmount := 100.00

	testAuth := &gateway.Authorization{
		Id: "test",
		Amount: testAmount,
		Currency: "EUR",
		CreditCard: &bank.CreditCard{
			Number: "4000 0000 0000 0123",
			Expiry: "12/22",
			Cvv: "123",
		},
	}

	testRefundRequest := &requestParams{
		Id: "test",
		Amount: testAmount,
	}

	testResp := &actionsResponse{
		Amount: testAmount,
		Currency: "EUR",
	}

	testRefundRequestJSON, _ := json.Marshal(testRefundRequest)
	testRespJSON, _ := json.Marshal(testResp)

	tests := []struct{
		body []byte
		authReturned *MockAuthorization
		refundCreated *gateway.Refund
		err error
		expectedCode int
		expectedBody string
		description string
	}{
		{
			testRefundRequestJSON,
			new(MockAuthorization),
			&gateway.Refund{
				Authorization: testAuth,
				Amount: testAmount,
			},
			nil,
			200,
			string(testRespJSON),
			"OK - Refund Created",
		},
		{
			testRefundRequestJSON,
			nil,
			nil,
			nil,
			400,
			"Wrong auth Id\n",
			"Error - Auth Id nil, wrong Id",
		},
		{
			testRefundRequestJSON,
			new(MockAuthorization),
			nil,
			errors.New("Capture Error - Something went wrong"),
			400,
			"Capture Error - Something went wrong\n",
			"Error - Capture Error - Something went wrong",
		},
	}

	for _, iterTest := range tests {
		req, err := http.NewRequest("POST", "/refund", bytes.NewBuffer(iterTest.body))
		assert.NoError(err)

		mockAuth := iterTest.authReturned
		if mockAuth != nil {
			mockAuth.On("GetCurrency").Return("EUR")
			mockAuth.On("Refund", testAmount, "EUR").Return(iterTest.refundCreated, iterTest.err)

			mockAuth.MethodCalled("Refund", testAmount, "EUR")
			mockAuth.AssertNumberOfCalls(t, "Refund", 1)
		}
		
		testDB := new(MockDB)
		db.DB = testDB

		testDB.On("FetchItem", mock.Anything).Return(mockAuth)
	
		w := httptest.NewRecorder()
		RefundHandler(w, req)
	
		assert.Equal(iterTest.expectedCode, w.Code, iterTest.description)
		assert.Equal(iterTest.expectedBody, w.Body.String(), iterTest.description)
	}
}

func TestVoidHandler(t *testing.T) {
	assert := assert.New(t)

	testVoidRequest := &requestParams{
		Id: "test",
	}

	testResp := &actionsResponse{
		Amount: 0.0,
		Currency: "EUR",
	}

	testVoidRequestJSON, _ := json.Marshal(testVoidRequest)
	testRespJSON, _ := json.Marshal(testResp)

	tests := []struct{
		body []byte
		authReturned *MockAuthorization
		err error
		expectedCode int
		expectedBody string
		description string
	}{
		{
			testVoidRequestJSON,
			new(MockAuthorization),
			nil,
			200,
			string(testRespJSON),
			"OK - Refund Created",
		},
		{
			testVoidRequestJSON,
			nil,
			nil,
			400,
			"Wrong auth Id\n",
			"Error - Auth Id nil, wrong Id",
		},
		{
			testVoidRequestJSON,
			new(MockAuthorization),
			errors.New("Void Error - Something went wrong"),
			400,
			"Void Error - Something went wrong\n",
			"Error - Void Error - Something went wrong",
		},
	}

	for _, iterTest := range tests {
		req, err := http.NewRequest("POST", "/void", bytes.NewBuffer(iterTest.body))
		assert.NoError(err)

		mockAuth := iterTest.authReturned
		if mockAuth != nil {
			mockAuth.On("GetCurrency").Return("EUR")
			mockAuth.On("Void").Return(iterTest.err)

			mockAuth.MethodCalled("Void")
			mockAuth.AssertNumberOfCalls(t, "Void", 1)
		}
		
		testDB := new(MockDB)
		db.DB = testDB

		testDB.On("FetchItem", mock.Anything).Return(mockAuth)
	
		w := httptest.NewRecorder()
		VoidHandler(w, req)
	
		assert.Equal(iterTest.expectedCode, w.Code, iterTest.description)
		assert.Equal(iterTest.expectedBody, w.Body.String(), iterTest.description)
	}
}
