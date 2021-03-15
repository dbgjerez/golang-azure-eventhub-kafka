package main

import (
	"context"
	"golang-azure-eventhub-kafka/adapters"
	"golang-azure-eventhub-kafka/controllers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	srv := initServer()
	go listenAndServe(srv)
	kafkaClient := adapters.NewConnection()
	ctxKafka, cancelKafka := context.WithCancel(context.Background())
	go kafkaClient.Suscribe(ctxKafka)

	/** Wait for exit. We need stop server and kafka thread */
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	cancelKafka()
	kafkaClient.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}

func listenAndServe(srv *http.Server) {
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}

func initServer() *http.Server {
	router := gin.Default()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", controllers.HealthControllerHandler())
	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"msg": "Not found"})
	})
	return srv
}
