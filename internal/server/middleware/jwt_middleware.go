package middleware

import (
	"fmt"
	"hng-stage2/helpers"
	"net/http"
	"os"
	"strings"

	jwt "github.com/golang-jwt/jwt/v4"
	_ "github.com/joho/godotenv/autoload"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization")

		tokenString := strings.Split(tokenHeader, " ")[1]
		if tokenString == "" {
			http.Error(w, "Missing or invalid token", http.StatusUnauthorized)
			return
		}

		tkn, err := helpers.ValidateJWT(tokenString)

		if err != nil {
			helpers.RespondWithError(w, http.StatusUnauthorized, helpers.ErrorResponse{
				Status:     http.StatusText(http.StatusUnauthorized),
				Message:    "Invalid token",
				StatusCode: http.StatusUnauthorized,
			})
			return
		}

		if err := tkn.Valid(); err != nil {
			helpers.RespondWithError(w, http.StatusUnauthorized, helpers.ErrorResponse{
				Status:     http.StatusText(http.StatusUnauthorized),
				Message:    "Invalid token",
				StatusCode: http.StatusUnauthorized,
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
}
