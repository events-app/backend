package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/auth0/go-jwt-middleware"
	"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
	"github.com/events-app/events-api/handlers"
	"github.com/events-app/events-api/internal/platform/web"
	"github.com/gorilla/mux"
)

const key = "KLHkjhsd*h67r3gJhjuds"
const maxUploadSize = 2 * 1048576 // bytes = 2 mb
const uploadPath = "./uploaded-files"
const serverPort = "8000"

func main() {
	r := mux.NewRouter()
	// use middleware handler
	r.Use(handlers.HeaderMiddleware)
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(key), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err string) {
			web.ErrorJSON(w, err, http.StatusInternalServerError)
		},
	})

	r.HandleFunc("/", handlers.Info).Methods("GET")
	r.HandleFunc("/api/v1/health", handlers.HealthCheck).Methods("GET")
	r.Handle("/api/v1/cards/secured", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(handlers.SecuredContent)),
	)).Methods("GET")
	r.HandleFunc("/api/v1/cards/{name}", handlers.GetCard).Methods("GET")
	r.HandleFunc("/api/v1/cards", handlers.GetCards).Methods("GET")
	r.HandleFunc("/api/v1/cards", handlers.AddCard).Methods("POST")
	r.HandleFunc("/api/v1/cards/{name}", handlers.UpdateCard).Methods("PUT")
	r.HandleFunc("/api/v1/cards/{name}", handlers.DeleteCard).Methods("DELETE")

	r.HandleFunc("/api/v1/login", Login).Methods("POST")

	r.HandleFunc("/api/v1/menus/{name}", handlers.GetMenu).Methods("GET")
	r.HandleFunc("/api/v1/menus", handlers.GetMenus).Methods("GET")
	r.HandleFunc("/api/v1/menus", handlers.AddMenu).Methods("POST")
	r.HandleFunc("/api/v1/menus/{name}", handlers.UpdateMenu).Methods("PUT")
	r.HandleFunc("/api/v1/menus/{name}", handlers.DeleteMenu).Methods("DELETE")
	r.HandleFunc("/api/v1/upload", handlers.UploadFile(uploadPath, maxUploadSize)).Methods("POST")
	
	// r.PathPrefix("/files/").Handler(http.FileServer(http.Dir(uploadPath)))
	fs := http.FileServer(http.Dir(uploadPath))
	// --- r.PathPrefix("/files/").Handler(http.StripPrefix("files/", fs))
	// r.Handle("/files", http.StripPrefix("/files", fs)).Methods("GET")
	r.HandleFunc("/files", handlers.GetFiles(uploadPath)).Methods("GET")
	r.Handle("/files/{file}", http.StripPrefix("/files", fs)).Methods("GET")
	// http.Handle("/files/", http.StripPrefix("/files", fs))

	// temporary handlers for backward compatibility with frontend
	r.Handle("/api/v1/content/secured", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(handlers.SecuredContent)),
	)).Methods("GET")
	r.HandleFunc("/api/v1/content/{name}", handlers.GetCard).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = serverPort
	}
	log.Println("Listening on port", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), r); err != nil {
		log.Printf("error: listing and serving: %s", err)
		return
	}
}
