package api

import (
	"hng-stage2/handlers"
	"hng-stage2/middleware"
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
	router := mux.NewRouter()

	log.Println("JSON API server running on port: ", s.addr)

	publicRouter := router.PathPrefix("/auth").Subrouter()

	publicRouter.HandleFunc("/register", handlers.HandleRegister).Methods("POST")
	publicRouter.HandleFunc("/login", handlers.HandleLogin).Methods("POST")

	protectedRouter := router.PathPrefix("/api").Subrouter()

	protectedRouter.Use(middleware.JWTAuthMiddleware)
	protectedRouter.HandleFunc("/users/{id}", handlers.HandleGetUserByID).Methods("GET")
	protectedRouter.HandleFunc("/organisations", handlers.HandleGetUserOrganisations).Methods("GET")
	protectedRouter.HandleFunc("/organisations/{orgId}", handlers.HandleGetOrgByID).Methods("GET")
	protectedRouter.HandleFunc("/organisations", handlers.HandleCreateOrg).Methods("POST")
	protectedRouter.HandleFunc("/organisations/{orgId}/users", handlers.HandleAddUserToOrganisation).Methods("POST")

	log.Fatal(http.ListenAndServe(s.addr, router))
}
