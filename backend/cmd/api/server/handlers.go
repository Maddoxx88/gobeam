package main

import (
	"encoding/json"
	"log"
	"net/http"

	db_model "github.com/Maddoxx88/gobeam/database"
)

func (app *application) UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	//Limit of file size is 2GB
	err := r.ParseMultipartForm(2 << 30)
	if err != nil {
		http.Error(w, "Unable to parse form request ", http.StatusBadRequest)
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to return form-data", http.StatusBadRequest)
		return
	}
	defer file.Close()

	//Getting the filename
	fileName := header.Filename
	//Getting the extension
	extension := GetExtension(fileName)

	return_struct := db_model.File{}

	byte_file, err := ReadFileToBytes(file)

	if err != nil {
		log.Fatal("Error converting file to bytes for inserting into database ", err)
	}

	insert_query := `INSERT INTO files (name,extension,data) VALUES (?,?,?) RETURNING id,name,extension`

	err = app.db.QueryRow(insert_query, fileName, extension, byte_file).Scan(&return_struct.Id, &return_struct.Name, &return_struct.Extension)

	if err != nil {
		log.Fatal("Error inserting file to database as ", err, " ", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response, err := json.Marshal(return_struct)
	if err != nil {
		log.Fatal("Error marshaling the json while uploading the file", err)
	}
	w.Write(response)
}

func (app *application) RetrieveFileHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) GetAllFilesHandler(w http.ResponseWriter, r *http.Request) {}

func (app *application) DeleteFileHandler(w http.ResponseWriter, r *http.Request) {}
