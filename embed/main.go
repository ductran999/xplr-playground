package main

import (
	"log"
	"net/http"
	"play-ground/embed/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	hdl, err := handler.NewManifestHandler()
	if err != nil {
		log.Fatalln("error", err)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	r.POST("/deployment", hdl.GetDeploymentManifest)
	r.POST("/service", hdl.GetServiceManifest)

	// Start server
	if err := r.Run(":8080"); err != http.ErrServerClosed {
		log.Fatalln("start server error", err)
	}
}
