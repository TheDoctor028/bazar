package api

import (
	"github.com/TheDoctor028/bazar/internal/config/application"
	"github.com/TheDoctor028/bazar/internal/config/db"
	"github.com/TheDoctor028/bazar/models/users"
	"github.com/qor/admin"
	"github.com/qor/qor"
)

// New new home app
func New(config *Config) *App {
	if config.Prefix == "" {
		config.Prefix = "/api"
	}
	return &App{Config: config}
}

// App home app
type App struct {
	Config *Config
}

// Config home config struct
type Config struct {
	Prefix string
}

// ConfigureApplication configure application
func (app App) ConfigureApplication(application *application.Application) {
	API := admin.New(&qor.Config{DB: db.DB})

	API.AddResource(&users.User{})
	// User := API.AddResource(&users.User{})
	// userOrders, _ := User.AddSubResource("Orders")
	// userOrders.AddSubResource("OrderItems", &admin.Config{Name: "Items"})

	application.Router.Mount(app.Config.Prefix, API.NewServeMux(app.Config.Prefix))
}
