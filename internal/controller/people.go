package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gogoalish/timetracker/internal/logger"
	"github.com/gogoalish/timetracker/internal/service"

	"go.uber.org/zap"
)

type PeopleController struct {
	svc service.PeopleService
}

func NewPeopleController(svc service.PeopleService) *PeopleController {
	return &PeopleController{
		svc: svc,
	}
}

type createPersonReq struct {
	PassportNumber int `json:"passport_number" binding:"required"`
	PassportSerie  int `json:"passport_serie" binding:"required"`
}

var ErrNoLogger = errors.New("logger not found in context")

// Create godoc
// @Summary Create a new person
// @Description Create a new person with given passport details
// @Tags People
// @Accept json
// @Produce json
// @Param person body createPersonReq true "Person details"
// @Success 200 {integer} int "Person ID"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /people/create [post]
func (c *PeopleController) Create(ctx *gin.Context) {
	l, ok := logger.FromContext(ctx.Request.Context())
	if !ok {
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrNoLogger))
		return
	}

	var req createPersonReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		l.Error("PeopleCntrl - Create - binding error", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	l.Info("Received request to create person", zap.String("passport_serie", fmt.Sprintf("%d", req.PassportSerie)), zap.String("passport_number", fmt.Sprintf("%d", req.PassportNumber)))

	id, err := c.svc.CreatePerson(ctx, req.PassportSerie, req.PassportNumber)
	if err != nil {
		l.Error("PeopleCntrl - Create - CreatePerson error", zap.Error(err))
		if errors.Is(err, service.ErrAlreadyExists) {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	l.Info("Person created successfully", zap.Int32("id", id))
	ctx.JSON(http.StatusOK, id)
}

type listPeopleReq struct {
	Limit          *int32 `form:"limit" binding:"omitempty,min=1"`
	Page           *int32 `form:"page" binding:"omitempty,min=1"`
	PassportSerie  *int32 `form:"passport_serie" binding:"omitempty,min=1"`
	PassportNumber *int32 `form:"passport_number" binding:"omitempty,min=1"`
	Surname        string `form:"surname"`
	Name           string `form:"name"`
	Patronymic     string `form:"patronymic"`
	Address        string `form:"address"`
}

// List godoc
// @Summary List people
// @Description List people with optional filters
// @Tags People
// @Accept json
// @Produce json
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param passport_serie query int false "Passport Serie"
// @Param passport_number query int false "Passport Number"
// @Param surname query string false "Surname"
// @Param name query string false "Name"
// @Param patronymic query string false "Patronymic"
// @Param address query string false "Address"
// @Success 200 {array} service.Person "List of people"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /people/list [get]
func (c *PeopleController) List(ctx *gin.Context) {
	l, ok := logger.FromContext(ctx.Request.Context())
	if !ok {
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrNoLogger))
		return
	}

	var req listPeopleReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		l.Error("PeopleCntrl - List - binding error", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	l.Debug("Listing people with filters", zap.Any("filters", req))

	people, err := c.svc.ListPeople(ctx, service.Filter{
		Limit:          req.Limit,
		Offset:         req.Page,
		PassportSerie:  req.PassportSerie,
		PassportNumber: req.PassportNumber,
		Surname:        req.Surname,
		Name:           req.Name,
		Patronymic:     req.Patronymic,
		Address:        req.Address,
	})
	if err != nil {
		l.Error("PeopleCntrl - List - ListPeople error", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	l.Info("People listed successfully", zap.Int("count", len(people)))
	ctx.JSON(http.StatusOK, people)
}

type updatePersonReq struct {
	ID             int32  `json:"id" binding:"required,min=1"`
	PassportSerie  *int32 `json:"passport_serie" binding:"omitempty,min=1"`
	PassportNumber *int32 `json:"passport_number" binding:"omitempty,min=1"`
	Surname        string `json:"surname"`
	Name           string `json:"name"`
	Patronymic     string `json:"patronymic"`
	Address        string `json:"address"`
}

// Update godoc
// @Summary Update a person
// @Description Update a person's details
// @Tags People
// @Accept json
// @Produce json
// @Param person body updatePersonReq true "Person details"
// @Success 200 "Success"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /people/update [put]
func (c *PeopleController) Update(ctx *gin.Context) {
	l, ok := logger.FromContext(ctx.Request.Context())
	if !ok {
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrNoLogger))
		return
	}

	var req updatePersonReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		l.Error("PeopleCntrl - Update - binding error", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	l.Debug("Updating person with ID", zap.Int32("id", req.ID))

	err := c.svc.UpdatePerson(ctx, service.UpdatedPerson{
		ID:             req.ID,
		PassportSerie:  req.PassportSerie,
		PassportNumber: req.PassportNumber,
		Surname:        req.Surname,
		Name:           req.Name,
		Patronymic:     req.Patronymic,
		Address:        req.Address,
	})
	if err != nil {
		l.Error("PeopleCntrl - Update - UpdatePerson error", zap.Error(err))
		if errors.Is(err, service.ErrNoResult) {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	l.Info("Person updated successfully", zap.Int32("id", req.ID))
	ctx.Status(http.StatusOK)
}

type deletePersonReq struct {
	ID int32 `json:"id" binding:"required,min=1"`
}

// Delete godoc
// @Summary Delete a person
// @Description Delete a person by ID
// @Tags People
// @Accept json
// @Produce json
// @Param person body deletePersonReq true "Person ID"
// @Success 200 "Success"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /people/delete [delete]
func (c *PeopleController) Delete(ctx *gin.Context) {
	l, ok := logger.FromContext(ctx.Request.Context())
	if !ok {
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrNoLogger))
		return
	}

	var req deletePersonReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		l.Error("PeopleCntrl - Delete - binding error", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	l.Debug("Deleting person with ID", zap.Int32("id", req.ID))

	err := c.svc.DeletePerson(ctx, req.ID)
	if err != nil {
		l.Error("PeopleCntrl - Delete - DeletePerson error", zap.Error(err))
		if errors.Is(err, service.ErrNoResult) {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	l.Info("Person deleted successfully", zap.Int32("id", req.ID))
	ctx.Status(http.StatusOK)
}
