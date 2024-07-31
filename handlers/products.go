package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iadityanath8/gopi/database"
	"github.com/iadityanath8/gopi/models"
)

func GetProduct(c *gin.Context) {
	var product *models.Product
	if err := database.DB.First(&product, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product Not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func GetProducts(c *gin.Context) {

	var products []models.Product
	name := c.Query("name")

	if name != "" {
		if err := database.DB.Where("code = ?", name).First(&products).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Produgt not found !!"})
			return
		}
		c.JSON(http.StatusOK, products)
	}

	database.DB.Find(&products)
	c.JSON(http.StatusOK, products)
}

func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := database.DB.Create(&product)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

func UpdataProduct(c *gin.Context) {
	var product models.Product
	if err := database.DB.First(&product, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Save(&product)
	c.JSON(http.StatusOK, product)
}

func DeleteProduct(c *gin.Context) {
	var product models.Product
	if err := database.DB.First(&product, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
	}

	database.DB.Delete(&product)
	c.Status(http.StatusNoContent)
}
