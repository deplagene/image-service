package api

import (
	"database/sql"
	"log/slog"
	"net/http"
	"teach-stack/cmd/metrics"
	"teach-stack/configs"
	"teach-stack/services/image"
	"teach-stack/services/rabbitmq"
	"teach-stack/services/s3"
	"teach-stack/utils"
)

type ApiStruct struct {
	addr string
	db   *sql.DB
}

func NewApi(addr string, db *sql.DB) *ApiStruct {
	return &ApiStruct{
		addr: addr,
		db:   db,
	}
}

func (a *ApiStruct) Run() error {
	router := http.NewServeMux()

	imageS3storage := s3.NewMinioProvider(configs.Envs.S3Url, configs.Envs.S3User, configs.Envs.S3Password, false)
	imageS3storage.Connect()
	imageStore := image.NewStore(a.db)
	imageService := image.NewService(imageStore, imageS3storage)
	conn, err := rabbitmq.Connect(configs.Envs.BrokerUser, configs.Envs.BrokerPassword, configs.Envs.BrokerHost)
	if err != nil {
		slog.Error("cannot connect to rabbitmq", utils.Err(err))
		return err
	}
	defer conn.Close()
	rb, err := rabbitmq.NewRabbitMqClient(conn)
	if err != nil {
		slog.Error("failed to create a new rabbitmq server", utils.Err(err))
		return err
	}
	defer rb.Close()

	imageHandlers := image.NewHandler(imageStore, imageS3storage, rb, imageService)
	imageHandlers.RegisterRoutes(router)

	handler := corsMiddleware(router)

	server := http.Server{
		Addr:    a.addr,
		Handler: handler,
	}

	go func() {
		imgConsumer := image.NewConsumer(imageService, rb)
		err = imgConsumer.Listen()
		if err != nil {
			slog.Error("failed to listen a consumer", utils.Err(err))
		}
	}()

	go func() {
		_ = metrics.Listen("localhost:8080")
	}()

	slog.Info("Server run successful running", slog.String("address", a.addr))
	return server.ListenAndServe()
}