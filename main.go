package main

import (
	"log"
	"tanya_dokter_app/app/router"
	"tanya_dokter_app/config"
	"time"

	"github.com/labstack/echo/v4"
	"gopkg.in/tylerb/graceful.v1"
)

func main() {
	app := echo.New()
	router.Init(app)
	config.Database()

	app.Server.Addr = "127.0.0.1:" + config.LoadConfig().Port
	log.Printf("Server: " + config.LoadConfig().BaseUrl)
	log.Printf("Documentation: " + config.LoadConfig().BaseUrl + "/docs")
	graceful.ListenAndServe(app.Server, 5*time.Second)

}
