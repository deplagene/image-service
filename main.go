package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"deplagene/image-service/configs"
	"deplagene/image-service/db"
	"deplagene/image-service/services/images"
	"deplagene/image-service/types"

	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
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

	// MongoDb
	// todo add function's
	client, err := db.NewMongoClient(configs.Envs.MongoDbConnUrl)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.TODO()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	store := images.NewStore(client)

	err = store.Create(ctx, "some-url.com")
	if err != nil {
		log.Fatal(err)
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
