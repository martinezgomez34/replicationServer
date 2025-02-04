package serverReplication

import (
	"fmt"
	"github.com/gin-gonic/gin"
)
func Run() {
	r := gin.Default()

	r.GET("/replication", getReplicatedProduct)
	r.DELETE("/replication/:id", deleteReplicatedProduct)
	r.PUT("/replication/:id", updateReplicatedProduct)

	if err := r.Run(":8081"); err != nil {
		fmt.Println("Error: Replication Server hasn't begun")
	}
}
