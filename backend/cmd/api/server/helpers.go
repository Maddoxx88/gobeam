package main

import (
	"fmt"
	"io"
	"mime/multipart"
	"strings"
)

func GetNameAndExtension(f_name string) (string, string) {
	arr := strings.Split(f_name, ".")
	name := ""
	for i, v := range arr {
		if i != (len(arr) - 1) {
			name += v
		}
	}
	return name, arr[len(arr)-1]
}

func ReadFileToBytes(file multipart.File) ([]byte, error) {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error Reading file ", err)
		return nil, err
	}

	return fileBytes, nil
}
