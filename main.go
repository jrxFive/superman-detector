package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/jrxfive/superman-detector/handlers/healthz"
	v1 "github.com/jrxfive/superman-detector/handlers/v1"
	"github.com/jrxfive/superman-detector/internal/pkg/settings"
	"github.com/jrxfive/superman-detector/internal/pkg/signals"
	customMiddleware "github.com/jrxfive/superman-detector/middleware"
	"github.com/jrxfive/superman-detector/models"
	"github.com/jrxfive/superman-detector/pkg/geoip"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func signalMonitor(signalChannel chan os.Signal, e *echo.Echo) {
	s := <-signalChannel
	e.Logger.Warnf("signal interrupt detected:%s attempting graceful shutdown", s.String())
	err := e.Shutdown(context.Background())
	if err != nil {
		e.Logger.Warnf("failed to gracefully shutdown", s.String())
	}
	os.Exit(0)
}

func main() {
	//Configuration
	s := settings.NewSettings()

	//Echo
	e := echo.New()

	//Statsd Telemetry
	statsdClient, err := statsd.New(s.StatsdAddress, func(options *statsd.Options) error {
		options.Namespace = s.StatsdNamespace
		options.Tags = []string{fmt.Sprintf("app:%s", s.StatsdNamespace)}
		options.BufferPoolSize = s.StatsdBufferPoolSize
		options.Telemetry = false
		return nil
	})
	defer func() {
		err = statsdClient.Close()
		e.Logger.Error(err)
	}()

	//Middleware
	e.Use(middleware.Logger())

	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return uuid.New().String()
		},
	}))

	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "",
		ContentTypeNosniff:    "",
		XFrameOptions:         "",
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: "default-src 'self'",
	}))

	e.Use(middleware.BodyLimit(s.RequestBodyLimit))

	e.Use(customMiddleware.NewStats(statsdClient).Process)

	//Signal Monitor
	smc := signals.NewSignalMonitoringChannel()
	go signalMonitor(smc, e)

	//Databases
	geoDB, err := geoip.NewDefaultLocator(s)
	if err != nil {
		e.Logger.Fatalf("failed to open geo database err:%s", err.Error())
	}
	defer func() {
		err = geoDB.Close()
		e.Logger.Error(err)
	}()

	db, err := gorm.Open(s.SqlDialect, s.SqlConnectionString)
	if err != nil {
		e.Logger.Fatalf("failed to open database err:%s", err.Error())
	}
	defer func() {
		err = db.Close()
		e.Logger.Error(err)
	}()

	//Database create if missing, no-op if created
	if !db.HasTable(&models.LoginEvent{}) {
		db.CreateTable(&models.LoginEvent{})
	}

	//Handler Creation
	health := healthz.NewHealthz(db.DB(), statsdClient)

	//v1 Handler Creation
	login := v1.NewLogin(db, geoDB, statsdClient, s)

	//Handlers Registration
	e.GET("/healthz", health.GetHealthz)
	e.HEAD("/healthz", health.HeadHealthz)

	v1Group := e.Group("/v1")
	v1Group.POST("", login.PostLogin)

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", s.ServicePort),
		ReadTimeout:       s.ServerReadTimeoutSeconds,
		ReadHeaderTimeout: s.ServerReadTimeoutSeconds,
		WriteTimeout:      s.ServerWriteTimeoutSeconds,
	}

	e.Logger.Fatal(e.StartServer(server))
}
