package handlers

import (
	"encoding/json"
	"hng-stage2/helpers"
	"hng-stage2/models"
	"net/http"

	"github.com/gorilla/mux"
)

func HandleCreateOrg(w http.ResponseWriter, r *http.Request) {
	var org models.Organisation

	if err := json.NewDecoder(r.Body).Decode(&org); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := helpers.Validate.Struct(org); err != nil {
		helpers.RespondWithValidationError(w, err)
		return

	}

	err := models.CreateOrg(&org)

	if err != nil {
		http.Error(w, "Client error", http.StatusInternalServerError)
		return
	}

	resp := H{
		"status":  "success",
		"message": "Organisation created successfully",
		"data":    org,
	}

	helpers.WriteJSON(w, http.StatusCreated, resp)
}

func HandleGetUserByID(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]
	user, err := models.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "User retrieved successfully",
		"data":    user,
	}
	helpers.WriteJSON(w, http.StatusOK, response)
}

func HandleGetUserOrganisations(w http.ResponseWriter, r *http.Request) {
	claims, err := helpers.ValidateJWTFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	orgs, err := models.GetUserOrgsByID(claims.UserID)
	if err != nil {
		http.Error(w, "Organisations not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Organisations retrieved successfully",
		"data":    orgs,
	}
	helpers.WriteJSON(w, http.StatusOK, response)
}

func HandleGetOrgByID(w http.ResponseWriter, r *http.Request) {
	orgID := mux.Vars(r)["orgId"]
	org, err := models.GetByOrgID(orgID)
	if err != nil {
		http.Error(w, "Organisation not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Organisation retrieved successfully",
		"data":    org,
	}
	helpers.WriteJSON(w, http.StatusOK, response)
}

func HandleAddUserToOrganisation(w http.ResponseWriter, r *http.Request) {
	orgID := mux.Vars(r)["orgId"]
	var req struct {
		UserID string `json:"userId" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := helpers.Validate.Struct(req); err != nil {
		helpers.RespondWithValidationError(w, err)
		return
	}

	err := models.AddUserToOrganisation(orgID, req.UserID)
	if err != nil {
		http.Error(w, "Failed to add user to organisation", http.StatusInternalServerError)
		return
	}

	response := H{
		"status":  "success",
		"message": "User added to organisation successfully",
	}
	helpers.WriteJSON(w, http.StatusOK, response)
}
