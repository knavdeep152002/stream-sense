package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Auth struct {
	DB *gorm.DB
}

func CreateAuth(db *gorm.DB) *Auth {
	return &Auth{
		DB: db,
	}
}

// @Summary		Create a user
// @Description	Create a user
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			user	body		authInput	true	"User"
// @Success		200		{string}	token
// @Router			/auth/register [post]
func (a Auth) CreateUser(c *gin.Context) {
	var authInput authInput
	if err := c.ShouldBindJSON(&authInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := a.createUserInDB(&authInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := a.getSignedToken(&authInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// @Summary		Login
// @Description	Login
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			user	body		authInput	true	"User"
// @Success		200		{string}	token
// @Router			/auth/login [post]
func (a *Auth) Login(c *gin.Context) {

	var authInput authInput
	if err := c.ShouldBindJSON(&authInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := a.getSignedToken(&authInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
