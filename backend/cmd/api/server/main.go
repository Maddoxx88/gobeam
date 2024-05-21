package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"github.com/Maddoxx88/gobeam/database"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type config struct {
	version string
	port    string
}

type application struct {
	config config
	db     *sql.DB
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Issues reading from the env file ,", err.Error())
	}

	//getting the required values from env files
	address := os.Getenv("APP_Addr")
	port := os.Getenv("APP_Port")

	//opening a sqllite connection
	db, err := sql.Open("sqlite3", os.Getenv("DB"))

	if err != nil {
		fmt.Println("Error establishing a connection")
		log.Fatal(err)
	}

	defer db.Close()

	//Making sure the database connection is available
	err = db.Ping()
	if err != nil {
		log.Fatal("Cannot connect to database ", err)
	}

	//Initializing the database and creating the table and trigger
	database.InitDB(db)

	router := mux.NewRouter() //initializing the router

	//initializing the application structure
	app := &application{
		config: config{version: "1.0.0", port: port},
		db:     db,
	}

	//handling the endpoints
	router.HandleFunc("/upload", app.UploadFileHandler).Methods("POST")
	router.HandleFunc("/download/{file-name}", app.RetrieveFileHandler).Methods("GET")
	router.HandleFunc("/list", app.GetAllFilesHandler).Methods("GET")
	router.HandleFunc("/download/{file-name}", app.DeleteFileHandler).Methods("DELETE")

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
