package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/andrewwebber/tinyurl/pkg/infra/db/couchbase"
	usescases "github.com/andrewwebber/tinyurl/pkg/usecases"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	var (
		baseURL         = flag.String("baseurl", "http://localhost:8080", "TinyURL base url for all shortened urls")
		clusterURL      = flag.String("cluster-url", "couchbase://localhost", "Couchbase cluster connection string")
		clusterUsername = flag.String("cluster-username", "Administrator", "Couchbase cluster authentication username")
		clusterPassword = flag.String("cluster-password", "password", "Couchbase cluster authentication password")
		bucketName      = flag.String("bucket-name", "tinyurl", "Couchbase bucket to store tiny urls")
	)

	flag.Parse()
	db, err := couchbase.New(*clusterURL, *clusterUsername, *clusterPassword, *bucketName)
	if err != nil {
		log.Fatal(err)
	}

	repository := couchbase.NewEntitiesRepository(db)
	tinyUrl := usescases.NewTinyURL(*baseURL, repository, usescases.XIDURLShortener)

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	log.Println(e.POST("/client/v1/tinyurl/shorten", shorten(tinyUrl)))
	e.GET("/:shorturl", serve(tinyUrl))

	// Start server
	go func() {
		if err := e.Start(":8080"); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

// Handlers
type RequestShortURL struct {
	URL string `json:"url" form:"url" query:"url"`
}

type ShortUrlResponse struct {
	Short string `json:"short" form:"short" query:"short"`
	URL   string `json:"url" form:"url" query:"url"`
}

func shorten(tinyUrl usescases.TinyURL) echo.HandlerFunc {
	return func(c echo.Context) error {
		r := new(RequestShortURL)
		if err := c.Bind(r); err != nil {
			log.Println(err)
			return err
		}

		s, err := tinyUrl.ShortenURL(r.URL)
		if err != nil {
			log.Println(err)
			return err
		}

		return c.JSON(http.StatusOK, &ShortUrlResponse{Short: s.Short, URL: r.URL})
	}
}

func serve(tinyURL usescases.TinyURL) echo.HandlerFunc {
	return func(c echo.Context) error {
		s := c.Param("shorturl")
		if s == "" {
			return echo.NewHTTPError(http.StatusNotFound)
		}

		u, err := tinyURL.URL(s)
		if err != nil {
			log.Println(err)
			return err
		}

		return c.Redirect(http.StatusTemporaryRedirect, u)
	}
}
