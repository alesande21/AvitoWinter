package app

import (
	"AvitoWinter/internal/colorAttribute"
	config2 "AvitoWinter/internal/config"
	http3 "AvitoWinter/internal/controllers/http"
	"AvitoWinter/internal/database"
	repository2 "AvitoWinter/internal/repository"
	service2 "AvitoWinter/internal/service"
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	log2 "github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func RunApp() error {

	// Настройка логера
	SetLevel("debug", "console")
	log2.Info("Настройка логера...")

	//Загрузка конфига
	log.Println("Загрузка конфига для базы данных...")
	config, err := config2.GetDefaultConfig()
	if err != nil {
		return fmt.Errorf("-> config2.GetDefaultConfig%w", err)
	}

	// Инициализация базы данных
	log.Println("Инициализация базы данных...")
	var conn *database.DBConnection
	conn, err = database.Open(config.GetDBsConfig())
	if err != nil {
		return fmt.Errorf("-> database2.Open%w", err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log2.Infof("RunConsumer-> conn.Close:%s", err)
		}
	}()

	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()

	// Инициализация репозитория
	log2.Info("Инициализация репозитория...")
	var postgresRep database.DBRepository
	postgresRep, err = database.CreatePostgresRepository(conn.GetConn)
	if err != nil {
		return fmt.Errorf("-> database2.CreatePostgresRepository%w", err)
	}

	// Инициализация кеша
	//log2.Info("Инициализация кеш клиента...")
	//redisClient, err := cache2.NewRedisClientAdapter(ctx)
	//if err != nil {
	//	return fmt.Errorf("->  cache2.NewRedisClientAdapterr%v", err)
	//}
	//cacheClient := cache2.NewCacheClient(redisClient)

	// Инициализация сервиса
	log2.Info("Инициализация сервиса...")
	shopRepo := repository2.NewShopRepo(postgresRep)
	//repoWithCache := repository2.NewUserRepoWithCache(shopRepo, redisClient)
	shopService := service2.NewShopService(shopRepo)

	log2.Info("Загрузка настроек для сервера...")
	var serverAddress http3.ServerAddress
	err = serverAddress.UpdateEnvAddress()
	if err != nil {
		return fmt.Errorf("-> serverAddress.UpdateEnvAddress%w", err)
	}

	log2.Info("Инициализация и старт сервера...")
	swagger, err := http3.GetSwagger()
	if err != nil {
		return fmt.Errorf("->  http2.GetSwagger%w", err)
	}
	swagger.Servers = nil

	userServer := http3.NewUserServer(shopService)
	r := mux.NewRouter()
	//r.Use(middleware.OapiRequestValidator(swagger))

	handler := http3.HandlerWithOptions(userServer, http3.GorillaServerOptions{
		BaseRouter: r,
		BaseURL:    "",
	})

	s := &http.Server{
		Addr:    serverAddress.EnvAddress,
		Handler: handler,
	}

	//s := &http.Server{
	//	Addr:    serverAddress.EnvAddress,
	//	Handler: r,
	//}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	defer close(interrupt)

	shutDownChan := make(chan error, 1)
	defer close(shutDownChan)

	go func() {
		shutDownChan <- s.ListenAndServe()
	}()

	log2.Infof("Подключнеие установлено -> %s", colorAttribute.ColorString(colorAttribute.FgYellow, serverAddress.EnvAddress))

	select {
	case sig := <-interrupt:
		log2.Infof("Приложение прерывается: %s", sig)
		ctxShutDown, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)

		//cancel()

		defer cancelShutdown()
		err := s.Shutdown(ctxShutDown)
		if err != nil {
			return fmt.Errorf("-> s.Shutdown: %w", err)
		}

		log2.Info("Сервер завершил работу")
	case err := <-shutDownChan:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf(": ошибка при запуске сервера: %w", err)
		}
	}

	return nil
}
