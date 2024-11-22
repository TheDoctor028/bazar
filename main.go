package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/TheDoctor028/bazar/app/account"
	adminapp "github.com/TheDoctor028/bazar/app/admin"
	"github.com/TheDoctor028/bazar/app/api"
	"github.com/TheDoctor028/bazar/app/home"
	"github.com/TheDoctor028/bazar/app/pages"
	"github.com/TheDoctor028/bazar/app/static"
	"github.com/TheDoctor028/bazar/internal/config"
	"github.com/TheDoctor028/bazar/internal/config/application"
	"github.com/TheDoctor028/bazar/internal/config/auth"
	"github.com/TheDoctor028/bazar/internal/config/bindatafs"
	"github.com/TheDoctor028/bazar/internal/db"
	"github.com/TheDoctor028/bazar/utils/funcmapmaker"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/qor/admin"
	"github.com/qor/publish2"
	"github.com/qor/qor"
	"github.com/qor/qor/utils"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	cmdLine := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	compileTemplate := cmdLine.Bool("compile-templates", false, "Compile Templates")
	cmdLine.Parse(os.Args[1:])

	var (
		Router = chi.NewRouter()
		Admin  = admin.New(&admin.AdminConfig{
			SiteName: "QOR DEMO",
			Auth:     auth.AdminAuth{},
			DB:       db.DB.Set(publish2.VisibleMode, publish2.ModeOff).Set(publish2.ScheduleMode, publish2.ModeOff),
		})
		Application = application.New(&application.Config{
			Router: Router,
			Admin:  Admin,
			DB:     db.DB,
		})
	)

	funcmapmaker.AddFuncMapMaker(auth.Auth.Config.Render)

	Router.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Add("Access-Control-Allow-Origin", "*")
			handler.ServeHTTP(w, req)
		})
	})

	Router.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			req.Header.Del("Authorization")
			handler.ServeHTTP(w, req)
		})
	})

	Router.Use(middleware.RealIP)
	Router.Use(middleware.Logger)
	Router.Use(middleware.Recoverer)
	Router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			var (
				tx         = db.DB
				qorContext = &qor.Context{Request: req, Writer: w}
			)

			if locale := utils.GetLocale(qorContext); locale != "" {
				tx = tx.Set("l10n:locale", locale)
			}

			ctx := context.WithValue(req.Context(), utils.ContextDBName, publish2.PreviewByDB(tx, qorContext))
			next.ServeHTTP(w, req.WithContext(ctx))
		})
	})

	Application.Use(api.New(&api.Config{}))
	Application.Use(adminapp.New(&adminapp.Config{}))
	Application.Use(home.New(&home.Config{}))
	Application.Use(account.New(&account.Config{}))
	Application.Use(pages.New(&pages.Config{}))
	Application.Use(static.New(&static.Config{
		Prefixs: []string{"/system"},
		Handler: utils.FileServer(http.Dir(filepath.Join(config.Root, "public"))),
	}))
	Application.Use(static.New(&static.Config{
		Prefixs: []string{"javascripts", "stylesheets", "images", "dist", "fonts", "vendors", "favicon.ico"},
		Handler: bindatafs.AssetFS.FileServer(http.Dir("public"), "javascripts", "stylesheets", "images", "dist", "fonts", "vendors", "favicon.ico"),
	}))

	if *compileTemplate {
		bindatafs.AssetFS.Compile()
	} else {
		fmt.Printf("Listening on: %v\n", config.Config.Port)
		if config.Config.HTTPS {
			if err := http.ListenAndServeTLS(fmt.Sprintf(":%d", config.Config.Port), "config/local_certs/server.crt", "config/local_certs/server.key", Application.NewServeMux()); err != nil {
				panic(err)
			}
		} else {
			if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), Application.NewServeMux()); err != nil {
				panic(err)
			}
		}
	}
}
