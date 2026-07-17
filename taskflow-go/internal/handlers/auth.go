package handlers

import (
	"net/http"
	"taskflow/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
		FullName string `json:"full_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, NewAppError(http.StatusBadRequest, err.Error()))
		return
	}

	user, token, err := h.service.Register(c.Request.Context(), req.Email, req.Password, req.FullName)
	if err != nil {
		RespondError(c, NewAppError(http.StatusBadRequest, err.Error()))
		return
	}

	RespondJSON(c, http.StatusCreated, gin.H{"user": user, "token": token})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, NewAppError(http.StatusBadRequest, err.Error()))
		return
	}

	user, token, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		RespondError(c, NewAppError(http.StatusUnauthorized, err.Error()))
		return
	}

	RespondJSON(c, http.StatusOK, gin.H{"user": user, "token": token})
}
