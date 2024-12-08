package main

import (
	"log"
	router "tanya_dokter_app/app/router"
	"tanya_dokter_app/config"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"gopkg.in/tylerb/graceful.v1"
)

// @title Tanya Dokter API
// @version V1.2412.081710
// @description API documentation by Nizom

// @securityDefinitions.apikey JwtToken
// @in header
// @name Authorization

func main() {
	app := echo.New()
	router.Init(app)
	config.Database()

	app.Server.Addr = "127.0.0.1:" + config.LoadConfig().Port
	log.Printf("Server: " + config.LoadConfig().BaseUrl)
	log.Printf("Documentation: " + config.LoadConfig().BaseUrl + "/docs")
	graceful.ListenAndServe(app.Server, 5*time.Second)
}

// func Main(w http.ResponseWriter, r *http.Request) {
// 	e := Start()

// 	e.ServeHTTP(w, r)
// }

// func Start() *echo.Echo {
// 	app := echo.New()

// 	config.Database()

// 	router.Init(app)

// 	return app
// }
