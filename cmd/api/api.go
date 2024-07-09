package main

import (
	"hng-stage2/internal/database"
	"hng-stage2/internal/repository"
	"hng-stage2/internal/server/controllers"
	"hng-stage2/internal/server/middleware"
	"hng-stage2/internal/service"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
}

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		addr: listenAddr,
	}
}

func (s *APIServer) Run() {

	database.Connect()

	userRepository := repository.NewUserSQLRepository(database.Database)
	organisationRepository := repository.NewOrganisationSQLRepository(database.Database)

	userService := service.NewUserService(userRepository, organisationRepository)
	orgService := service.NewOrganisationService(organisationRepository)

	authController := controllers.NewAuthController(*userService)
	organisationController := controllers.NewOrganisationController(*orgService)

	router := mux.NewRouter()

	log.Println("JSON API server running on port: ", s.addr)

	publicRouter := router.PathPrefix("/auth").Subrouter()

	publicRouter.HandleFunc("/register", authController.Register).Methods("POST")
	publicRouter.HandleFunc("/login", authController.Login).Methods("POST")

	protectedRouter := router.PathPrefix("/api").Subrouter()

	protectedRouter.Use(middleware.JWTAuthMiddleware)
	protectedRouter.HandleFunc("/users/{id}", organisationController.GetUserByID).Methods("GET")
	protectedRouter.HandleFunc("/organisations", organisationController.GetUserOrganisations).Methods("GET")
	protectedRouter.HandleFunc("/organisations/{orgId}", organisationController.GetOrgByID).Methods("GET")
	protectedRouter.HandleFunc("/organisations", organisationController.CreateOrganisation).Methods("POST")
	protectedRouter.HandleFunc("/organisations/{orgId}/users", organisationController.AddUserToOrganisation).Methods("POST")

	log.Fatal(http.ListenAndServe(s.addr, router))
}
