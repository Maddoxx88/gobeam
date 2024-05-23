package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	db_model "github.com/Maddoxx88/gobeam/database"
	"github.com/gorilla/mux"
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

	//Getting the filename and extension
	fileName, extension := GetNameAndExtension(header.Filename)

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
	vars := mux.Vars(r)
	fileName := vars["file-name"]
	response := db_model.File{}
	query := `SELECT id,name,extension,data FROM FILES WHERE name=?`
	err := app.db.QueryRow(query, fileName).Scan(&response.Id, &response.Name, &response.Extension, &response.Data)
	if err != nil {
		log.Fatal("Error in executing the query ", err)
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", http.DetectContentType(response.Data))
	w.Header().Set("Content-Length", strconv.Itoa(len(response.Data)))

	_, err = w.Write(response.Data)
	if err != nil {
		log.Fatal("Error writing file to resposne ", err)
	}
}

func (app *application) GetAllFilesHandler(w http.ResponseWriter, r *http.Request) {
	query := `SELECT name,extension from FILES`
	rows, err := app.db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	response_files := []db_model.File{}

	for rows.Next() {
		var name, extension string
		if err := rows.Scan(&name, &extension); err != nil {
			log.Fatal(err)
		}
		temp := db_model.File{}
		temp.Name = name
		temp.Extension = extension
		response_files = append(response_files, temp)
	}

	response, err := json.Marshal(response_files)
	if err != nil {
		log.Fatal("Error in marshaling(converting resposne to JSON) the response ", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (app *application) DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	query := `DELETE FROM FILES WHERE Id=?`
	_, err := app.db.Exec(query, id)
	if err != nil {
		log.Fatal("Error while deleting the file from database with id ", id, " ", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprint("Successfully deleted the data with id ", id)))
}
