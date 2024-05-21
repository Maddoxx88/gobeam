package main

import (
	"fmt"
	"io"
	"mime/multipart"
)

func GetExtension(f_name string) string {
	var temp string = ""
	for index, value := range f_name {
		t := string(value)
		i := index + 1
		if t == "." {
			for {
				if i == len(f_name) {
					break
				}
				if string(f_name[i]) == "." {
					temp = ""
					break
				}
				temp += string(f_name[i])
				i++
			}
		}
	}
	return temp
}

func ReadFileToBytes(file multipart.File) ([]byte, error) {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error Reading file ", err)
		return nil, err
	}

	return fileBytes, nil
}
