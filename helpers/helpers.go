package helpers

import (
	"encoding/json"
	"errors"
	"hng-stage2/models"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"

	jwt "github.com/golang-jwt/jwt/v4"
)

var Validate *validator.Validate = validator.New()

type JWTClaims struct {
	Email  string `json:"email"`
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

func GenerateJWT(user models.User) (string, error) {

	secret := os.Getenv("JWT_SECRET")
	now := time.Now()
	claims := &JWTClaims{
		Email:  user.Email,
		UserID: user.ID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: now.Add(time.Hour * 24)},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secret)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

type ErrorResponse struct {
	Status     string  `json:"status,omitempty"`
	Message    string  `json:"message,omitempty"`
	StatusCode int     `json:"statusCode,omitempty"`
	Errors     []Error `json:"errors,omitempty"`
}

type Error struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func RespondWithValidationError(w http.ResponseWriter, err error) {
	var errors []Error
	for _, err := range err.(validator.ValidationErrors) {
		var element Error
		element.Field = err.StructNamespace()
		element.Message = err.Tag()
		errors = append(errors, element)
	}

	response := map[string]interface{}{
		"errors": errors,
	}
	WriteJSON(w, http.StatusUnprocessableEntity, response)
}

func RespondWithError(w http.ResponseWriter, status int, statusText, message string, statusCode int) {
	response := map[string]interface{}{
		"status":     statusText,
		"message":    message,
		"statusCode": statusCode,
	}

	WriteJSON(w, status, response)
}

func ValidateJWT(tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}

	jwtKey := os.Getenv("JWT_SECRET")

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("invalid token signature")
		}
		return nil, errors.New("invalid token")
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func ValidateJWTFromRequest(r *http.Request) (*JWTClaims, error) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return nil, errors.New("missing or invalid token")
	}

	return ValidateJWT(tokenString)
}

func FormatValidationError(err error) []Error {
	var errors []Error
	for _, err := range err.(validator.ValidationErrors) {
		var element Error
		element.Field = err.Field()
		element.Message = "Required Input"
		errors = append(errors, element)
	}
	return errors
}
