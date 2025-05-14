package router

import (
	"net/http"
	"auth/internal/handler"
)

func SetupRoutes() {
	http.HandleFunc("POST /register", handler.Register)
	http.HandleFunc("POST /login", handler.Login)

	// api.Use(jwtware.New(jwtware.Config{
	// 	SigningKey:  jwtware.SigningKey{Key: []byte(config.Env.SecretKey)},
	// 	TokenLookup: "header:Authorization,cookie:Authorization",
	// }))

	// api.Get("/projects", handler.GetProjects)
	// api.Get("/projects/:id", handler.GetProject)
	//
	// api.Get("/tasks", handler.GetTasks)
	// api.Get("/tasks/:id", handler.GetTask)
	//
	// api.Get("/solutions", handler.GetSolutions)
	// api.Get("/solutions/:id", handler.GetSolution)
	//
	// api.Post("/projects", handler.LoadProject)
	// api.Put("/projects/:id", handler.LoadProject)
	// api.Delete("/projects/:id", handler.DeleteProject)
	// api.Post("/solutions", handler.SubmitSolution)
	//
	// api.Get("/delayed-tasks", handler.GetDelayedTasks)
	// api.Get("/delayed-tasks/:id", handler.GetDelayedTask)
}