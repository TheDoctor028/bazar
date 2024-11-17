package main

import (
	"fmt"
	"github.com/TheDoctor028/bazar/internal/utils"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/qor/admin"
	"net/http"
)

// Define a GORM-backend model
type User struct {
	gorm.Model
	Name string
}

// Define another GORM-backend model
type Product struct {
	gorm.Model
	Name        string
	Description string
}

func main() {
	// Set up the database
	DB, _ := gorm.Open("sqlite3", "demo.db")
	DB.AutoMigrate(&User{}, &Product{})

	// Initalize
	Admin := admin.New(&admin.AdminConfig{
		SiteName: utils.GetEnvOrDefault("SITE_NAME", "QOR Admin"),
		DB:       DB,
	})

	// Create resources from GORM-backend model
	Admin.AddResource(&User{})
	Admin.AddResource(&Product{})

	// Initalize an HTTP request multiplexer
	mux := http.NewServeMux()

	// Mount admin to the mux
	Admin.MountTo("/admin", mux)

	PORT := utils.GetEnvOrDefault("PORT", "8888")
	fmt.Println("Listening on: ", PORT)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", PORT), mux); err != nil {
		panic(err)
	}
}
