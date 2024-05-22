package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2023_1_Technokaif/cmd/api/init/server"
	codeDelivery "github.com/yarikTri/network-channel-layer/cmd/delivery/http"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/yarikTri/network-channel-layer/docs"

	swaggerFiles "github.com/swaggo/files" // swagger embed files
	swagger "github.com/swaggo/gin-swagger"

	flog "github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

// @title		КР СТ АСОИУ Сервис Канального Уровня
// @version		0.1.0
// @description	Сервис имитирует передачу данных через ненадёжную сеть с защитой данных с помощью кодировки Хэмминга [15, 11].

// @contact.name   Yaroslav Kuzmin
// @contact.email  yarik1448kuzmin@gmail.com

// @host localhost:8082
// @schemes https http
// @BasePath /

func main() {
	listenEndpoint := "localhost:8082"

	reqIdGetterMock := func(context.Context) (uint32, error) { return 0, nil }
	flogger, err := flog.NewFLogger(reqIdGetterMock)
	if err != nil {
		log.Fatalf("logger can not be defined: %v\n", err)
	}

	router := gin.Default()
	router.POST("/code", codeDelivery.Code)
	router.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler))

	var srv server.Server
	if err := srv.Init(listenEndpoint, router); err != nil {
		flogger.Errorf("error while launching server: %v", err)
	}

	go func() {
		if err := srv.Run(); err != nil {
			flogger.Errorf("server error: %v", err)
			os.Exit(1)
		}
	}()
	flogger.Info("trying to launch server")

	timer := time.AfterFunc(1*time.Second, func() {
		flogger.Infof("server launched at %s", listenEndpoint)
	})
	defer timer.Stop()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
}
