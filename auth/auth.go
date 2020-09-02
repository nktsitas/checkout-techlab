package auth

import (
	// "fmt"
	"net/http"
	"time"
	"os"
	"errors"
	"encoding/json"
	log "github.com/sirupsen/logrus"

	"github.com/dgrijalva/jwt-go"

)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type tokenResponse struct {
	AccessToken string `json:"access_token" example:"generated.jwt.token"`
}

type User struct {
	Username string
	Password string
}

var dummyUser User = User{"Checkout", "Checkout"}

// Refund godoc
// @Summary Logins a user and provides an authentication token
// @Description Logins a user and provides an authentication token
// @Tags status
// @Accept  json
// @Produce  json
// @Param Credentials body loginRequest true "User Credentials"
// @Success 200 {object} tokenResponse
// @Router /login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.WithField("err", err).Error("Login - Error reading body")
		http.Error(w, "Can't read body", http.StatusUnprocessableEntity)
		return
	}

	if dummyUser.Username != req.Username || dummyUser.Password != req.Password {
		log.Error("Login - Wrong Username or Password")
		http.Error(w, "Wrong Username or Password", http.StatusUnauthorized)
		return
	}

	token, err := GenerateToken(req.Username)
	if err != nil {
		log.WithField("err", err).Error("Login - Error Generating Token")
		http.Error(w, "Error Generating Token", http.StatusInternalServerError)
		return
	}

	resp := tokenResponse{
		AccessToken: token,
	}
	// tokenResponse.AccessToken = token
	respJSON, err := json.Marshal(resp)

	if err != nil {
		log.WithField("err", err).Error("Login - Error Marshaling Token Response")
		http.Error(w, "Error Marshaling Token Response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
  w.Write([]byte(respJSON))
}

func GenerateToken(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["client"] = username
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	accessSecret := os.Getenv("ACCESS_SECRET")
	
	tokenString, err := token.SignedString([]byte(accessSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func Authenticate(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, errors.New("Authentication Error")
				}
				return []byte(os.Getenv("ACCESS_SECRET")), nil
			})

			if err != nil {
				log.WithField("err", err).Error("Authenticate - Authentication Error")
				// hide too much info
				http.Error(w, "Authenticate - Authentication Error", http.StatusBadRequest)
				return
			}

			if token.Valid {
				inner.ServeHTTP(w, r)
			} else {
				log.Error("Authenticate - Authentication Error")
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		} else {
			log.Error("Authenticate - No Token Header provided")
			http.Error(w, "Authenticate - No Token Header provided", http.StatusBadRequest)
			return
		}
	})
}
