package handlers

import (
	"encoding/json"
	"hng-stage2/database"
	"hng-stage2/helpers"
	"hng-stage2/models"
	"net/http"

	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/crypto/bcrypt"
)

type H map[string]any

type RegisterParams struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	Phone     string `json:"phone"`
}

type LoginParams struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthSuccessResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    H      `json:"data"`
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {

	var params RegisterParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Bad request", "Invalid request payload", 400)
		return
	}

	if err := helpers.Validate.Struct(params); err != nil {
		helpers.RespondWithValidationError(w, err)
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)

	user := models.User{
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Password:  string(hashedPassword),
		Email:     params.Email,
		Phone:     params.Phone,
	}

	err := database.Database.Create(&user).Error

	if err != nil {
		helpers.WriteJSON(w, http.StatusUnauthorized, helpers.ErrorResponse{
			Status:     "Bad Request",
			StatusCode: http.StatusUnauthorized,
			Message:    "Authentication failed",
		})
		return
	}

	token, err := helpers.GenerateJWT(user)

	if err != nil {
		helpers.WriteJSON(w, http.StatusUnauthorized, helpers.ErrorResponse{
			Status:     "Bad Request",
			StatusCode: http.StatusUnauthorized,
			Message:    "Authentication failed",
		})
		return
	}

	helpers.WriteJSON(w, http.StatusCreated, AuthSuccessResponse{
		Status:  "success",
		Message: "Registration successful",
		Data: H{
			"accessToken": token,
			"user":        user,
		},
	})

}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var params *LoginParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {

		return
	}

	if err := helpers.Validate.Struct(params); err != nil {
		helpers.RespondWithValidationError(w, err)
		return
	}
	user, err := models.FindUserByEmail(params.Email)

	if err != nil {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	err = user.ValidatePassword(params.Password)

	if err != nil {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	token, err := helpers.GenerateJWT(*user)
	if err != nil {
		http.Error(w, "Authentication failed", http.StatusInternalServerError)

		return
	}
	helpers.WriteJSON(w, http.StatusOK, AuthSuccessResponse{
		Status:  "success",
		Message: "Registration successful",
		Data: H{
			"accessToken": token,
			"user":        user,
		},
	})

}
