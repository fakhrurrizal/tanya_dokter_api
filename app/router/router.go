package router

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	"tanya_dokter_app/app/controllers"
	"tanya_dokter_app/app/middlewares"
	"tanya_dokter_app/config"

	_ "tanya_dokter_app/docs"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func Init(app *echo.Echo) {
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("*.html")),
	}
	app.Renderer = renderer
	app.Use(middlewares.Cors())
	app.Use(middlewares.Secure())
	app.Use(middlewares.Gzip())
	app.Use(middlewares.Recover())
	app.Use(middlewares.Logger())

	app.GET("/", controllers.Index)
	app.GET("/test", controllers.Test)
	app.GET("/version", controllers.Version)
	app.GET("/swagger/*", echoSwagger.WrapHandler)
	app.GET("/docs", func(c echo.Context) error {
		err := c.Render(http.StatusOK, "docs.html", map[string]interface{}{
			"BaseUrl": config.LoadConfig().BaseUrl,
			"Title":   "Api Documentation of " + config.LoadConfig().AppName,
		})
		fmt.Println("err:", err)
		return err
	})
	app.Static("/assets", "assets")

	api := app.Group("/v1", middlewares.StripHTMLMiddleware, middlewares.CheckAPIKey())
	{
		auth := api.Group("/auth")
		{
			auth.POST("/signin", controllers.SignIn)
			auth.POST("/signup", controllers.SignUp)
			auth.POST("/forgot-password", controllers.ForgotPassword)

			// auth.POST("/signin/google/mobile", controllers.GoogleSignInMobile)
			auth.POST("/reset-password", controllers.ResetPassword)
			auth.POST("/resend-email-verification", controllers.ResendEmailVerification)
			auth.POST("/email-verification", controllers.EmailVerification)
			auth.GET("/user", controllers.GetSignInUser, middlewares.Auth())
			// auth.PUT("/user-profile", controllers.UpdateUserProfileByID, middlewares.Auth())
		}
		role := api.Group("/role")
		{
			role.POST("", controllers.CreateRole, middlewares.Auth())
			role.GET("", controllers.GetRoles)
			role.GET("/:id", controllers.GetRoleByID)
			role.DELETE("/:id", controllers.DeleteRoleByID, middlewares.Auth())
			role.PUT("/:id", controllers.UpdateRoleByID, middlewares.Auth())
		}
		user := api.Group("/user")
		{
			user.POST("", controllers.CreateUser, middlewares.Auth())
			user.GET("", controllers.GetUsers)
			user.GET("/:id", controllers.GetUserByID)
			user.DELETE("/:id", controllers.DeleteUserByID, middlewares.Auth())
			user.PUT("/:id", controllers.UpdateUserByID, middlewares.Auth())
		}
		category_specialist := api.Group("/category-specialist")
		{
			category_specialist.POST("", controllers.CreateCategorySpecialist, middlewares.Auth())
			category_specialist.GET("", controllers.GetCategorySpecialists)
			category_specialist.GET("/:id", controllers.GetCategorySpecialistByID)
			category_specialist.DELETE("/:id", controllers.DeleteCategorySpecialistByID, middlewares.Auth())
			category_specialist.PUT("/:id", controllers.UpdateCategorySpecialistByID, middlewares.Auth())
		}
		data_drugs := api.Group("/data-drugs")
		{
			data_drugs.POST("", controllers.CreateDataDrugs, middlewares.Auth())
			data_drugs.GET("", controllers.GetDataDrugs)
			data_drugs.GET("/:id", controllers.GetDataDrugsByID)
			data_drugs.DELETE("/:id", controllers.DeleteDataDrugsByID, middlewares.Auth())
			data_drugs.PUT("/:id", controllers.UpdateDataDrugsByID, middlewares.Auth())
		}
		files := api.Group("/file", middlewares.Auth())
		{
			files.POST("", controllers.UploadFile)
			files.GET("", controllers.GetFile)
		}
		chat := api.Group("/chat")
		{
			chat.GET("/messages/:user_pengirim/:user_penerima", controllers.GetMessagesByUsersHandler)
			chat.GET("/ws/:user_pengirim/:user_penerima", controllers.HandleWebSocket) // WebSocket untuk komunikasi real-time
			chat.POST("/send", controllers.SendMessageHandler)                         // Mengirim pesan baru
		}

	}
	log.Printf("Server started...")
}

type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}
