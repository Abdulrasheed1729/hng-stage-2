package handlers

import (
	"encoding/json"
	"fmt"
	"hng-stage2/database"
	"hng-stage2/models"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/badoux/checkmail"
	jwt "github.com/golang-jwt/jwt/v4"
)

type H map[string]any

type RegisterParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Phone     string `json:"phone"`
}

type LoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthSuccessResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    H      `json:"data"`
}

type AuthUnsuccessfulResponse struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

func HandleRegister(w http.ResponseWriter, r *http.Request) error {

	var params *RegisterParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return err
	}

	err := validateEmail(params.Email)

	if err != nil {
		WriteJSON(w, http.StatusUnprocessableEntity, H{"errors": []H{
			{
				"field":   "email",
				"message": "invalid email",
			},
		},
		})
		return err
	}

	user := models.User{
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Password:  params.Password,
		Email:     params.Email,
		Phone:     params.Phone,
	}

	err = validateEmail(params.Email)
	if err != nil {
		WriteJSON(w, http.StatusUnprocessableEntity, H{
			"errors": []H{
				{
					"field":   "email",
					"message": "invalid email",
				},
			},
		})

		return err
	}

	err = ValidateName(params.FirstName)
	if err != nil {
		WriteJSON(w, http.StatusUnprocessableEntity, H{
			"errors": []H{
				{
					"field":   "firstName",
					"message": "invalid name",
				},
			},
		})

		return err
	}

	err = ValidateName(params.LastName)
	if err != nil {
		WriteJSON(w, http.StatusUnprocessableEntity, H{
			"errors": []H{
				{
					"field":   "lastName",
					"message": "invalid name",
				},
			},
		})

		return err
	}

	err = user.ValidatePassword(params.Password)
	if err != nil {
		WriteJSON(w, http.StatusUnprocessableEntity, H{
			"errors": []H{
				{
					"field":   "password",
					"message": "invalid password",
				},
			},
		})

		return err
	}

	err = ValidatePhoneNumber(params.Phone)
	if err != nil {
		WriteJSON(w, http.StatusUnprocessableEntity, H{
			"errors": []H{
				{
					"field":   "phone",
					"message": "invalid phone number",
				},
			},
		})

		return err
	}

	err = database.Database.Create(&user).Error

	if err != nil {
		return err
	}

	token, _ := GenerateJWT(user)

	WriteJSON(w, http.StatusCreated, AuthSuccessResponse{
		Status:  "success",
		Message: "Registration successful",
		Data: H{
			"accessToken": token,
			"user":        user,
		},
	})

	return nil

}

func HandleLogin(w http.ResponseWriter, r *http.Request) error {

	if r.Method != "POST" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}
	var params *LoginParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return err
	}

	err := validateEmail(params.Email)
	if err != nil {
		WriteJSON(w, http.StatusUnprocessableEntity, H{
			"errors": []H{
				{
					"field":   "email",
					"message": "invalid email",
				},
			},
		})

		return err
	}

	user, err := models.FindUserByEmail(params.Email)

	if err != nil {
		WriteJSON(w, http.StatusUnprocessableEntity, H{
			"errors": []H{
				{
					"field":   "email",
					"message": "user does not exist",
				},
			},
		})

		return err
	}
	err = user.ValidatePassword(params.Password)

	if err != nil {
		WriteJSON(w, http.StatusUnprocessableEntity, H{
			"errors": []H{
				{
					"field":   "email",
					"message": "user does not exist",
				},
			},
		})
		return err
	}

	token, err := GenerateJWT(*user)

	if err != nil {
		WriteJSON(w, http.StatusUnauthorized, AuthUnsuccessfulResponse{
			Status:     "Bad request",
			Message:    "Authentication failed",
			StatusCode: http.StatusUnauthorized,
		})
		return err
	}

	WriteJSON(w, http.StatusOK, AuthSuccessResponse{
		Status:  "success",
		Message: "Registration successful",
		Data: H{
			"accessToken": token,
			"user":        user,
		},
	})

	return nil
}

func GenerateJWT(user models.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"id":  user.ID,
		"iat": time.Now().Unix(),
		"eat": time.Now().Add(time.Second * time.Duration(2000)).Unix(),
	})
	return token.SignedString(secret)
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

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func validateEmail(email string) error {
	return checkmail.ValidateFormat(email)
}

func ValidateName(name string) error {
	nameRegex := `^[a-zA-Z0-9_-]+$`
	_, err := regexp.MatchString(nameRegex, name)
	return err
}

func ValidatePhoneNumber(phoneNumber string) error {
	phoneNumberRegex := `^(?:(?:\(?(?:\+([1-9]\d{1,3}))?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?)?$`
	_, err := regexp.MatchString(phoneNumberRegex, phoneNumber)
	return err
}
