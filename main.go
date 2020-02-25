package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/jrxfive/superman-detector/handlers/healthz"
	v1 "github.com/jrxfive/superman-detector/handlers/v1"
	"github.com/jrxfive/superman-detector/internal/pkg/settings"
	"github.com/jrxfive/superman-detector/internal/pkg/signals"
	"github.com/jrxfive/superman-detector/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/oschwald/geoip2-golang"
)

func signalMonitor(signalChannel chan os.Signal, e *echo.Echo) {
	s := <-signalChannel
	e.Logger.Warnf("signal interrupt detected:%s attempting graceful shutdown", s.String())
	err := e.Shutdown(context.Background())
	if err != nil {
		e.Logger.Warnf("failed to gracefully shutdown", s.String())
	}
}

func main() {
	//Configuration
	s := settings.NewSettings()

	//Echo
	e := echo.New()

	//Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return uuid.New().String()
		},
	}))

	//Signal Monitor
	smc := signals.NewSignalMonitoringChannel()
	go signalMonitor(smc, e)

	//Databases
	geoDB, err := geoip2.Open(s.GeoIPDatabaseFileLocation)
	if err != nil {
		e.Logger.Fatalf("failed to open geo database err:%s", err.Error())
	}
	defer func() {
		err = geoDB.Close()
		e.Logger.Error(err)
	}()

	db, err := gorm.Open("sqlite3", "/tmp/superman.db")
	if err != nil {
		e.Logger.Fatalf("failed to open database err:%s", err.Error())
	}
	defer func() {
		err = db.Close()
		e.Logger.Error(err)
	}()

	//Database create if missing, no-op if created
	db.CreateTable(&models.LoginEvent{})

	//Handler Creation
	health := healthz.NewHealthz(db.DB())

	//v1 Handler Creation
	login := v1.NewLogin(db, geoDB, s)

	//Handlers Registation
	e.GET("/healthz", health.GetHealthz)
	e.HEAD("/healthz", health.HeadHealthz)

	loginGroup := e.Group("/v1")
	loginGroup.POST("", login.PostLogin)

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", s.ServicePort),
		ReadTimeout:       s.ServerReadTimeoutSeconds,
		ReadHeaderTimeout: s.ServerReadTimeoutSeconds,
		WriteTimeout:      s.ServerWriteTimeoutSeconds,
	}

	e.Logger.Fatal(e.StartServer(server))
}
