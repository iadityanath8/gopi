package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iadityanath8/gopi/database"
	"github.com/iadityanath8/gopi/middleware"
	"github.com/iadityanath8/gopi/models"
)

func AddItemToCart(c *gin.Context) {
	// local model
	var requestBody models.RequestBody
	// local model end

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	curr_usr := middleware.GetCurrentUser(c)
	var user models.User

	if err := database.DB.First(&user, "username = ?", curr_usr).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var product models.Product
	if err := database.DB.First(&product, requestBody.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// check product stock
	if product.Stock < requestBody.Quantity {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not enough Stock"})
		return
	}

	var cart models.Cart
	if err := database.DB.First(&cart, "user_id = ?", user.ID).Error; err != nil {
		cart.UserID = user.ID
		database.DB.Create(&cart)
	}

	// check if product is already in a cart
	var existingCartItem models.CartItem
	if err := database.DB.Where("cart_id = ? AND product_id = ?", cart.ID, product.ID).First(&existingCartItem).Error; err == nil {
		// Update quantity if product is already in cart
		existingCartItem.Quantity += requestBody.Quantity
		if err := database.DB.Save(&existingCartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update cart item"})
			return
		}
	} else {
		cartItem := models.CartItem{
			ProductID: product.ID,
			Quantity:  requestBody.Quantity,
			CartID:    cart.ID,
		}
		if err := database.DB.Create(&cartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add item to cart"})
			return
		}
	}

	product.Stock -= requestBody.Quantity
	if err := database.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update the stock price"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Item added to cart"})
}

func RemoveItemCart(c *gin.Context) {
	var requestBody models.RequestBodyDel

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	curr_usr := middleware.GetCurrentUser(c)
	var user models.User

	if err := database.DB.First(&user, "username = ?", curr_usr).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var cart models.Cart
	if err := database.DB.First(&cart, "user_id = ?", user.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	var cartItem models.CartItem
	if err := database.DB.Where("cart_id = ? AND product_id = ?", cart.ID, requestBody.ProductID).First(&cartItem).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found in cart"})
		return
	}

	var product models.Product
	if err := database.DB.First(&product, cartItem.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product Not Found"})
	}

	tx := database.DB.Begin()

	product.Stock += cartItem.Quantity
	if err := tx.Save(&product).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update product stock"})
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart"})
}

func GetCartItems(c *gin.Context) {
	curr_usr := middleware.GetCurrentUser(c)
	var user models.User

	if err := database.DB.First(&user, "username = ?", curr_usr).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User Not Found"})
		return
	}

	var cart models.Cart
	if err := database.DB.Preload("Items").First(&cart, "user_id = ?", user.ID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"items": []models.CartItemResponse{}})
		return
	}

	var cartItems []models.CartItemResponse

	for _, item := range cart.Items {
		var product models.Product
		if err := database.DB.First(&product, item.ProductID).Error; err != nil {
			continue
		}
		cartItems = append(cartItems, models.CartItemResponse{
			ProductId: item.ProductID,
			Name:      product.Name,
			Price:     float64(product.Price),
			Quantity:  item.Quantity,
		})
	}

	c.JSON(http.StatusOK, gin.H{"items": cartItems})
}
