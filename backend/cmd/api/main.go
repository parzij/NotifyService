package main

import (
	"os"
	"net/http"
	"github.com/gin-gonic/gin"
	"NotifyService/internal/router"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	gin.SetMode(os.Getenv("GIN_MODE"))

    port := os.Getenv("APP_PORT")

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	router.Router(r)

    srv := &http.Server{Addr: ":" + port, Handler: r}
    srv.ListenAndServe()
}
