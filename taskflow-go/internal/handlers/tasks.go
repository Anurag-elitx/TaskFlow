package handlers

import (
	"net/http"
	"strconv"
	"taskflow/internal/models"
	"taskflow/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TaskHandler struct {
	service *services.TaskService
}

func NewTaskHandler(service *services.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	userIDStr := c.MustGet("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)
	teamID, err := uuid.Parse(c.Param("teamId"))
	if err != nil {
		RespondError(c, NewAppError(http.StatusBadRequest, "Invalid team ID"))
		return
	}

	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		RespondError(c, NewAppError(http.StatusBadRequest, err.Error()))
		return
	}
	task.TeamID = teamID
	task.CreatedBy = userID

	if err := h.service.CreateTask(c.Request.Context(), &task); err != nil {
		RespondError(c, NewAppError(http.StatusForbidden, err.Error()))
		return
	}
	RespondJSON(c, http.StatusCreated, task)
}

func (h *TaskHandler) GetTasks(c *gin.Context) {
	userIDStr := c.MustGet("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)
	teamID, err := uuid.Parse(c.Param("teamId"))
	if err != nil {
		RespondError(c, NewAppError(http.StatusBadRequest, "Invalid team ID"))
		return
	}

	filters := models.TaskFilters{}
	if status := c.Query("status"); status != "" {
		filters.Status = &status
	}
	if priority := c.Query("priority"); priority != "" {
		filters.Priority = &priority
	}
	if assigneeIDStr := c.Query("assignee_id"); assigneeIDStr != "" {
		id, err := uuid.Parse(assigneeIDStr)
		if err == nil {
			filters.AssigneeID = &id
		}
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	filters.Page = page
	filters.Limit = limit

	tasks, total, err := h.service.GetTasksByTeam(c.Request.Context(), teamID, userID, filters)
	if err != nil {
		RespondError(c, NewAppError(http.StatusForbidden, err.Error()))
		return
	}

	totalPages := total / limit
	if total%limit != 0 {
		totalPages++
	}

	res := PaginatedResponse{
		Data:       tasks,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
	RespondJSON(c, http.StatusOK, res)
}

func (h *TaskHandler) GetTask(c *gin.Context) {
	userIDStr := c.MustGet("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		RespondError(c, NewAppError(http.StatusBadRequest, "Invalid task ID"))
		return
	}

	task, err := h.service.GetTaskByID(c.Request.Context(), taskID, userID)
	if err != nil {
		RespondError(c, NewAppError(http.StatusNotFound, err.Error()))
		return
	}
	RespondJSON(c, http.StatusOK, task)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	userIDStr := c.MustGet("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		RespondError(c, NewAppError(http.StatusBadRequest, "Invalid task ID"))
		return
	}

	var updates models.Task
	if err := c.ShouldBindJSON(&updates); err != nil {
		RespondError(c, NewAppError(http.StatusBadRequest, err.Error()))
		return
	}

	task, err := h.service.UpdateTask(c.Request.Context(), taskID, userID, &updates)
	if err != nil {
		RespondError(c, NewAppError(http.StatusForbidden, err.Error()))
		return
	}
	RespondJSON(c, http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	userIDStr := c.MustGet("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		RespondError(c, NewAppError(http.StatusBadRequest, "Invalid task ID"))
		return
	}

	err = h.service.DeleteTask(c.Request.Context(), taskID, userID)
	if err != nil {
		RespondError(c, NewAppError(http.StatusForbidden, err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}
