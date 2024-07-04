package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gogoalish/timetracker/config"
	"github.com/gogoalish/timetracker/internal/clients"
	"github.com/gogoalish/timetracker/internal/controller"
	"github.com/gogoalish/timetracker/internal/logger"
	"github.com/gogoalish/timetracker/internal/repo"
	"github.com/gogoalish/timetracker/internal/server"
	"github.com/gogoalish/timetracker/internal/service"
	_ "github.com/lib/pq"
)

func main() {

	cfg, err := config.New()
	if err != nil {
		log.Fatal("error init config", err)
	}

	l, err := logger.New()
	if err != nil {
		log.Fatal("error init logger", err)
	}
	defer l.Sync()

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		l.Fatal(fmt.Sprintf("error db conn init: %s", err))
	}
	err = db.Ping()
	if err != nil {
		l.Fatal(fmt.Sprint("error db conn ping: ", err))
	}
	defer db.Close()

	peopleRepo := repo.NewPeopleRepo(db)
	apiClient, err := clients.NewAPIService(cfg)
	if err != nil {
		l.Fatal(fmt.Sprint("error api client init: ", err))
	}
	peopleSvc := service.NewPeopleService(peopleRepo, apiClient)

	tasksRepo := repo.NewTasksRepo(db)
	tasksSvc := service.NewTasksService(tasksRepo)

	peopleController := controller.NewPeopleController(peopleSvc)
	tasksController := controller.NewTasksController(tasksSvc)

	router := server.NewRouter(peopleController, tasksController, l)
	httpServer := server.New(cfg, router)
	l.Info(fmt.Sprintf("server is listening on: http://%s:%s", cfg.Host, cfg.Port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("main - signal:" + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Sprint("main - httpServer.Notify: ", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Sprint("main - httpServer.Shutdown: ", err))
	}

}
