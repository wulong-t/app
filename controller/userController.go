package controller

import (
	"Server/database"
	"Server/models"
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetUserByID
// @Summary Get user by ID
// @Description Get user by ID
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Security BearerAuth
// @Success 200 {object} models.UserModel
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /user/getuser/{id} [get]
func GetUserByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	userID := c.Param("id")
	user, err := gorm.G[models.UserModel](database.DB).Where("id = ?", userID).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// UpdateUser
// @Summary Update user
// @Description Update user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body models.UpdateUser true "User"
// @Security BearerAuth
// @Success 200 {object} models.UpdateUser
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /user/updateuser/{id} [patch]
func UpdateUser(c *gin.Context) {
	updateID := c.Param("id")
	userID := c.GetString("userID")
	if updateID != userID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to update this user"})
		return
	}
	var user models.UpdateUser
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	updateUser := map[string]any{}
	if user.Name != "" {
		updateUser["name"] = user.Name
	}
	if user.Bio != "" {
		updateUser["bio"] = user.Bio
	}
	if user.ImageURL != "" {
		updateUser["image_url"] = user.ImageURL
	}
	result := database.DB.WithContext(ctx).Model(&models.UserModel{}).Where("id = ?", updateID).Updates(updateUser)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"details": result.Error.Error(),
		})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}
