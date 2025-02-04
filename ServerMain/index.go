package serverMain

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)


func Run() {
	r := gin.Default()

	r.POST("/product", createProduct)
	r.GET("/product", getProduct)
	r.GET("/cambios", getCambios)
	r.GET("/productReplication", productReplication)
	r.DELETE("/product/:id", deleteProduct)
	r.PUT("/product/:id", updateProduct)
	

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 5 * time.Minute,
		IdleTimeout:  1 * time.Hour,
	}

	if err := srv.ListenAndServe(); err != nil {
		fmt.Println("Error: Server Main hasn't begun")
	}
}
