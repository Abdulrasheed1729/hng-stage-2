package controllers

import (
	"encoding/json"
	"hng-stage2/helpers"
	"hng-stage2/internal/models"
	"hng-stage2/internal/service"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/go-playground/validator/v10"
)

type OrganisationController struct {
	organisationService service.OrganisationService
	validate            *validator.Validate
}

func NewOrganisationController(organisationService service.OrganisationService) *OrganisationController {
	return &OrganisationController{organisationService: organisationService, validate: validator.New()}
}

func (controller OrganisationController) CreateOrganisation(w http.ResponseWriter, r *http.Request) {
	var org models.Organisation

	if err := json.NewDecoder(r.Body).Decode(&org); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, helpers.ErrorResponse{
			Status:     "Bad Request",
			Message:    "Invalid request payload",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	if err := helpers.Validate.Struct(org); err != nil {
		helpers.RespondWithValidationError(w, err)
		return

	}

	err := controller.organisationService.Create(&org)

	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, helpers.ErrorResponse{
			Status:     "Bad Request",
			Message:    "Client error",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	resp := H{
		"status":  "success",
		"message": "Organisation created successfully",
		"data":    org,
	}

	helpers.WriteJSON(w, http.StatusCreated, resp)
}

func (controller OrganisationController) GetUserByID(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]

	_, err := helpers.ValidateJWTFromRequest(r)

	if err != nil {

		helpers.RespondWithError(w, http.StatusUnauthorized, helpers.ErrorResponse{
			Status:     http.StatusText(http.StatusUnauthorized),
			Message:    "Invalid token",
			StatusCode: http.StatusUnauthorized,
		})
		return
	}

	// if userID != claims.UserID {
	// 	helpers.RespondWithError(w, http.StatusUnauthorized, helpers.ErrorResponse{
	// 		Status:     "Bad request",
	// 		Message:    "Client error",
	// 		StatusCode: http.StatusUnauthorized,
	// 	})
	// 	return
	// }

	user, err := controller.organisationService.GetUserByID(userID)
	if err != nil {
		helpers.RespondWithError(w, http.StatusNotFound, helpers.ErrorResponse{
			Status:     "Not found",
			Message:    "User not found",
			StatusCode: http.StatusNotFound,
		})
		return
	}

	userResponse := H{
		"userId":    user.UserID,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"email":     user.Email,
	}

	helpers.WriteJSON(w, http.StatusOK, userResponse)
}

func (controller OrganisationController) GetUserOrganisations(w http.ResponseWriter, r *http.Request) {
	claims, err := helpers.ValidateJWTFromRequest(r)
	if err != nil {
		// helpers.RespondWithError(w, http.StatusBadRequest, helpers.ErrorResponse{
		// 	Status:     "Bad request",
		// 	Message:    "Authentication failure",
		// 	StatusCode: http.StatusBadRequest,
		// })
		helpers.RespondWithError(w, http.StatusUnauthorized, helpers.ErrorResponse{
			Status:     http.StatusText(http.StatusUnauthorized),
			Message:    "Invalid token",
			StatusCode: http.StatusUnauthorized,
		})
		return
	}

	orgs, err := controller.organisationService.GetByUser(claims.UserID)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, helpers.ErrorResponse{
			Status:     "Bad request",
			Message:    "Authentication failure",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	response := H{
		"status":  "success",
		"message": "Organisations retrieved successfully",
		"data":    orgs,
	}
	helpers.WriteJSON(w, http.StatusOK, response)
}

func (controller OrganisationController) GetOrgByID(w http.ResponseWriter, r *http.Request) {
	orgID := mux.Vars(r)["orgId"]

	org, err := controller.organisationService.GetByID(orgID)

	if err != nil {
		helpers.RespondWithError(w, http.StatusNotFound, helpers.ErrorResponse{
			Status:     "Not found",
			Message:    "Organisation not found",
			StatusCode: http.StatusNotFound,
		})
		return
	}

	response := H{
		"status":  "success",
		"message": "Organisation retrieved successfully",
		"data":    org,
	}
	helpers.WriteJSON(w, http.StatusOK, response)
}

type AddUserToOrganisationParam struct {
	UserID string `json:"userId" validate:"required"`
}

func (controller OrganisationController) AddUserToOrganisation(w http.ResponseWriter, r *http.Request) {
	orgID := mux.Vars(r)["orgId"]

	var param AddUserToOrganisationParam

	if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := helpers.Validate.Struct(param); err != nil {
		helpers.RespondWithValidationError(w, err)
		return
	}

	err := controller.organisationService.AddUserToOrganisation(orgID, param.UserID)
	if err != nil {
		http.Error(w, "Failed to add user to organisation", http.StatusInternalServerError)
		helpers.RespondWithError(w, http.StatusNotFound, helpers.ErrorResponse{
			Status:     "Client error",
			Message:    "Failed to add user to organisation",
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	response := H{
		"status":  "success",
		"message": "User added to organisation successfully",
	}
	helpers.WriteJSON(w, http.StatusOK, response)
}
