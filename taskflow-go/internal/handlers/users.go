package handlers

import (
	"net/http"
	"taskflow/internal/models"
	"taskflow/internal/repository"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	repo repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	user, err := h.repo.GetByID(c.Request.Context(), userID)
	if err != nil || user == nil {
		RespondError(c, NewAppError(http.StatusNotFound, "User not found"))
		return
	}

	RespondJSON(c, http.StatusOK, user)
}

func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	var req struct {
		FullName string `json:"full_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, NewAppError(http.StatusBadRequest, err.Error()))
		return
	}

	user, err := h.repo.GetByID(c.Request.Context(), userID)
	if err != nil || user == nil {
		RespondError(c, NewAppError(http.StatusNotFound, "User not found"))
		return
	}

	user.FullName = req.FullName
	if err := h.repo.Update(c.Request.Context(), user); err != nil {
		RespondError(c, NewAppError(http.StatusInternalServerError, err.Error()))
		return
	}

	RespondJSON(c, http.StatusOK, user)
}
