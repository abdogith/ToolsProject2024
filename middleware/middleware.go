package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Claims defines the structure of the JWT claims, including user information like ID and role.
type Claims struct {
	UserID int    `json:"userID"` // user identifier
	Role   string `json:"role"`   // user role (e.g., admin, user, courier)
	jwt.StandardClaims
}

// JwtKey is the secret key used for signing JWT tokens. Ensure this is kept safe in production.
var JwtKey = []byte("your_secret_key") // replace with a more secure key in production

// GenerateToken creates a JWT token with a userID and role, and sets an expiration time of 24 hours.
func GenerateToken(userID int, role string) (string, error) {
	// setting the expiration time for the token to 24 hours from now
	expirationTime := time.Now().Add(24 * time.Hour)

	// creating the claims (payload) for the JWT token
	claims := &Claims{
		UserID: userID,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), // setting the expiration time
		},
	}

	// creating a new token with the claims and signing it using the HS256 method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// signing the token with the secret key and returning the signed token
	return token.SignedString(JwtKey)
}

// AuthMiddleware is a function that validates the JWT token passed in the Authorization header.
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// extracting the token from the Authorization header
		tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

		// creating an empty claims object
		claims := &Claims{}

		// parsing the token and verifying its validity
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// returning the secret key to verify the token
			return JwtKey, nil
		})

		// if there's an error or the token is invalid, return an Unauthorized error
		if err != nil || !token.Valid {
			http.Error(w, "unauthorized access", http.StatusUnauthorized)
			return
		}

		// if the token is valid, call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
