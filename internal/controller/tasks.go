package controller

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gogoalish/timetracker/internal/logger"
	"github.com/gogoalish/timetracker/internal/service"
	"go.uber.org/zap"
)

type TasksController struct {
	svc service.TasksService
}

func NewTasksController(svc service.TasksService) *TasksController {
	return &TasksController{
		svc: svc,
	}
}

type createTaskReq struct {
	Description string `json:"description" binding:"required"`
	UserID      int    `json:"user_id" binding:"required,min=1"`
}

// Create godoc
// @Summary Create a new task
// @Description Create a new task with a specific user ID and description
// @Tags Tasks
// @Accept  json
// @Produce  json
// @Param   task  body  createTaskReq  true  "Task description and user ID"
// @Success 200 {integer} int "Task ID"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /tasks/create [post]
func (c *TasksController) Create(ctx *gin.Context) {
	l, ok := logger.FromContext(ctx.Request.Context())
	if !ok {
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrNoLogger))
		return
	}

	var req createTaskReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		l.Error("TasksController - Create - binding error", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := c.svc.CreateTask(ctx, req.UserID, req.Description)
	if err != nil {
		l.Error("TasksController - Create - CreateTask error", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	l.Info("Task created successfully", zap.Int("task_id", int(id)))
	ctx.JSON(http.StatusOK, id)
}

type taskStartReq struct {
	ID int `json:"id" binding:"required,min=1"`
}

// Start godoc
// @Summary Start a task
// @Description Start a task by its ID
// @Tags Tasks
// @Accept  json
// @Produce  json
// @Param   task  body  taskStartReq  true  "Task ID"
// @Success 200
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /tasks/start [put]
func (c *TasksController) Start(ctx *gin.Context) {
	l, ok := logger.FromContext(ctx.Request.Context())
	if !ok {
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrNoLogger))
		return
	}

	var req taskStartReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		l.Error("TasksController - Start - binding error", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	l.Debug("Starting task", zap.Int("task_id", req.ID))

	err := c.svc.StartTask(ctx, req.ID)
	if err != nil {
		l.Error("TasksController - Start - StartTask error", zap.Error(err))
		if errors.Is(err, service.ErrNoResult) {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	l.Info("Task started successfully", zap.Int("task_id", req.ID))
	ctx.Status(http.StatusOK)
}

type taskEndReq struct {
	ID int `json:"id" binding:"required,min=1"`
}

// End godoc
// @Summary End a task
// @Description End a task by its ID
// @Tags Tasks
// @Accept  json
// @Produce  json
// @Param   task  body  taskEndReq  true  "Task ID"
// @Success 200
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /tasks/end [post]
func (c *TasksController) End(ctx *gin.Context) {
	l, ok := logger.FromContext(ctx.Request.Context())
	if !ok {
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrNoLogger))
		return
	}

	var req taskEndReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		l.Error("TasksController - End - binding error", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	l.Debug("Ending task", zap.Int("task_id", req.ID))

	err := c.svc.EndTask(ctx, req.ID)
	if err != nil {
		l.Error("TasksController - End - EndTask error", zap.Error(err))
		if errors.Is(err, service.ErrNoResult) {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	l.Info("Task ended successfully", zap.Int("task_id", req.ID))
	ctx.Status(http.StatusOK)
}

type getOrderedTasksReq struct {
	UserID int    `json:"user_id" binding:"required,min=1"`
	FromDT string `json:"from_dt" binding:"required"`
	ToDT   string `json:"to_dt" binding:"required"`
}

const dateLayout = "2006-01-02 15:04:05"

// Ordered godoc
// @Summary Get ordered tasks
// @Description Get ordered tasks by user ID and date range
// @Tags Tasks
// @Accept  json
// @Produce  json
// @Param   tasks  body  getOrderedTasksReq  true  "User ID and date range"
// @Success 200 {array} service.Task "List of tasks"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /tasks/ordered [get]
func (c *TasksController) Ordered(ctx *gin.Context) {
	l, ok := logger.FromContext(ctx.Request.Context())
	if !ok {
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrNoLogger))
		return
	}

	var req getOrderedTasksReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		l.Error("TasksController - Ordered - binding error", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	from, err := time.Parse(dateLayout, req.FromDT)
	if err != nil {
		l.Error("TasksController - Ordered - time parsing error for from_dt", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	to, err := time.Parse(dateLayout, req.ToDT)
	if err != nil {
		l.Error("TasksController - Ordered - time parsing error for to_dt", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	l.Debug("Fetching ordered tasks", zap.Int("user_id", req.UserID), zap.String("from_dt", req.FromDT), zap.String("to_dt", req.ToDT))

	tasks, err := c.svc.GetOrderedTasks(ctx, req.UserID, from, to)
	if err != nil {
		l.Error("TasksController - Ordered - GetOrderedTasks error", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	l.Info("Ordered tasks fetched successfully", zap.Int("user_id", req.UserID), zap.Int("task_count", len(tasks)))
	ctx.JSON(http.StatusOK, tasks)
}
