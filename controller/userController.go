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
