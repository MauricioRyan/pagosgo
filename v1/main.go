package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	cors "github.com/itsjamie/gin-cors"
	"github.com/mauricioryan/pagosgo/routes"
	"github.com/mauricioryan/pagosgo/tools/env"
)

func main() {
	if len(os.Args) > 1 {
		env.Load(os.Args[1])
	}

	// Hoy gin usa v8, para actualizar gin validator a v9.
	binding.Validator = new(defaultValidator)

	server := gin.Default()

	server.Use(gzip.Gzip(gzip.DefaultCompression))

	// cuando tengoque tener configurado cors??
	server.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	//server.Use(static.Serve("/", static.LocalFile(env.Get().WWWWPath, true)))

	server.GET("/v1/pagos", routes.Pagos)
	server.POST("/v1/pagos", routes.NewPago)
	//server.GET("/v1/pagos/:pagoID", routes.GetPago)

	server.Run(fmt.Sprintf(":%d", env.Get().Port))
}
