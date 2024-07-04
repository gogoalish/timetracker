package server

import (
	"github.com/gin-gonic/gin"
	"github.com/gogoalish/timetracker/internal/controller"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

func NewRouter(peopleCntrl *controller.PeopleController, taskCntrl *controller.TasksController, l *zap.Logger) *gin.Engine {
	router := gin.New()
	router.Use(RequestLogger(l))
	people := router.Group("/people")
	{
		people.POST("/create", peopleCntrl.Create)
		people.GET("/list", peopleCntrl.List)
		people.PUT("/update", peopleCntrl.Update)
		people.DELETE("/delete", peopleCntrl.Delete)
	}

	tasks := router.Group("/tasks")
	{
		tasks.POST("/create", taskCntrl.Create)
		tasks.POST("/start", taskCntrl.Start)
		tasks.POST("/update", taskCntrl.End)
		tasks.GET("/ordered", taskCntrl.Ordered)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}
