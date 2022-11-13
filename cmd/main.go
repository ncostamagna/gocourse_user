package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/ncostamagna/gocourse_user/internal/user"
	"github.com/ncostamagna/gocourse_user/pkg/bootstrap"
)

func main() {

	router := mux.NewRouter()
	_ = godotenv.Load()
	l := bootstrap.InitLogger()

	db, err := bootstrap.DBConnection()
	if err != nil {
		l.Fatal(err)
	}

	pagLimDef := os.Getenv("PAGINATOR_LIMIT_DEFAULT")
	if pagLimDef == "" {
		l.Fatal("paginator limit default is required")
	}

	userRepo := user.NewRepo(l, db)
	userSrv := user.NewService(l, userRepo)
	userEnd := user.MakeEndpoints(userSrv, user.Config{LimPageDef: pagLimDef})

	router.HandleFunc("/users", userEnd.Create).Methods("POST")
	router.HandleFunc("/users/{id}", userEnd.Get).Methods("GET")
	router.HandleFunc("/users", userEnd.GetAll).Methods("GET")
	router.HandleFunc("/users/{id}", userEnd.Update).Methods("PATCH")
	router.HandleFunc("/users/{id}", userEnd.Delete).Methods("DELETE")

	port := os.Getenv("PORT")
	address := fmt.Sprintf("127.0.0.1:%s", port)

	srv := &http.Server{
		Handler:      router,
		Addr:         address,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
