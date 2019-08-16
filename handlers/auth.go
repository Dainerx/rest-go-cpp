package handlers

import (
	"encoding/json"
	"net/http"
	"rest-go-cpp/models"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("BETELL_SECRET_KEY_2019_API")

// Create a struct to read the username and password from the request body
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	User *models.User
	jwt.StandardClaims
}

func authenificationTokenFromCookie(r *http.Request) (string, error) {
	// Get token string from client stored as cookie
	c, err := r.Cookie("token")
	if err != nil {
		return "", err
	} else {
		return c.Value, nil
	}
}

func Auth(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if ok, _, err := models.UserExists(creds.Email, creds.Password); !ok || err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(24 * time.Hour)
	// Get the user id
	_, user, _ := models.UserExists(creds.Email, creds.Password)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string based on our secret key
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the client cookie for "token" as the JWT we just generated with expiration time
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}

func Authenticated(w http.ResponseWriter, r *http.Request) (bool, *models.User) {
	tokenString, err := authenificationTokenFromCookie(r)
	if err != nil {
		return false, nil
	}
	claims := &Claims{}
	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	token, err := jwt.ParseWithClaims(tokenString, claims, func(tokenString *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return false, nil
		}
		return false, nil
	}
	if !token.Valid {
		return false, nil
	}
	return true, claims.User
}
