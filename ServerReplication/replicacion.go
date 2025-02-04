package serverReplication

import (
	"fmt"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

type Product struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Amount string `json:"amount"`
	CodeBar string `json:"codeBar"`
}

var bdReplication []Product

func getReplicatedUsers(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	name := c.DefaultQuery("name", "")
	amount := c.DefaultQuery("amount", "")
	codeBar := c.DefaultQuery("codeBar", "")
	accion := c.DefaultQuery("accion", "")

	if id != "" && name != "" && amount != "" && codeBar != "" && accion != "" {
		productID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
			return
		}
		for i, product := range bdReplication {
			if product.ID == productID {
				bdReplication[i] = Product{ID: productID, Name: name, Amount: amount, CodeBar: codeBar}
				fmt.Println("Producto actualizado:", bdReplication[i])
				c.JSON(http.StatusOK, gin.H{"mensaje": "Producto actualizado:", "producto": bdReplication[i]})
				return
			}
		}

		newProduct := Product{ID: productID, Name: name, Amount: amount, CodeBar: codeBar}
		bdReplication = append(bdReplication, newProduct)
		fmt.Println("Producto replicado:", newProduct)
		c.JSON(http.StatusOK, bdReplication)
		return
	}

	c.JSON(http.StatusOK, bdReplication)
}

func deleteReplicatedProduct(c *gin.Context) {
	id := c.Param("id")
	productID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	for i, product := range bdReplication {
		if product.ID == productID {
			bdReplication = append(bdReplication[:i], bdReplication[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"mensaje": "Producto eliminado:", "producto": product})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
}

func updateReplicatedProduct(c *gin.Context) {
	id := c.Param("id")
	productID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	var updatedProduct Product
	if err := c.ShouldBindJSON(&updatedProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	for i, product := range bdReplication {
		if product.ID == productID {
			bdReplication[i] = updatedProduct
			c.JSON(http.StatusOK, gin.H{"mensaje": "Producto actualizado:", "producto": bdReplication[i]})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
}
