package router

import (
	"html/template"
	"io"
	"log"

	"tanya_dokter_app/app/controllers"
	"tanya_dokter_app/app/middlewares"

	_ "tanya_dokter_app/docs"

	"github.com/labstack/echo/v4"
)

func Init(app *echo.Echo) {
	// renderer := &TemplateRenderer{
	// 	templates: template.Must(template.ParseGlob("*.html")),
	// }
	// app.Renderer = renderer

	app.Use(middlewares.Cors())
	app.Use(middlewares.Secure())
	app.Use(middlewares.Gzip())
	app.Use(middlewares.Recover())

	app.GET("/", controllers.Index)
	app.GET("/test", controllers.Test)
	app.GET("/version", controllers.Version)
	// app.GET("/swagger/*", echoSwagger.WrapHandler)
	// app.GET("/docs", func(c echo.Context) error {
	// 	config := config.LoadConfig()

	// 	data := map[string]interface{}{
	// 		"BaseUrl": config.BaseUrl,
	// 		"Title":   "API Documentation of " + config.AppName,
	// 	}

	// 	if err := c.Render(http.StatusOK, "docs.html", data); err != nil {
	// 		fmt.Println("Render error:", err)
	// 		return err
	// 	}

	// 	return nil
	// })
	app.Static("/assets", "assets")

	api := app.Group("/v1", middlewares.StripHTMLMiddleware, middlewares.CheckAPIKey())
	{
		auth := api.Group("/auth")
		{
			auth.GET("/csrf", controllers.CSRFToken)
			auth.POST("/signup", controllers.SignUp)
			auth.POST("/signin", controllers.SignIn)
			// auth.POST("/forgot-password", controllers.ForgotPassword)

			// auth.POST("/signin/google/mobile", controllers.GoogleSignInMobile)
			// auth.POST("/reset-password", controllers.ResetPassword)
			auth.POST("/resend-email-verification", controllers.ResendEmailVerification)
			auth.POST("/email-verification", controllers.EmailVerification)
			// auth.GET("/user", controllers.GetSignInUser, middlewares.Auth())
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
