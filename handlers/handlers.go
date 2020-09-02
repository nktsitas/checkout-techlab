package handlers

import (
	"net/http"
	"io"
	// "time"
	"io/ioutil"
	// "strconv"
	"reflect"
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"github.com/nktsitas/checkout-techlab/gateway"
	"github.com/nktsitas/checkout-techlab/db"
)

type requestParams struct {
	Id string `json:"id" example:"unique_authorization_id"`
	Amount float64 `json:"amount" example:"100.00"`
}

type voidRequestParams struct {
	Id string `json:"id" example:"unique_authorization_id"`
}

type authResponse struct {
	Id string `json:"id" example:"unique_authorization_id"`
	Amount float64 `json:"amount" example:"100.00"`
	Currency string `json:"currency" example:"EUR"`
}

type actionsResponse struct {
	Amount float64 `json:"amount" example:"100.00"`
	Currency string `json:"currency" example:"EUR"`
}

// --- --- ---

// Ping godoc
// @Summary Get a server status update
// @Description Get a server status update
// @Tags status
// @Accept  json
// @Produce  json
// @Success 200 {string} Pong!
// @Router /status/ping [get]
func Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")

	io.WriteString(w, `Pong!`)
}

// CreateAuthorization godoc
// @Summary Creates a new authorization
// @Description Creates a new authorization
// @Tags status
// @Accept  json
// @Produce  json
// @Param authorization body gateway.Authorization true "Create authorization"
// @Param Token header string true "generated.jwt.token"
// @Success 200 {object} authResponse
// @Router /authorize [post]
func CreateAuthorizationHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
			log.WithField("err", err).Error("CreateAuthorizationHandler - Error reading body")
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
	}

	salt := gateway.Gateway.GetSalt()
	
	auth, err := gateway.Gateway.NewAuthorization(body, salt)
	if err != nil {
		log.WithField("err", err).Error("CreateAuthorizationHandler - Error Creating Authorization")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := &authResponse{
		Id: auth.Id,
		Amount: auth.Amount,
		Currency: auth.GetCurrency(),
	}

	writeResponse(w, resp)
}

// Capture godoc
// @Summary Captures amount from authorization
// @Description Captures amount from authorization
// @Tags status
// @Accept  json
// @Produce  json
// @Param captureRequest body requestParams true "Capture Amount"
// @Param Token header string true "generated.jwt.token"
// @Success 200 {object} actionsResponse
// @Router /capture [post]
func CaptureHandler(w http.ResponseWriter, r *http.Request) {
	var req requestParams
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
			log.WithField("err", err).Error("CaptureHandler - Error reading body")
			http.Error(w, "Can't read body", http.StatusBadRequest)
			return
	}

	authI := db.DB.FetchItem(req.Id)

	if authI == nil || reflect.ValueOf(authI).IsNil() {
		log.WithField("id", req.Id).Error("CaptureHandler - Wrong auth Id")
		http.Error(w, "Wrong auth Id", http.StatusBadRequest)
		return
	}

	auth := authI.(gateway.AuthorizationI)
	capture, err := auth.Capture(req.Amount, auth.GetCurrency())
	if err != nil {
		log.WithField("err", err).Error("CaptureHandler - Error in Capture")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := &actionsResponse{
		Amount: capture.Amount,
		Currency: auth.GetCurrency(),
	}

	writeResponse(w, resp)
}

// Void godoc
// @Summary Voids a transaction without charging the user
// @Description Voids a transaction without charging the user
// @Tags status
// @Accept  json
// @Produce  json
// @Param voidRequest body voidRequestParams true "Refund Amount"
// @Param Token header string true "generated.jwt.token"
// @Success 200 {object} actionsResponse
// @Router /void [post]
func VoidHandler(w http.ResponseWriter, r *http.Request) {
	var req voidRequestParams
	var err error
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
			log.WithField("err", err).Error("VoidHandler - Error reading body")
			http.Error(w, "Can't read body", http.StatusBadRequest)
			return
	}

	authI := db.DB.FetchItem(req.Id)

	if authI == nil || reflect.ValueOf(authI).IsNil() {
		log.WithField("id", req.Id).Error("VoidHandler - Wrong auth Id")
		http.Error(w, "Wrong auth Id", http.StatusBadRequest)
		return
	}

	auth := authI.(gateway.AuthorizationI)
	err = auth.Void()
	if err != nil {
		log.WithField("err", err).Error("VoidHandler - Error executing void")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := &actionsResponse{
		Amount: 0,
		Currency: auth.GetCurrency(),
	}

	writeResponse(w, resp)
}

// Refund godoc
// @Summary Refunds a previously captured amount from authorization
// @Description Refunds a previously captured amount from authorization
// @Tags status
// @Accept  json
// @Produce  json
// @Param refundRequest body requestParams true "Refund Amount"
// @Param Token header string true "generated.jwt.token"
// @Success 200 {object} actionsResponse
// @Router /refund [post]
func RefundHandler(w http.ResponseWriter, r *http.Request) {
	var req requestParams

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
			log.WithField("err", err).Error("RefundHandler - Error reading body")
			http.Error(w, "Can't read body", http.StatusBadRequest)
			return
	}

	authI := db.DB.FetchItem(req.Id)

	if authI == nil || reflect.ValueOf(authI).IsNil() {
		log.WithField("id", req.Id).Error("RefundHandler - Wrong auth Id")
		http.Error(w, "Wrong auth Id", http.StatusBadRequest)
		return
	}

	auth := authI.(gateway.AuthorizationI)
	refund, err := auth.Refund(req.Amount, auth.GetCurrency())
	if err != nil {
		log.WithField("err", err).Error("RefundHandler - Error executing refund")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := &actionsResponse{
		Amount: refund.Amount,
		Currency: auth.GetCurrency(),
	}

	writeResponse(w, resp)
}

func writeResponse(w http.ResponseWriter, resp interface{}) {
	respJSON, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
  w.Write(respJSON)
}
