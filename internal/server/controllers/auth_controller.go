package controllers

import (
	"encoding/json"
	"hng-stage2/helpers"
	"hng-stage2/internal/models"
	"hng-stage2/internal/service"
	"net/http"

	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
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

type SuccessResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    H      `json:"data"`
}

type AuthController struct {
	userService service.UserService
	validate    *validator.Validate
}

func NewAuthController(userService service.UserService) *AuthController {
	return &AuthController{userService: userService, validate: validator.New()}
}

func (controller *AuthController) Register(w http.ResponseWriter, r *http.Request) {

	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, helpers.ErrorResponse{
			Status:     "Bad Request",
			Message:    "Invalid request payload",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	if err := helpers.Validate.Struct(user); err != nil {
		helpers.RespondWithValidationError(w, err)
		return
	}

	existingUser, _ := controller.userService.GetUserByEmail(user.Email)

	// if err != nil {

	// 	if err != gorm.ErrRecordNotFound {
	// 		panic(err)
	// 	}
	// }
	// 	// if err == gorm.ErrRecordNotFound {

	// 	// 	helpers.RespondWithError(w, http.StatusInternalServerError, helpers.ErrorResponse{
	// 	// 		Status:     "Bad Request",
	// 	// 		Message:    "Client error",
	// 	// 		StatusCode: http.StatusInternalServerError,
	// 	// 	})
	// 	// 	return
	// 	// } else if err != gorm.ErrRecordNotFound {

	// 	// 	helpers.RespondWithError(w, http.StatusBadRequest, helpers.ErrorResponse{
	// 	// 		Status:     "Bad Request",
	// 	// 		Message:    "Email already exists",
	// 	// 		StatusCode: http.StatusBadRequest,
	// 	// 	})
	// 	// }

	// }

	if existingUser != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, helpers.ErrorResponse{
			Status:     "Bad Request",
			Message:    "Email already exists",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	newUser, err := controller.userService.Register(&user)

	if err != nil {

		helpers.RespondWithError(w, http.StatusBadRequest, helpers.ErrorResponse{
			Status:     "Bad Request",
			Message:    "Registration unsuccessful",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	token, err := helpers.GenerateJWT(*newUser)

	if err != nil {
		helpers.WriteJSON(w, http.StatusUnauthorized, helpers.ErrorResponse{
			Status:     "Bad Request",
			StatusCode: http.StatusUnauthorized,
			Message:    "Authentication failed",
		})
		return
	}

	userResponse := H{
		"userId":    newUser.UserID,
		"firstName": newUser.FirstName,
		"lastName":  newUser.LastName,
		"email":     newUser.Email,
	}

	helpers.WriteJSON(w, http.StatusCreated, SuccessResponse{
		Status:  "success",
		Message: "Registration successful",
		Data: H{
			"accessToken": token,
			"user":        userResponse,
		},
	})

}

func (controller AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var params LoginParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {

		return
	}

	if err := helpers.Validate.Struct(&params); err != nil {
		helpers.RespondWithValidationError(w, err)
		return
	}
	user, err := controller.userService.Login(params.Email, params.Password)

	if err != nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, helpers.ErrorResponse{
			Status:     "Bad request",
			Message:    "Authentication failed",
			StatusCode: http.StatusUnauthorized,
		})
		return
	}

	token, err := helpers.GenerateJWT(*user)
	if err != nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, helpers.ErrorResponse{
			Status:     "Bad request",
			Message:    "Authentication failed",
			StatusCode: http.StatusUnauthorized,
		})
		return
	}

	userResponse := H{
		"userId":    user.UserID,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"email":     user.Email,
	}

	helpers.WriteJSON(w, http.StatusOK, SuccessResponse{
		Status:  "success",
		Message: "Registration successful",
		Data: H{
			"accessToken": token,
			"user":        userResponse,
		},
	})

}
