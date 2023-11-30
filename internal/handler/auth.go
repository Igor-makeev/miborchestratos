package handler

import (
	"miborchestrator/internal/entities"
	customerrors "miborchestrator/internal/entities/custom_errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) register(c *gin.Context) {
	var input entities.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Service.Authorization.CreateUser(c.Request.Context(), input)
	if err != nil {
		if _, ok := err.(*customerrors.LoginConflict); ok {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	token, err := h.Service.Authorization.GenerateToken(c.Request.Context(), input.Login, input.Password)
	if err != nil {
		if _, ok := err.(*customerrors.InvalidLoginOrPassword); ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header(authorizationHeader, token)
	c.Status(http.StatusOK)
}

func (h *Handler) login(c *gin.Context) {
	var input entities.SignInInput

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.Service.Authorization.GenerateToken(c.Request.Context(), input.Login, input.Password)
	if err != nil {
		if _, ok := err.(*customerrors.InvalidLoginOrPassword); ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header(authorizationHeader, token)
	c.Status(http.StatusOK)
}
