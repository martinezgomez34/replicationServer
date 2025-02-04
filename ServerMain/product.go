package serverMain

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

type Product struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Amount  string `json:"amount"`
	CodeBar string `json:"codeBar"`
}

type Cambio struct {
	Accion  string  `json:"accion"`
	Product Product `json:"product"`
}

var (
	bd      []Product
	cambios []Cambio
)

func replicationServer(product Product, accion string) {
	url := fmt.Sprintf("http://localhost:8081/replication?id=%d&name=%s&amount=%s&codeBar=%s&accion=%s",
		product.ID,
		url.QueryEscape(product.Name),
		url.QueryEscape(product.Amount),
		url.QueryEscape(product.CodeBar),
		url.QueryEscape(accion),
	)

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.Status)
	}
}

func productReplication(c *gin.Context) {
	if len(cambios) > 0 {
		lastChange := cambios[len(cambios)-1]
		replicationServer(lastChange.Product, lastChange.Accion)
		c.JSON(http.StatusOK, gin.H{"mensaje": "Hay cambios", "producto": lastChange.Product})
	} else {
		c.JSON(http.StatusOK, gin.H{"mensaje": "No hay cambios nuevos"})
	}
}

func createProduct(c *gin.Context) {
	var newProduct Product
	if err := c.ShouldBindJSON(&newProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newProduct.ID = int64(len(bd) + 1)
	bd = append(bd, newProduct)

	cambios = append(cambios, Cambio{Accion: "create", Product: newProduct})
	replicationServer(newProduct, "create")

	c.JSON(http.StatusCreated, newProduct)
}

func getProduct(c *gin.Context) {
	c.JSON(http.StatusOK, bd)
}

func getCambios(c *gin.Context) {
	if len(cambios) == 0 {
		c.JSON(http.StatusOK, gin.H{"mensaje": "No hay cambios nuevos"})
		return
	}

	response := cambios
	cambios = []Cambio{}

	c.JSON(http.StatusOK, response)
}

func deleteProduct(c *gin.Context) {
	id := c.Param("id")
	var indexToRemove int = -1

	for i, product := range bd {
		if fmt.Sprintf("%d", product.ID) == id {
			indexToRemove = i
			break
		}
	}

	if indexToRemove == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
		return
	}

	productToDelete := bd[indexToRemove]
	bd = append(bd[:indexToRemove], bd[indexToRemove+1:]...)

	client := &http.Client{Timeout: 10 * time.Second}
	replicationURL := fmt.Sprintf("http://localhost:8081/replication/%d", productToDelete.ID)
	req, err := http.NewRequest("DELETE", replicationURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear"})
		return
	}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar"})
		return
	}

	cambios = append(cambios, Cambio{Accion: "delete", Product: productToDelete})

	c.JSON(http.StatusOK, gin.H{"mensaje": "Producto eliminado", "producto": productToDelete})
}


func updateProduct(c *gin.Context) {
	id := c.Param("id")
	var updatedProduct Product
	if err := c.ShouldBindJSON(&updatedProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var productToUpdate *Product
	for i, product := range bd {
		if fmt.Sprintf("%d", product.ID) == id {
			productToUpdate = &bd[i]
			break
		}
	}

	if productToUpdate == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
		return
	}

	productToUpdate.Name = updatedProduct.Name
	productToUpdate.Amount = updatedProduct.Amount
	productToUpdate.CodeBar = updatedProduct.CodeBar

	cambios = append(cambios, Cambio{Accion: "update", Product: *productToUpdate})
	replicationServer(*productToUpdate, "update")

	c.JSON(http.StatusOK, gin.H{"mensaje": "Producto actualizado", "producto": *productToUpdate})
}