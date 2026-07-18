package handlers

import (
	"net/http"
	"taskflow/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TeamHandler struct {
	service *services.TeamService
}

func NewTeamHandler(service *services.TeamService) *TeamHandler {
	return &TeamHandler{service: service}
}

func (h *TeamHandler) CreateTeam(c *gin.Context) {
	userIDStr := c.MustGet("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		RespondError(c, NewAppError(http.StatusBadRequest, "Invalid user ID"))
		return
	}

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, NewAppError(http.StatusBadRequest, err.Error()))
		return
	}

	team, err := h.service.CreateTeam(c.Request.Context(), req.Name, req.Description, userID)
	if err != nil {
		RespondError(c, NewAppError(http.StatusInternalServerError, err.Error()))
		return
	}
	RespondJSON(c, http.StatusCreated, team)
}

func (h *TeamHandler) GetTeams(c *gin.Context) {
	userIDStr := c.MustGet("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		RespondError(c, NewAppError(http.StatusBadRequest, "Invalid user ID"))
		return
	}

	teams, err := h.service.GetUserTeams(c.Request.Context(), userID)
	if err != nil {
		RespondError(c, NewAppError(http.StatusInternalServerError, err.Error()))
		return
	}
	RespondJSON(c, http.StatusOK, teams)
}

func (h *TeamHandler) AddMember(c *gin.Context) {
	ownerIDStr := c.MustGet("user_id").(string)
	ownerID, err := uuid.Parse(ownerIDStr)
	if err != nil {
		RespondError(c, NewAppError(http.StatusBadRequest, "Invalid user ID"))
		return
	}

	teamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		RespondError(c, NewAppError(http.StatusBadRequest, "Invalid team ID"))
		return
	}

	var req struct {
		UserID uuid.UUID `json:"user_id" binding:"required"`
		Role   string    `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, NewAppError(http.StatusBadRequest, err.Error()))
		return
	}

	err = h.service.AddMember(c.Request.Context(), teamID, ownerID, req.UserID, req.Role)
	if err != nil {
		RespondError(c, NewAppError(http.StatusForbidden, err.Error()))
		return
	}
	RespondJSON(c, http.StatusOK, gin.H{"message": "Member added successfully"})
}
