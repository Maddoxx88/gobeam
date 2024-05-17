package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	handlers "github.com/Maddoxx88/gobeam/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type config struct {
	version string
	port    string
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Issues reading from the env file ,", err.Error())
	}

	//opening a sqllite connection
	db, err := sql.Open("sqlite3", os.Getenv("DB"))

	if err != nil {
		fmt.Println("Error establishing a connection")
		log.Fatal(err)
	}

	//Make sure the server is running till it gets instructuction to kill
	// killcheck := make(chan os.Signal, 1)
	// signal.Notify(killcheck, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	router := mux.NewRouter() //initializing the router

	//handling the endpoints
	router.HandleFunc("/upload", handlers.UploadFileHandler).Methods("POST")
	router.HandleFunc("/download/{file-name}", handlers.RetrieveFileHandler).Methods("GET")
	router.HandleFunc("/list", handlers.GetAllFilesHandler).Methods("GET")
	router.HandleFunc("/download/{file-name}", handlers.DeleteFileHandler).Methods("DELETE")

	//getting the required values from env files
	address := os.Getenv("APP_Addr")
	port := os.Getenv("APP_Port")

	//starting the server

	fmt.Println("Starting the server on the port ", port)
	srv := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf("%s%s", address, port),
	}

	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal("Issue with starting the server at port ", port, " ", err.Error())
	}
}
