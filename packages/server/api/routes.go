package routes

import (
	"area-server/api/middlewares"
	appletr "area-server/api/routes/applet"
	appletcontextr "area-server/api/routes/applet/context"
	appletnewr "area-server/api/routes/applet/new"
	authr "area-server/api/routes/auth"
	authenticatorsr "area-server/api/routes/authenticators"
	authorizationr "area-server/api/routes/authorization"
	servicesr "area-server/api/routes/services"
	storer "area-server/api/routes/store"
	userr "area-server/api/routes/user"
	"area-server/authenticators"
	"area-server/config"
	"time"

	services "area-server/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// It initializes the routes of the server
func Init(app *fiber.App) {

	// Routes

	// TODO: Admin routes (Admin Dashboard)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).SendString("Area Server - v0.1.0")
	})

	// About.json
	app.Get("/about.json", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"client": fiber.Map{
				"host": c.Hostname(),
			},
			"server": fiber.Map{
				"current_time":   time.Now(),
				"authenticators": authenticators.List,
				"services":       services.List,
			},
		})
	})

	// Change the middleware to use the token middleware if problems with sessions
	cmiddleware := middlewares.SessionMiddleware
	if config.CFG.Mode == config.Token {
		cmiddleware = middlewares.TokenMiddleware
	}

	// Static
	app.Static("/assets", "./assets")
	app.Static("/avatars", "./avatars")

	// Auth
	auth := app.Group("/auth")
	auth.Post("/login", authr.Login)
	auth.Post("/register", authr.Register)
	auth.Post("/external", authr.ExternalAuth)

	// Connected User
	user := app.Group("/me", cmiddleware)
	user.Get("/", userr.GetUser)
	user.Post("/logout", userr.LogoutUser)
	user.Put("/", userr.UpdateUser)
	user.Delete("/", userr.DeleteUser)

	// Avatar
	avatar := user.Group("/avatar")
	avatar.Get("/", userr.GetAvatar)
	avatar.Put("/", userr.UpdateAvatar)

	store := app.Group("/store") // List all public applets
	store.Get("/", storer.GetStoreApplets)

	// Applet (Area)
	applet := app.Group("/applet", cmiddleware)

	appletnew := applet.Group("/new")
	appletnew.Get("/", appletnewr.GetNewApplet)                                             // Default public=false
	appletnew.Put("/", middlewares.NewAppletMiddleware, appletnewr.AddStateToNewApplet)     // Default public=false
	appletnew.Post("/", middlewares.NewAppletMiddleware, appletnewr.SubmitNewApplet)        // Default public=false
	appletnew.Delete("/", middlewares.NewAppletMiddleware, appletnewr.DeleteStateNewApplet) // Default public=false

	applet.Get("/", appletr.GetApplets) // Get all applets of the user + (BONUS) Get all public applets, algorithm to evalutate score of applets

	appletcurrent := applet.Group("/:applet_id")
	appletcurrent.Get("/", appletr.GetApplet) // Get an applet by id
	//appletcurrent.Put("/", appletr.UpdateApplet)           // Update an applet by id
	appletcurrent.Patch("/", appletr.UpdateAppletActivity) // Update an applet activity by id
	appletcurrent.Delete("/", appletr.DeleteApplet)        // Delete an applet by id
	appletcurrent.Get("/reactions", appletr.GetAppletReactions)

	appletcurrent.Put("/start", appletr.StartApplet) // Start an applet by id
	appletcurrent.Put("/stop", appletr.StopApplet)   // Stop an applet by id

	appletlogs := app.Group("/logs/:applet_id")
	appletlogs.Get("/", websocket.New(appletcontextr.GetAppletLogs))

	// Authorization
	authorization := app.Group("/authorization", cmiddleware)
	authorization.Get("/", authorizationr.GetAuthorizations)
	authorization.Post("/", authorizationr.CreateAuthorization)
	authorization.Delete("/:name", authorizationr.DeleteAuthorization)

	authorization.Get("/services", authorizationr.GetAuthorizedServices)

	// Authenticators routes
	authenticatorsL := app.Group("/authenticators")
	authenticatorsL.Get("/", authenticatorsr.GetAuthenticators)
	authenticatorsL.Get("/:authenticator", authenticatorsr.GetAuthenticator)

	// Service routes
	servicesL := app.Group("/services")
	servicesL.Get("/", servicesr.GetServices)

	serviceRoutes := servicesL.Group("/:service")
	serviceRoutes.Get("/", servicesr.GetService)
	serviceRoutes.Get("/actions", servicesr.GetServiceActions)
	serviceRoutes.Get("/actions/:action", servicesr.GetServiceAction)
	serviceRoutes.Get("/reactions", servicesr.GetServiceReactions)
	serviceRoutes.Get("/reactions/:reaction", servicesr.GetServiceReaction)
	serviceRoutes.Get("/api", servicesr.GetApiEndpoints)

	for _, service := range services.List {
		serviceR2 := servicesL.Group("/" + service.Name + "/api")
		for _, route := range service.Routes {
			if route.NeedAuth == true {
				serviceR2.Add(route.Method, route.Endpoint, cmiddleware, route.Handler)
			} else {
				serviceR2.Add(route.Method, route.Endpoint, route.Handler)
			}
		}
	}
}
