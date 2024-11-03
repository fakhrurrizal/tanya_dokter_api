package handler

import (
	"net/http"
	"tanya_dokter_app/app/router"
	"tanya_dokter_app/config"

	"github.com/labstack/echo/v4"
	_ "github.com/joho/godotenv/autoload"
)

// func main() {
// 	app := echo.New()
// 	router.Init(app)
// 	config.Database()

// 	app.Server.Addr = "127.0.0.1:" + config.LoadConfig().Port
// 	log.Printf("Server: " + config.LoadConfig().BaseUrl)
// 	log.Printf("Documentation: " + config.LoadConfig().BaseUrl + "/docs")
// 	graceful.ListenAndServe(app.Server, 5*time.Second)

// }

func Main(w http.ResponseWriter, r *http.Request) {
	e := Start()

	e.ServeHTTP(w, r)
}

func Start() *echo.Echo {
	app := echo.New()

	config.Database()

	router.Init(app)

	return app
}
