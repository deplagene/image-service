package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"deplagene/image-service/configs"
	"deplagene/image-service/types"
)

func main() {
	resp, err := CheckImageForNsfwContent(types.ImgPath, configs.Envs.NsfwApiUrl)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	var nsfwResult types.NsfwResult

	if err := json.NewDecoder(resp.Body).Decode(&nsfwResult); err != nil {
		log.Fatal(err)
	}

	if nsfwResult.IsNsfw == true && nsfwResult.ConfidencePercentage >= types.MinConfidenceThreshold {
		fmt.Printf("Confidence: %.2f%%\n", nsfwResult.ConfidencePercentage)
	} else {
		fmt.Printf("%+v\n", nsfwResult)
	}
}

func CheckImageForNsfwContent(imgPath string, apiUrl string) (*http.Response, error) {
	const op = "services.image.CheckImageForNsfwContent"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	file, err := os.Open(imgPath)
	if err != nil {
		return nil, fmt.Errorf("%s:%v", op, err)
	}

	defer file.Close()

	part, err := writer.CreateFormFile("file", imgPath)
	if err != nil {
		return nil, fmt.Errorf("%s:%v", op, err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("%s:%v", op, err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("%s:%v", op, err)
	}

	req, err := http.NewRequest("POST", apiUrl, body)
	if err != nil {
		return nil, fmt.Errorf("%s:%v", op, err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)

	return resp, nil
}
