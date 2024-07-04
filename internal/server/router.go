package server

import (
	"github.com/gin-gonic/gin"
	"github.com/gogoalish/timetracker/internal/controller"
	"go.uber.org/zap"
)

func NewRouter(peopleCntrl *controller.PeopleController, taskCntrl *controller.TasksController, l *zap.Logger) *gin.Engine {
	router := gin.New()
	router.Use(RequestLogger(l))
	people := router.Group("/people")
	{
		people.POST("/create", peopleCntrl.Create)
		people.GET("/list", peopleCntrl.List)
		people.PATCH("/update", peopleCntrl.Update)
		people.DELETE("/delete", peopleCntrl.Delete)
	}

	tasks := router.Group("/tasks")
	{
		tasks.POST("/create", taskCntrl.Create)
		tasks.PATCH("/start", taskCntrl.Start)
		tasks.PATCH("/update", taskCntrl.End)
		tasks.GET("/ordered", taskCntrl.Ordered)
	}
	return router
}
